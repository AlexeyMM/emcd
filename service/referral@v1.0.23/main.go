// The main package is the entry point of the program.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.emcdtech.com/emcd/sdk/app"
	skdConfig "code.emcdtech.com/emcd/sdk/config"
	sdkErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"
	"code.emcdtech.com/emcd/service/coin/protocol/coin"
	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	promocodePB "code.emcdtech.com/emcd/service/promocode/protocol/promocode"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	coinClient "code.emcdtech.com/emcd/service/referral/internal/clients/coin"
	"code.emcdtech.com/emcd/service/referral/internal/clients/promocode"
	"code.emcdtech.com/emcd/service/referral/internal/config"
	"code.emcdtech.com/emcd/service/referral/internal/healthchecker"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
	"code.emcdtech.com/emcd/service/referral/internal/server"
	"code.emcdtech.com/emcd/service/referral/internal/service"
	"code.emcdtech.com/emcd/service/referral/protocol/default_settings"
	"code.emcdtech.com/emcd/service/referral/protocol/default_users_settings"
	"code.emcdtech.com/emcd/service/referral/protocol/default_whitelabel_settings"
	pbReferal "code.emcdtech.com/emcd/service/referral/protocol/referral"
	"code.emcdtech.com/emcd/service/referral/protocol/reward"
)

const shutdownTimeout = 5 * time.Second

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    string = "referral"
	serviceVersion string = "local-build"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	log.Info(ctx, "starting referral server app")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Info(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	cfg, err := skdConfig.New[config.Config]()
	if err != nil {
		log.Fatal(ctx, "new config: %s", err.Error())
	}

	deps, err := NewDeps(ctx, cfg)
	if err != nil {
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
		WithMetrics(deps.Infrastructure.HTTPSrv).
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
		log.Fatal(ctx, "abnormal stopped referral server app: %s", err.Error())
	}
	log.Info(ctx, "graceful stopped referral server app")
}

// Deps dependencies
type Deps struct {
	Infrastructure struct {
		Tracer     trace.TracerProvider
		HTTPSrv    *echo.Echo
		GRPCSrv    app.GRPCServerFactory
		DB         *pgxpool.Pool
		Transactor pgTx.PgxTransactor
	}
	GRPCClient struct {
		profile   profile.ProfileServiceClient
		promoCode *promocode.Client
		coin      *coinClient.Client
	}
	Repository struct {
		DefaultSettings           repository.DefaultSettings
		DefaultUserSettings       *repository.DefaultUserSettings
		DefaultWhitelabelSettings repository.DefaultWhitelabelSettings
		Referral                  repository.Referral
		Profile                   repository.Profile
	}
	Service struct {
		DefaultSettings           service.DefaultSettings
		DefaultUsersSettings      *service.DefaultUsersSettings
		DefaultWhitelabelSettings service.DefaultWhitelabelSettings
		Referral                  server.ReferralService
		Reward                    service.Reward
	}
	Workers    []app.Worker
	Collectors []prometheus.Collector
}

