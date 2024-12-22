package main

import (
	"code.emcdtech.com/emcd/service/profile/internal/pkg/nonce"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	walletPb "code.emcdtech.com/emcd/blockchain/wallet/protocol/emcd_wallet_pb"
	"code.emcdtech.com/emcd/sdk/app"
	sdkConfig "code.emcdtech.com/emcd/sdk/config"
	sdkErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	userAccountRepository "code.emcdtech.com/emcd/service/accounting/repository"
	coinPb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	emailPb "code.emcdtech.com/emcd/service/email/protocol/email"
	"code.emcdtech.com/emcd/service/referral/protocol/default_settings"
	"code.emcdtech.com/emcd/service/referral/protocol/default_whitelabel_settings"
	"code.emcdtech.com/emcd/service/referral/protocol/referral"
	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
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

	"code.emcdtech.com/emcd/service/profile/internal/config"
	"code.emcdtech.com/emcd/service/profile/internal/consumer"
	"code.emcdtech.com/emcd/service/profile/internal/healthchecker"
	"code.emcdtech.com/emcd/service/profile/internal/jobs"
	"code.emcdtech.com/emcd/service/profile/internal/notification"
	"code.emcdtech.com/emcd/service/profile/internal/pkg/aes"
	"code.emcdtech.com/emcd/service/profile/internal/repository"
	"code.emcdtech.com/emcd/service/profile/internal/server"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	pb "code.emcdtech.com/emcd/service/profile/protocol/profile"
)

const (
	shutdownTimeout                  = 5 * time.Second
	cacheUpdateCoinValidatorInterval = 15 * time.Minute
)

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    string = "profile"
	serviceVersion string = "local-build"
)

func main() {
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.CreateClientNamedContext(ctx, serviceName)

	log.Info(ctx, "starting profile server app")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Info(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	cfg, err := sdkConfig.New[config.Config]()
	if err != nil {
		// nolint: gocritic
		log.Fatal(ctx, "new config: %s", err.Error())
		return
	}

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
		WithMetrics(deps.Infrastructure.HTTPSrv).
		WithGRPC(
			func(opts ...grpc.ServerOption) *grpc.Server {
				grpcSrv = deps.Infrastructure.GRPCSrv(opts...)
				return grpcSrv
			},
			healthchecker.NewCommon(deps.Infrastructure.DB),
			healthchecker.NewCommon(deps.Infrastructure.EmcdDB),
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
		log.Fatal(ctx, "abnormal stopped profile server app: %s", err.Error())
	}
	log.Info(ctx, "graceful stopped profile server app")
}

// Deps dependencies
type Deps struct {
	Infrastructure struct {
		Tracer           trace.TracerProvider
		HTTPSrv          *echo.Echo
		GRPCSrv          app.GRPCServerFactory
		DB               *pgxpool.Pool
		EmcdDB           *pgxpool.Pool
		DBTransactor     pgTx.PgxTransactor
		EmcdDBTransactor pgTx.PgxTransactor
	}
	GRPCClient struct {
		Whitelabel              whitelabel.WhitelabelServiceClient
		Referral                referral.ReferralServiceClient
		Email                   emailPb.EmailServiceClient
		DefaultReferralSettings default_settings.DefaultSettingsServiceClient
		DefaultWLSettings       default_whitelabel_settings.DefaultWhitelabelSettingsServiceClient
		Coin                    coinPb.CoinServiceClient
		UserAccount             userAccountPb.UserAccountServiceClient
		Wallet                  walletPb.WalletServiceClient
	}
	Repository struct {
		Profile              repository.Profile
		OldUsers             repository.OldUsers
		UserLogs             repository.UserLogs
		KYC                  repository.Kyc
		NotificationSettings repository.NotificationSettings
		ProfileLog           repository.ProfileLog
		Whitelabel           repository.Whitelabel
		Referral             repository.Referral
		EMail                repository.Email
		Coin                 repository.Coin
		coinValidator        coinValidatorRepo.CoinValidatorRepository
		Wallet               repository.Wallet
	}
	Service struct {
		ProfileAdminLog service.ProfileLog
		Profile         service.Profile
	}
	Workers    []app.Worker
	Collectors []prometheus.Collector
}

func NewDeps(ctx context.Context, cfg config.Config) (Deps, error) {
	aes.Init(cfg.APIKeyConfig.Salt, cfg.APIKeyConfig.Secret) // ???
	var deps Deps
	if err := deps.makeInfrastructure(ctx, cfg); err != nil {
		return deps, fmt.Errorf("configurer infrastructure deps: %w", err)
	}
	if err := deps.makeGRPCClients(cfg); err != nil {
		return deps, fmt.Errorf("configurer grpc clients deps: %w", err)
	}
	deps.makeRepository()
	deps.makeService(cfg)
	deps.makeWorker()

	deps.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		pbProfileController := server.NewProfile(
			deps.Service.Profile,
			deps.Repository.OldUsers,
		)
		pbAdminController := server.NewAdmin(
			deps.Service.Profile,
			deps.Service.ProfileAdminLog,
		)

		opts = append(opts,
			grpc.ChainUnaryInterceptor(
				sdkErr.ServerUnaryInterceptor,
				log.ServerUnaryNamedInterceptor(serviceName),
			),
		)
		grpcSrv := grpc.NewServer(opts...)
		pb.RegisterProfileServiceServer(grpcSrv, pbProfileController)
		pb.RegisterAdminProfileServiceServer(grpcSrv, pbAdminController)
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
		return fmt.Errorf("new pgxpool internal db: %w", err)
	}
	d.Infrastructure.EmcdDB, err = sdkConfig.NewPGXPoolFromDSN(ctx, d.Infrastructure.Tracer, cfg.EmcdPostgresConnectionString)
	if err != nil {
		return fmt.Errorf("new pgxpool emcd db: %w", err)
	}

	d.Infrastructure.DBTransactor = pgTx.NewPgxTransactor(d.Infrastructure.DB)
	d.Infrastructure.EmcdDBTransactor = pgTx.NewPgxTransactor(d.Infrastructure.EmcdDB)

	d.Infrastructure.HTTPSrv = echo.New()
	d.Infrastructure.HTTPSrv.HideBanner = true
	d.Infrastructure.HTTPSrv.HidePort = true
	return nil
}

