package app

import (
	"context"
	"net/http"
	"net/http/pprof"

	sdkErr "code.emcdtech.com/emcd/sdk/error"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	noopMetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	noopTrace "go.opentelemetry.io/otel/trace/noop"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServerFactory = func(opts ...grpc.ServerOption) *grpc.Server

type Application struct {
	metricsRegistry  *prometheus.Registry
	getTraceProvider func() trace.TracerProvider
	getMeterProvider func() metric.MeterProvider
	applyGrpc        func(trace.TracerProvider, metric.MeterProvider)
	applyPprof       func()
	applyMetrics     func()
	applyHttp        func(provider trace.TracerProvider)
	workers          Workers
}

func New() *Application {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	app := &Application{
		metricsRegistry: registry,
		getTraceProvider: func() trace.TracerProvider {
			return noopTrace.NewTracerProvider()
		},
		getMeterProvider: func() metric.MeterProvider {
			return noopMetric.NewMeterProvider()
		},
		applyGrpc:    func(trace.TracerProvider, metric.MeterProvider) {},
		applyPprof:   func() {},
		applyMetrics: func() {},
		applyHttp:    func(trace trace.TracerProvider) {},
	}
	return app
}

func (a *Application) WithWorker(workers ...Worker) *Application {
	a.workers = append(a.workers, workers...)
	return a
}

func (a *Application) WithPprof(srv *echo.Echo) *Application {
	a.applyPprof = func() {
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		srv.GET("/debug/pprof/*", echo.WrapHandler(pprofMux))
	}
	return a
}

func (a *Application) WithHTTP(srv *echo.Echo, middleware ...echo.MiddlewareFunc) *Application {
	a.applyHttp = func(provider trace.TracerProvider) {
		srv.Use(otelecho.Middleware("", otelecho.WithTracerProvider(provider)))
		srv.Use(middleware...)
	}
	return a
}

func (a *Application) WithMetrics(srv *echo.Echo, collectors ...prometheus.Collector) *Application {
	a.applyMetrics = func() {
		a.metricsRegistry.MustRegister(collectors...)
		// may be a panic, watch https://github.com/labstack/echo-contrib/blob/master/echoprometheus/prometheus.go#L138
		srv.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Skipper: func(c echo.Context) bool {
				return c.Path() == "/metrics"
			},
			Registerer: a.metricsRegistry,
		}))
		srv.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{Gatherer: a.metricsRegistry}))
	}
	return a
}

func (a *Application) WithTracing(tracer trace.TracerProvider) *Application {
	a.getTraceProvider = func() trace.TracerProvider {
		return tracer
	}
	return a
}

// WithGRPC ...
// Deprecated: please use [WithGRPCServer] instead of WithGRPC and [NewGRPCSrvWorker] combo.
func (a *Application) WithGRPC(newServer GRPCServerFactory, healthCheckers ...HealthChecker) *Application {
	grpcServerMetrics := grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram())
	a.metricsRegistry.MustRegister(grpcServerMetrics)

	a.applyGrpc = func(tp trace.TracerProvider, mp metric.MeterProvider) {
		srv := newServer(
			grpc.ChainUnaryInterceptor(
				sdkErr.ServerUnaryInterceptor,
				grpcServerMetrics.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				grpcServerMetrics.StreamServerInterceptor(),
			),
			grpc.StatsHandler(otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tp),
				// otelgrpc.WithMeterProvider(mp),
			)),
		)
		grpc_health_v1.RegisterHealthServer(srv, NewGrpcHealth(healthCheckers))
	}
	return a
}

func (a *Application) WithDeps(deps *Deps) *Application {
	return a.WithTracing(deps.Infrastructure.Tracer).
		WithPprof(deps.Infrastructure.HTTPSrv).
		WithMetrics(deps.Infrastructure.HTTPSrv, deps.Collectors...)
}

func (a *Application) Run(ctx context.Context) error {
	traceProvider := a.getTraceProvider()
	metricProvider := a.getMeterProvider()
	a.applyPprof()
	a.applyMetrics()
	a.applyGrpc(traceProvider, metricProvider)
	a.applyHttp(traceProvider)
	if len(a.workers) > 0 {
		group, ctxGroup := errgroup.WithContext(ctx)
		for _, w := range a.workers {
			_w := w
			group.Go(func() error {
				return _w.Run(ctxGroup)
			})
		}
		return group.Wait()
	}
	<-ctx.Done()
	return nil
}