func NewDeps(ctx context.Context, cfg config.Config) (Deps, error) {
	var deps Deps
	if err := deps.makeInfrastructure(ctx, cfg); err != nil {
		return deps, fmt.Errorf("configurer infrastructure deps: %w", err)
	}
	if err := deps.makeGRPCClients(cfg); err != nil {
		return deps, fmt.Errorf("configurer grpc clients deps: %w", err)
	}
	deps.makeRepository()
	deps.makeService()

	deps.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		pbDefaultSettingsController := server.NewDefaultSettings(deps.Service.DefaultSettings)
		pbDefaultUsersSettingsController := server.NewDefaultUsersSettings(deps.Service.DefaultUsersSettings)
		pbDefaultWhitelabelSettingsController := server.NewDefaultWhitelabelSettings(deps.Service.DefaultWhitelabelSettings)
		pbReferralController := server.NewReferralServer(
			deps.Service.Referral,
			deps.Service.DefaultSettings,
			deps.Service.DefaultWhitelabelSettings,
		)
		pbRewardController := server.NewReward(
			deps.Service.Reward,
			deps.Service.Referral,
		)
		opts = append(opts, grpc.UnaryInterceptor(sdkErr.ServerUnaryInterceptor))
		grpcSrv := grpc.NewServer(opts...)
		default_settings.RegisterDefaultSettingsServiceServer(grpcSrv, pbDefaultSettingsController)
		default_users_settings.RegisterDefaultUsersSettingsServiceServer(grpcSrv, pbDefaultUsersSettingsController)
		default_whitelabel_settings.RegisterDefaultWhitelabelSettingsServiceServer(grpcSrv, pbDefaultWhitelabelSettingsController)
		pbReferal.RegisterReferralServiceServer(grpcSrv, pbReferralController)
		reward.RegisterRewardServiceServer(grpcSrv, pbRewardController)
		reflection.Register(grpcSrv)
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
		return fmt.Errorf("pool to referral: %w", err)
	}
	d.Infrastructure.Transactor = pgTx.NewPgxTransactor(d.Infrastructure.DB)
	d.Infrastructure.HTTPSrv = echo.New()
	d.Infrastructure.HTTPSrv.HideBanner = true
	d.Infrastructure.HTTPSrv.HidePort = true
	return nil
}

func (d *Deps) makeRepository() {
	d.Repository.DefaultSettings = repository.NewDefaultSettings(d.Infrastructure.DB)
	d.Repository.DefaultUserSettings = repository.NewDefaultUserSettings(d.Infrastructure.DB)
	d.Repository.DefaultWhitelabelSettings = repository.NewDefaultWhitelabelSettings(d.Infrastructure.DB)
	d.Repository.Referral = repository.NewReferral(d.Infrastructure.Transactor, d.Infrastructure.DB)
	d.Repository.Profile = repository.NewProfileRepository(d.GRPCClient.profile)
}

func (d *Deps) makeService() {
	d.Service.DefaultSettings = service.NewDefaultSettings(d.Repository.DefaultSettings)
	d.Service.DefaultUsersSettings = service.NewDefaultUsersSettings(d.Repository.DefaultUserSettings,
		d.GRPCClient.promoCode, d.GRPCClient.coin)
	d.Service.DefaultWhitelabelSettings = service.NewDefaultWhitelabelSettings(d.Repository.DefaultWhitelabelSettings)
	d.Service.Referral = service.NewReferralService(d.Repository.Referral, d.Repository.Profile)
	d.Service.Reward = service.NewReward(
		d.Repository.Referral,
		d.Repository.Profile,
		d.Repository.DefaultSettings,
		d.Repository.DefaultWhitelabelSettings,
	)
}

func (d *Deps) makeGRPCClients(cfg config.Config) error {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(d.Infrastructure.Tracer),
			),
		),
	}
	conn, err := grpc.NewClient(cfg.GRPCClient.Profile, dialOptions...)
	if err != nil {
		return fmt.Errorf("new profile conn: %w", err)
	}
	d.GRPCClient.profile = profile.NewProfileServiceClient(conn)

	promoCodeConnection, err := grpc.NewClient(cfg.GRPCClient.PromoCode, dialOptions...)
	if err != nil {
		return fmt.Errorf("new promo code connetion: %w", err)
	}
	promoCodeClient := promocodePB.NewPromocodeClient(promoCodeConnection)
	d.GRPCClient.promoCode = promocode.NewClient(promoCodeClient)

	coinConnection, err := grpc.NewClient(cfg.GRPCClient.Coin, dialOptions...)
	if err != nil {
		return fmt.Errorf("new coin connetion: %w", err)
	}

	coinClientPB := coin.NewCoinServiceClient(coinConnection)
	d.GRPCClient.coin = coinClient.NewClient(coinClientPB)
	return nil
}
