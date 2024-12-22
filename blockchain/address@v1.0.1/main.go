package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	addressNodePb "code.emcdtech.com/emcd/blockchain/node/proto/address_node"
	nodeRepository "code.emcdtech.com/emcd/blockchain/node/repository_external"
	"code.emcdtech.com/emcd/sdk/app"
	sdkCfg "code.emcdtech.com/emcd/sdk/config"
	sdkErr "code.emcdtech.com/emcd/sdk/error"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	userAccountRepository "code.emcdtech.com/emcd/service/accounting/repository"
	coinPb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	profilePb "code.emcdtech.com/emcd/service/profile/protocol/profile"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"code.emcdtech.com/emcd/blockchain/address/internal/config"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/healthchecker"
	"code.emcdtech.com/emcd/blockchain/address/internal/metrics"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository/repository_migration"
	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	"code.emcdtech.com/emcd/blockchain/address/internal/worker"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

const shutdownTimeout = 2 * time.Second
const cacheUpdateCoinValidatorInterval = 15 * time.Minute
const migrationBatchInterval = 10 * time.Second
const migrationBatchIntervalPersonal = 1 * time.Minute

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    = "address"
	serviceVersion = "local-build"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = sdkLog.CreateClientNamedContext(ctx, serviceName)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		sdkLog.Info(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	cfg, err := sdkCfg.New[config.Config]()
	if err != nil {

		sdkLog.Panic(ctx, "failed load config: %v", err)

	}

	if err := cfg.Validate(); err != nil {

		sdkLog.Panic(ctx, "failed load config: %v", err)

	}

	if err = sdkLog.Init(ctx); err != nil {

		sdkLog.Panic(ctx, "failed init logger: %v", err)

	}

	tracerProvider, err := registerTracer(ctx, &cfg.Tracing, cfg.Environment.Name)
	if err != nil {

		sdkLog.Panic(ctx, "failed register tracer: %v", err)

	}

	dbPool, err := createPgPoolWithTracer(ctx, cfg.PGXPool.ConnectionString, tracerProvider)
	if err != nil {

		sdkLog.Panic(ctx, "failed create pgpool: %v", err)

	}
	defer func() {
		sdkLog.Info(context.Background(), "closing db pool")
		dbPool.Close()

	}()

	dbPoolMigrate, err := createPgPoolWithTracer(ctx, cfg.MigrationPostgresConnectionString, tracerProvider)
	if err != nil {

		sdkLog.Panic(ctx, "failed create migration pgpool: %v", err)

	}
	defer func() {
		sdkLog.Info(context.Background(), "closing migration db pool")
		dbPoolMigrate.Close()

	}()

	echoServer := echo.New()
	echoServer.HideBanner = true
	echoServer.HidePort = false

	d := newDeps()

	if err := d.createConnections(&cfg, tracerProvider); err != nil {

		sdkLog.Panic(ctx, "failed create connections: %v", err)

	}
	defer d.closeConnections()

	collector := metrics.NewAppMetrics()

	coinValidator := coinValidatorRepo.NewCoinValidatorRepository(d.getCoinServiceClient())
	addressRepo := repository.NewAddressRepository(dbPool)

	addressNodeRepo := nodeRepository.NewAddressNodeRepository(d.getAddressNodeServiceClient())
	userAccountRepo := userAccountRepository.NewUserAccountRepository(d.getUserAccountServiceClient(), coinValidator)

	migrationRepo := repository_migration.NewMigrationRepository(dbPoolMigrate, dbPool)
	migrationWorker := worker.NewMigrateAddress(coinValidator, migrationRepo, addressRepo)

	masterKeyIdMap := map[nodeCommon.NetworkGroupEnum][]string{
		nodeCommon.EthNetworkGroupId:  cfg.EthMasterKeys,
		nodeCommon.AlphNetworkGroupId: cfg.AlphMasterKeys,
	}
	rabbitRepo := repository.NewRabbitRepository(cfg.RabbitmqUrl)
	defer func() {
		if err := rabbitRepo.Close(); err != nil {
			sdkLog.Error(context.Background(), "failed close rabbit repo: %v", err)

		}
	}()

	rabbitService := service.NewRabbitService(cfg.RabbitmqExchangeName, rabbitRepo)
	addressService := service.NewAddressService(
		addressRepo,
		userAccountRepo,
		addressNodeRepo,
		d.getProfileServiceClient(),
		masterKeyIdMap,
		true,
		rabbitService)

	addressHandler := handler.NewAddressHandler(addressService, coinValidator, d.getCoinServiceClient(), cfg.GetIsNetworkOldWayMap())

	// selfExternalRepo := repositoryExternal.NewAddressRepository(coinValidator, d.getSelfServiceClient())

	a := app.New().
		WithTracing(tracerProvider).
		WithPprof(echoServer).
		WithMetrics(echoServer, collector).
		WithGRPC(d.serverFactory(addressHandler), // TODO: refactor -> deps.grpcServer
			healthchecker.NewCommon(dbPool),
		).
		WithWorker(
			app.NewHTTPSrvWorker(echoServer, cfg.HTTP.ListenAddr, shutdownTimeout),
			app.WorkerFn(
				func(ctx context.Context) error {

					return app.NewGRPCSrvWorker(d.grpcServer, cfg.GRPC.ListenAddr).Run(ctx) // TODO: refactor created by deps.createGrpcServer
				},
			),
			app.WorkerFn(
				func(ctx context.Context) error {
					wgFake := new(sync.WaitGroup)
					defer wgFake.Wait()

					wgFake.Add(1)
					coinValidator.Serve(ctx, wgFake, cacheUpdateCoinValidatorInterval)

					return nil
				}),
			// migrationWorker.GetWorkerAddress(),
			// migrationWorker.GetWorkerUserAccount(),
			// worker.NewGenerateCpuTempMetric(collector.CpuTemp),
			// worker.NewGenerateHDFailuresMetric(collector.HDFailures),
			worker.NewCronWorker(migrationBatchIntervalPersonal, migrationWorker.DoPersonalAddress),
			worker.NewCronWorker(migrationBatchInterval, migrationWorker.DoAddress),
			worker.NewCronWorker(migrationBatchInterval, migrationWorker.DoUserAccount),
		)

	if err := a.Run(ctx); err != nil {
		sdkLog.Error(ctx, "abnormal app shutdown: %v", err)

	}

	sdkLog.Info(ctx, "main tail is start gracefully shutdown...")
}

