package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"code.emcdtech.com/emcd/sdk/app"
	sdkCfg "code.emcdtech.com/emcd/sdk/config"
	sdkErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/email/internal/config"
	grpccontroller "code.emcdtech.com/emcd/service/email/internal/controller/grpc"
	"code.emcdtech.com/emcd/service/email/internal/healthchecker"
	"code.emcdtech.com/emcd/service/email/internal/mail_sender"
	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
	"code.emcdtech.com/emcd/service/email/internal/repository/delegate"
	"code.emcdtech.com/emcd/service/email/internal/repository/pg"
	"code.emcdtech.com/emcd/service/email/internal/service"
	pb "code.emcdtech.com/emcd/service/email/protocol/email"
)

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    string = "email"
	serviceVersion string = "local-build"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	err := log.Init(ctx)
	fmt.Println(err)
	log.Info(ctx, "starting email server app")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Info(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	cfg, err := sdkCfg.New[config.Config]()
	if err != nil {
		// nolint: gocritic
		log.Fatal(ctx, "get config: %s", err.Error())
		return
	}
	log.Info(ctx, "start with AvailableDomains: %+v", cfg.Domains)

	deps, err := NewDeps(ctx, cfg)
	if err != nil {
		// nolint: gocritic
		log.Fatal(ctx, "build deps: %s", err.Error())
		return
	}
	otel.SetTracerProvider(deps.Infrastructure.Tracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	defer func() {
		// оставляем пока так, до момента как в app появятся defer handler
		log.Info(ctx, "pgx pool is gracefully shutdown...")
		deps.Infrastructure.DB.Close()
		log.Info(ctx, "pgx pool has gracefully shutdown.")
	}()

	var grpcSrv *grpc.Server

	err = app.New().
		WithTracing(deps.Infrastructure.Tracer).
		WithPprof(deps.Infrastructure.HTTPSrv).
		WithMetrics(deps.Infrastructure.HTTPSrv /*, deps.Service.Email*/).
		WithGRPC(
			func(opts ...grpc.ServerOption) *grpc.Server {
				grpcSrv = deps.Infrastructure.GRPCSrv(opts...)
				return grpcSrv
			},
			healthchecker.NewCommon(deps.Infrastructure.DB),
		).
		WithWorker(
			app.NewHTTPSrvWorker(deps.Infrastructure.HTTPSrv, cfg.HTTP.ListenAddr, shutdownTimeout),
			app.WorkerFn(
				func(ctx context.Context) error {
					return app.NewGRPCSrvWorker(grpcSrv, cfg.GRPC.ListenAddr).Run(ctx)
				},
			),
		).
		WithWorker(deps.Workers...).
		Run(ctx)
	if err != nil {
		log.Panic(ctx, "abnormal stopped email server app: %s", err.Error())
	}
	log.Info(ctx, "graceful stopped email server app")
}

// Deps dependencies
type Deps struct {
	Infrastructure struct {
		Tracer  trace.TracerProvider
		HTTPSrv *echo.Echo
		GRPCSrv app.GRPCServerFactory
		DB      *pgxpool.Pool
	}
	GRPCClient struct {
		WhiteLabelClient whitelabel.WhitelabelServiceClient
		ProfileClient    profile.ProfileServiceClient
	}
	Repository struct {
		EmailMessages          repository.EmailMessages
		Template               repository.Template
		WhiteLabel             repository.Whitelabel
		WhiteLabelEventClients repository.WhiteLabelEventClients
		Profile                repository.Profile
		ProvideSettings        repository.ProvideSettings
	}
	EMailProvider struct {
		Mailgun *mail_sender.Mailgun
		SMTP    *mail_sender.Smtp
	}
	Service struct {
		Email       service.Email
		EmailSender *service.EmailSender
	}
	Workers []app.Worker
}

func NewDeps(ctx context.Context, cfg config.Config) (Deps, error) {
	var deps Deps
	if err := deps.makeInfrastructure(ctx, cfg); err != nil {
		return deps, fmt.Errorf("configurer infrastructure deps: %w", err)
	}
	if err := deps.makeGRPCClient(cfg); err != nil {
		return deps, fmt.Errorf("configurer grpc client: %w", err)
	}
	deps.makeRepository()
	deps.makeService(cfg)

	deps.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		pbEmailController := grpccontroller.NewEmail(deps.Service.Email)
		pbEmailProviderSettingsController := grpccontroller.NewSettingsController(deps.Repository.ProvideSettings)
		pbEmailTemplate := grpccontroller.NewTemplateController(deps.Repository.Template)
		opts = append(opts, grpc.UnaryInterceptor(sdkErr.ServerUnaryInterceptor))
		grpcSrv := grpc.NewServer(opts...)
		pb.RegisterEmailServiceServer(grpcSrv, pbEmailController)
		pb.RegisterEmailProviderSettingsServiceServer(grpcSrv, pbEmailProviderSettingsController)
		pb.RegisterEmailTemplateServiceServer(grpcSrv, pbEmailTemplate)
		if !cfg.Environment.IsProduction() {
			reflection.Register(grpcSrv)
		}
		return grpcSrv
	}
	return deps, nil
}