func (d *Deps) makeRepository() {
	d.Repository.Profile = repository.NewProfile(d.Infrastructure.DBTransactor, d.Infrastructure.EmcdDBTransactor)
	d.Repository.OldUsers = repository.NewOldUsers(d.Infrastructure.EmcdDBTransactor)
	d.Repository.UserLogs = repository.NewUserLogs(d.Infrastructure.EmcdDBTransactor)
	d.Repository.KYC = repository.NewKyc(d.Infrastructure.EmcdDBTransactor)
	d.Repository.NotificationSettings = repository.NewNotificationSettings(d.Infrastructure.DBTransactor)
	d.Repository.ProfileLog = repository.NewProfileLog(d.Infrastructure.DBTransactor)
	d.Repository.Whitelabel = repository.NewWhitelabel(d.GRPCClient.Whitelabel)
	d.Repository.Referral = repository.NewReferral(
		d.GRPCClient.Referral,
		d.GRPCClient.DefaultReferralSettings,
		d.GRPCClient.DefaultWLSettings,
	)
	d.Repository.EMail = repository.NewEmail(d.GRPCClient.Email)
	d.Repository.Coin = repository.NewCoin(d.GRPCClient.Coin)
	d.Repository.coinValidator = coinValidatorRepo.NewCoinValidatorRepository(d.GRPCClient.Coin)
	d.Repository.Wallet = repository.NewWallet(d.GRPCClient.Wallet)
}