func registerTracer(ctx context.Context, traceCfg *sdkCfg.APMTracing, serviceEnvironment string) (trace.TracerProvider, error) {
	if tracerTransport, err := transport.NewHTTPTransport(transport.HTTPTransportOptions(*traceCfg)); err != nil {

		return nil, fmt.Errorf("failed create transport for apm: %w", err)
	} else if tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:        serviceName,
		ServiceVersion:     serviceVersion,
		ServiceEnvironment: serviceEnvironment,
		Transport:          tracerTransport,
	}); err != nil {

		return nil, fmt.Errorf("failed create tracer: %w", err)
	} else if tracerProvider, err := apmotel.NewTracerProvider(apmotel.WithAPMTracer(tracer)); err != nil {

		return nil, fmt.Errorf("failed create tracer provider: %w", err)
	} else {
		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

		sdkLog.Info(ctx, "otel tracer successfully registered")

		return tracerProvider, nil
	}
}

func createPgPoolWithTracer(ctx context.Context, dsnString string, tracerProvider trace.TracerProvider) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(dsnString)
	if err != nil {

		return nil, fmt.Errorf("failed parse connection string: %w", err)
	} else {
		pgxCfg.ConnConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithIncludeQueryParameters(),
			otelpgx.WithDisableQuerySpanNamePrefix(),
			otelpgx.WithTracerProvider(tracerProvider),
		)
	}

	if pgPool, err := pgxpool.NewWithConfig(ctx, pgxCfg); err != nil {

		return nil, fmt.Errorf("failed create pool: %w", err)
	} else if err = pgPool.Ping(ctx); err != nil {

		return nil, fmt.Errorf("failed did ping: %w", err)
	} else {

		return pgPool, nil
	}
}

// Deps dependencies
type deps struct {
	selfConn       *grpc.ClientConn
	accountingConn *grpc.ClientConn
	profileConn    *grpc.ClientConn
	coinConn       *grpc.ClientConn
	nodeConn       *grpc.ClientConn

	grpcServer *grpc.Server // TODO: refactor
}