func (d *Deps) makeInfrastructure(ctx context.Context, cfg config.Config) (err error) {
	tracer, err := cfg.Tracing.New(serviceName, serviceVersion, cfg.Environment.Name)
	if err != nil {
		return fmt.Errorf("new apm tracer: %w", err)
	}
	d.Infrastructure.Tracer, err = apmotel.NewTracerProvider(apmotel.WithAPMTracer(tracer))
	if err != nil {
		return fmt.Errorf("new tracer provider: %w", err)
	}

	d.Infrastructure.DB, err = cfg.PGXPool.New(ctx, d.Infrastructure.Tracer)
	if err != nil {
		return fmt.Errorf("new pgxpool: %w", err)
	}
	d.Infrastructure.HTTPSrv = echo.New()
	d.Infrastructure.HTTPSrv.HideBanner = true
	d.Infrastructure.HTTPSrv.HidePort = true
	return nil
}

func (d *Deps) makeGRPCClient(cfg config.Config) error {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			sdkErr.ClientUnaryInterceptor,
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(d.Infrastructure.Tracer),
			),
		),
	}
	cc, err := grpc.NewClient(cfg.WhiteLabelAddress, options...)
	if err != nil {
		return fmt.Errorf("dial white lable client: %w", err)
	}
	d.GRPCClient.WhiteLabelClient = whitelabel.NewWhitelabelServiceClient(cc)

	cc, err = grpc.NewClient(cfg.ProfileAddress, options...)
	if err != nil {
		return fmt.Errorf("dial white profile: %w", err)
	}
	d.GRPCClient.ProfileClient = profile.NewProfileServiceClient(cc)
	return nil
}

func (d *Deps) makeRepository() {
	d.Repository.EmailMessages = pg.NewEmailMessages(d.Infrastructure.DB)
	d.Repository.ProvideSettings = pg.NewProvideSettingsStore(d.Infrastructure.DB)
	d.Repository.Template = pg.NewTemplateStore(d.Infrastructure.DB)
	d.Repository.WhiteLabel = delegate.NewWhitelabel(d.GRPCClient.WhiteLabelClient)
	d.Repository.WhiteLabelEventClients = delegate.NewWhiteLabelEventClients(d.GRPCClient.WhiteLabelClient)
	d.Repository.Profile = delegate.NewProfileClient(d.GRPCClient.ProfileClient)
}

func (d *Deps) makeService(cfg config.Config) {
	d.Service.EmailSender = service.NewEmailSender(
		d.Repository.Template,
		d.Repository.ProvideSettings,
		d.Repository.EmailMessages,
	)

	d.Service.Email = service.NewEmail(
		d.Service.EmailSender,
		d.Repository.EmailMessages,
		d.Repository.WhiteLabel,
		d.Repository.WhiteLabelEventClients,
		d.Repository.Profile,
		map[model.CodeTemplate]string{
			model.WalletChangedAddress: cfg.ChangeWalletAddressLink,
		},
		cfg.Domains,
	)
}
