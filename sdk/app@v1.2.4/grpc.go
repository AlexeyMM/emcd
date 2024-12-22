package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
)

// WithGRPCServer creates new GRPC server with default options, registers your services with the registrar func
// and adds workers to the Application.
//
// NOTE: at least for now you should register your health-checkers with
//
//	grpc_health_v1.RegisterHealthServer(srv, app.NewGrpcHealth(healthCheckers))
//
// Full example:
//
//	app.New().
//		WithTracing(tracerProvider).
//		WithGRPCServer(cfg.GRPC.ListenAddr, func(srv *grpc.Server) {
//			yourproto.RegisterYourServiceServer(srv, yourServer)
//			grpc_health_v1.RegisterHealthServer(srv, app.NewGrpcHealth(yourHealthCheckers))
//		}).
//		Run(ctx)
func (a *Application) WithGRPCServer(
	address string,
	registrar func(srv *grpc.Server),
	opts ...grpc.ServerOption,
) *Application {
	grpcServerMetrics := grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram())
	a.metricsRegistry.MustRegister(grpcServerMetrics)

	defaultOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			sdkError.ServerUnaryInterceptor,
			grpcServerMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcServerMetrics.StreamServerInterceptor(),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(a.getTraceProvider()))), // FIXME could be noop if WithGRPCServer called before WithTracing
	}

	srv := grpc.NewServer(append(defaultOptions, opts...)...)
	reflection.Register(srv)

	registrar(srv)

	a.workers = append(a.workers,
		WorkerFn(func(ctx context.Context) error {
			<-ctx.Done()
			srv.GracefulStop()
			return nil
		}),
		WorkerFn(func(ctx context.Context) error {
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return fmt.Errorf("open listener %s: %w", address, err)
			}

			log.SInfo(ctx, "grpc server started", map[string]any{"address": listener.Addr().String()})

			err = srv.Serve(listener)
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			} else if err != nil {
				return fmt.Errorf("grpc server error: %w", err)
			}

			log.Info(ctx, "grpc server shut down gracefully")

			return nil
		}),
	)

	return a
}
