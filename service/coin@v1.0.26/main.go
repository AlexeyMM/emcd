package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	pbOldNode "code.emcdtech.com/emcd/blockchain/node/proto"
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

	"code.emcdtech.com/emcd/sdk/app"
	sdkConfig "code.emcdtech.com/emcd/sdk/config"
	sdkErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/coin/internal/config"
	"code.emcdtech.com/emcd/service/coin/internal/healthchecker"
	"code.emcdtech.com/emcd/service/coin/internal/repository"
	"code.emcdtech.com/emcd/service/coin/internal/repository/pg"
	"code.emcdtech.com/emcd/service/coin/internal/server"
	"code.emcdtech.com/emcd/service/coin/internal/service"
	pb "code.emcdtech.com/emcd/service/coin/protocol/coin"
)

const shutdownTimeout = 2 * time.Second

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    string = "coin"
	serviceVersion string = "local-build"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	log.Info(ctx, "starting coin server app")
	cfg, err := sdkConfig.New[config.Config]()
	if err != nil {
		log.Fatal(ctx, "get config: %s", err.Error())
		return
	}
	deps, err := NewDeps(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "build deps: %s", err.Error())
		return
	}

	otel.SetTracerProvider(deps.Infrastructure.Tracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	//if err = InitLogger(cfg.Log.Level); err != nil {
	//	log.Fatal().Err(err).Send()
	//}

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
		WithWorker(deps.Workers...).
		WithWorker(
			app.NewHTTPSrvWorker(deps.Infrastructure.HTTPSrv, cfg.HTTP.ListenAddr, shutdownTimeout),
			app.WorkerFn(
				func(ctx context.Context) error {
					return app.NewGRPCSrvWorker(grpcSrv, cfg.GRPC.ListenAddr).Run(ctx)
				},
			),
		).
		Run(ctx)
	if err != nil {
		log.Fatal(ctx, "coin service shutdown: %s", err.Error())
	}
	log.Info(ctx, "coin service gracefully shutdown...")
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
		WhiteLabel whitelabel.WhitelabelServiceClient
		OldNode    pbOldNode.NodeClient
	}
	Repository struct {
		WhiteLabel repository.WhiteLabel
		Coin       repository.Coin
	}
	Service struct {
		Coin *service.Coin
	}
	Workers    []app.Worker
	Collectors []prometheus.Collector
}

func NewDeps(ctx context.Context, cfg config.Config) (Deps, error) {
	var deps Deps
	if err := deps.makeInfrastructure(ctx, cfg); err != nil {
		return deps, fmt.Errorf("infrastructure: %w", err)
	}
	if err := deps.makeGRPCClient(cfg); err != nil {
		return deps, fmt.Errorf("grpc client: %w", err)
	}
	deps.makeRepository()
	deps.makeService()

	deps.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		pbCoinController := server.NewCoin(deps.Service.Coin)
		opts = append(opts, grpc.UnaryInterceptor(sdkErr.ServerUnaryInterceptor))
		grpcSrv := grpc.NewServer(opts...)
		pb.RegisterCoinServiceServer(grpcSrv, pbCoinController)
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
		return fmt.Errorf("pool to coin: %w", err)
	}

	d.Infrastructure.HTTPSrv = echo.New()
	d.Infrastructure.HTTPSrv.HideBanner = true
	d.Infrastructure.HTTPSrv.HidePort = true
	return nil
}

func (d *Deps) makeGRPCClient(cfg config.Config) error {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(d.Infrastructure.Tracer),
			),
		),
		grpc.WithChainUnaryInterceptor(
			sdkErr.ClientUnaryInterceptor,
		),
	}

	cc, err := grpc.NewClient(cfg.WhiteLabelAddress, dialOptions...)
	if err != nil {
		return fmt.Errorf("dial white lable client: %w", err)
	}
	d.GRPCClient.WhiteLabel = whitelabel.NewWhitelabelServiceClient(cc)

	cc, err = grpc.NewClient(cfg.NodeAddress, dialOptions...)
	if err != nil {
		return fmt.Errorf("dial white profile client: %w", err)
	}

	d.GRPCClient.OldNode = pbOldNode.NewNodeClient(cc)
	return nil
}

func (d *Deps) makeRepository() {
	d.Repository.Coin = pg.NewCoinStore(d.Infrastructure.DB)
	d.Repository.WhiteLabel = d.GRPCClient.WhiteLabel
}

func (d *Deps) makeService() {
	d.Service.Coin = service.NewCoin(
		d.Repository.Coin,
		d.GRPCClient.OldNode,
		d.GRPCClient.WhiteLabel,
	)
}