func newDeps() *deps {

	return &deps{
		selfConn:       nil,
		accountingConn: nil,
		profileConn:    nil,
		coinConn:       nil,
		nodeConn:       nil,
		grpcServer:     nil,
	}
}

func (d *deps) createConnections(cfg *config.Config, tracerProvider trace.TracerProvider) error {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(tracerProvider),
			),
		),
		grpc.WithChainUnaryInterceptor(
			sdkLog.ClientUnaryInterceptor,
			sdkErr.ClientUnaryInterceptor,
		),
	}

	if selfConn, err := grpc.NewClient(cfg.GRPC.ListenAddr, dialOptions...); err != nil {

		return fmt.Errorf("new grpc client self: %w", err)
	} else {
		d.selfConn = selfConn

	}

	if accountingConn, err := grpc.NewClient(cfg.GrpcClient.AccountingAddress, dialOptions...); err != nil {

		return fmt.Errorf("new grpc client accounting: %w", err)
	} else {
		d.accountingConn = accountingConn

	}

	if profileConn, err := grpc.NewClient(cfg.GrpcClient.ProfileAddress, dialOptions...); err != nil {

		return fmt.Errorf("new grpc client profile: %w", err)
	} else {
		d.profileConn = profileConn

	}

	if coinConn, err := grpc.NewClient(cfg.GrpcClient.CoinAddress, dialOptions...); err != nil {

		return fmt.Errorf("new grpc client coin: %w", err)
	} else {
		d.coinConn = coinConn

	}

	if nodeConn, err := grpc.NewClient(cfg.GrpcClient.NodeAddress, dialOptions...); err != nil {

		return fmt.Errorf("new grpc client node: %w", err)
	} else {
		d.nodeConn = nodeConn

	}

	return nil
}

func (d *deps) closeConnections() {
	ctx := context.Background() // internal context, because external may be already closed

	if err := d.selfConn.Close(); err != nil {
		sdkLog.Error(ctx, "failed close self connection: %s", err)

	}

	if err := d.accountingConn.Close(); err != nil {
		sdkLog.Error(ctx, "failed close accounting connection: %s", err)

	}

	if err := d.profileConn.Close(); err != nil {
		sdkLog.Error(ctx, "failed close profile connection: %s", err)

	}

	if err := d.coinConn.Close(); err != nil {
		sdkLog.Error(ctx, "failed close coin connection: %s", err)

	}

	if err := d.nodeConn.Close(); err != nil {
		sdkLog.Error(ctx, "failed close node connection: %s", err)

	}
}

func (d *deps) serverFactory(addressHandler *handler.AddressHandler) func(opts ...grpc.ServerOption) *grpc.Server {

	return func(opts ...grpc.ServerOption) *grpc.Server {
		grpOptions := []grpc.ServerOption{grpc.ChainUnaryInterceptor(
			sdkErr.ServerUnaryInterceptor,
			sdkLog.ServerUnaryNamedInterceptor(serviceName),
		)}

		grpOptions = append(grpOptions, opts...)

		srv := grpc.NewServer(grpOptions...)

		// coinValidator := coinValidatorRepo.NewCoinValidatorRepository(d.getCoinServiceClient())
		//
		// pb.RegisterAddressServiceServer(srv, handler.NewAddressHandler(coinValidator))

		addressPb.RegisterAddressServiceServer(srv, addressHandler)
		reflection.Register(srv)

		d.grpcServer = srv // TODO: refactor, make this function pure

		return srv

	}
}

func (d *deps) getUserAccountServiceClient() userAccountPb.UserAccountServiceClient {

	return userAccountPb.NewUserAccountServiceClient(d.coinConn)
}

func (d *deps) getProfileServiceClient() profilePb.ProfileServiceClient {

	return profilePb.NewProfileServiceClient(d.profileConn)
}

func (d *deps) getCoinServiceClient() coinPb.CoinServiceClient {

	return coinPb.NewCoinServiceClient(d.coinConn)
}

func (d *deps) getAddressNodeServiceClient() addressNodePb.AddressNodeServiceClient {

	return addressNodePb.NewAddressNodeServiceClient(d.nodeConn)
}

// func (d *deps) getSelfServiceClient() addressPb.AddressServiceClient {
//
// 	return addressPb.NewAddressServiceClient(d.selfConn)
// }