func (d *Deps) makeService(cfg config.Config) {
	d.Service.ProfileAdminLog = service.NewProfileLog(d.Repository.ProfileLog)

	// пока тут временно оставил чтоб сократить время до получения результата
	analyticsNotifier := notification.NewAnalytics(cfg.SegmentKey)
	jobsClient := jobs.NewJobsClient(cfg.APIJobsConfig)

	// helper, via external repository
	userAccountRepo := userAccountRepository.NewUserAccountRepository(d.GRPCClient.UserAccount, d.Repository.coinValidator)
	userAccountService := service.NewUserAccountService(userAccountRepo, d.Repository.coinValidator)

	nonceStore := nonce.NewStore()

	d.Service.Profile = service.NewProfile(
		d.Repository.Profile,
		d.Repository.OldUsers,
		d.Repository.UserLogs,
		d.Repository.NotificationSettings,
		d.Repository.KYC,
		cfg.AccessSecret,
		analyticsNotifier,
		cfg.IdenfyRetryDelayMinutes,
		d.Repository.Referral,
		d.Repository.EMail,
		cfg.MinGetAllReferralsTake,
		jobsClient,
		userAccountService,
		d.Repository.Wallet,
		d.Repository.Coin,
		d.Repository.Whitelabel,
		nonceStore,
	)
}

func (d *Deps) makeGRPCClients(cfg config.Config) error {
	commonDialOption := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(d.Infrastructure.Tracer),
			),
		),
		grpc.WithChainUnaryInterceptor(
			log.ClientUnaryInterceptor,
			sdkErr.ClientUnaryInterceptor,
		),
	}

	cc, err := grpc.NewClient(cfg.GRPCClient.WhitelabelAddress, commonDialOption...)
	if err != nil {
		return fmt.Errorf("dial referral client: %w", err)
	}
	d.GRPCClient.Whitelabel = whitelabel.NewWhitelabelServiceClient(cc)

	cc, err = grpc.NewClient(cfg.GRPCClient.ReferralAddress, commonDialOption...)
	if err != nil {
		return fmt.Errorf("dial referral client: %w", err)
	}
	d.GRPCClient.Referral = referral.NewReferralServiceClient(cc)
	d.GRPCClient.DefaultReferralSettings = default_settings.NewDefaultSettingsServiceClient(cc)
	d.GRPCClient.DefaultWLSettings = default_whitelabel_settings.NewDefaultWhitelabelSettingsServiceClient(cc)

	cc, err = grpc.NewClient(cfg.GRPCClient.EmailAddress, commonDialOption...)
	if err != nil {
		return fmt.Errorf("dial email client: %w", err)
	}
	d.GRPCClient.Email = emailPb.NewEmailServiceClient(cc)

	cc, err = grpc.NewClient(cfg.GRPCClient.CoinAddress, commonDialOption...)
	if err != nil {
		return fmt.Errorf("dial coin client: %w", err)
	}
	d.GRPCClient.Coin = coinPb.NewCoinServiceClient(cc)

	cc, err = grpc.NewClient(cfg.GRPCClient.AccountingAddress, commonDialOption...)
	if err != nil {
		return fmt.Errorf("dial coin client: %w", err)
	}
	d.GRPCClient.UserAccount = userAccountPb.NewUserAccountServiceClient(cc)

	wc, err := grpc.NewClient(cfg.GRPCClient.WalletAddress, commonDialOption...)
	if err != nil {
		return fmt.Errorf("dial wallet client: %w", err)
	}
	d.GRPCClient.Wallet = walletPb.NewWalletServiceClient(wc)

	return nil
}

func (d *Deps) makeWorker() {
	d.Workers = app.Workers{
		app.WorkerFn(
			func(ctx context.Context) error {
				consumer.
					NewFinderNewUsers(
						d.Infrastructure.EmcdDBTransactor,
						d.GRPCClient.Whitelabel,
						d.Service.Profile,
					).
					Consume(ctx)
				return nil
			}),
	}

	coinValidatorWorker := app.WorkerFn(
		func(ctx context.Context) error {
			wg := new(sync.WaitGroup) // blank using
			defer wg.Wait()

			wg.Add(1)
			d.Repository.coinValidator.Serve(ctx, wg, cacheUpdateCoinValidatorInterval)

			return nil
		})

	d.Workers = append(d.Workers, coinValidatorWorker)
}
