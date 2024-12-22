package app

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/config"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	pg "code.emcdtech.com/emcd/sdk/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type DepsConfig struct {
	Environment               config.Environment
	Log                       config.Log
	GRPC                      config.GRPCServer
	Tracing                   config.APMTracing
	PGXPool                   pg.PGXPool
	HTTP                      config.HTTPServer
	HttpMetricShutdownTimeout time.Duration `env:"HTTP_METRIC_SHUTDOWN_TIMEOUT"`
}

type DB struct {
	*pgxpool.Pool
}

type Infrastructure struct {
	Tracer  trace.TracerProvider
	HTTPSrv *echo.Echo
	GRPCSrv GRPCServerFactory
	DB      DB
}

type Deps struct {
	Config         DepsConfig
	Infrastructure Infrastructure
	Workers        Workers
	Collectors     []prometheus.Collector
	DeferHandlers  []func(ctx context.Context)
}

func NewDeps(ctx context.Context, cfg DepsConfig, serviceName string, serviceVersion string) (Deps, error) {
	deps := Deps{
		Config: cfg,
	}

	if err := deps.makeInfrastructure(ctx, cfg, serviceName, serviceVersion); err != nil {
		return deps, fmt.Errorf("configurer infrastructure deps: %w", err)
	}

	deps.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpc.UnaryInterceptor(sdkError.ServerUnaryInterceptor))
		grpcSrv := grpc.NewServer(opts...)
		return grpcSrv
	}

	otel.SetTracerProvider(deps.Infrastructure.Tracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return deps, nil
}

func (d *Deps) makeInfrastructure(ctx context.Context, cfg DepsConfig, serviceName string, serviceVersion string) (err error) {
	tracer, err := cfg.Tracing.New(serviceName, serviceVersion, cfg.Environment.Name)
	if err != nil {
		return fmt.Errorf("new apm tracer: %w", err)
	}
	d.Infrastructure.Tracer, err = apmotel.NewTracerProvider(apmotel.WithAPMTracer(tracer))
	if err != nil {
		return fmt.Errorf("new tracer provider: %w", err)
	}

	pool, err := cfg.PGXPool.New(ctx, d.Infrastructure.Tracer)
	if err != nil {
		return fmt.Errorf("pool to referral statistic: %w", err)
	}
	d.Infrastructure.DB = DB{pool}

	d.Infrastructure.HTTPSrv = echo.New()
	d.Infrastructure.HTTPSrv.HideBanner = false
	d.Infrastructure.HTTPSrv.HidePort = false
	d.Infrastructure.HTTPSrv.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowHeaders:     []string{"*"},
			AllowMethods:     []string{"*"},
		}),
		middleware.RecoverWithConfig(
			middleware.RecoverConfig{
				LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
					log.Error(c.Request().Context(), "recover: %s, stack: %s", err.Error(), string(stack))
					return err
				},
			},
		),
	)
	return nil
}

func (d *DB) Check(ctx context.Context) HealthCheckServingStatus {
	if err := d.Ping(ctx); err != nil {
		return ServiceUnknownHealthCheckServingStatus
	}
	return ServingHealthCheckServingStatus
}

func (d *DB) ServiceName() string {
	return "pgxpool"
}
