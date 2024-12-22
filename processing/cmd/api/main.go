package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.emcdtech.com/b2b/processing/protocol/coinpb"

	"code.emcdtech.com/b2b/processing/internal/config"
	"code.emcdtech.com/b2b/processing/internal/controller"
	"code.emcdtech.com/b2b/processing/pkg/grpckit"
	"code.emcdtech.com/b2b/processing/protocol/adminpb"
	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
	"code.emcdtech.com/b2b/processing/protocol/merchantpb"

	"code.emcdtech.com/emcd/sdk/app"
	cfgSdk "code.emcdtech.com/emcd/sdk/config"
	"code.emcdtech.com/emcd/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const shutdownTimeout = 30 * time.Second

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    = "processing-grpc-api"
	serviceVersion = "local-build"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Info(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	cfg, err := cfgSdk.New[config.APIConfig]()
	if err != nil {
		log.Fatal(ctx, "new config: %s", err.Error())

		return
	}

	cfg.HttpMetricShutdownTimeout = shutdownTimeout

	deps, err := newDeps(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "new deps: %s", err.Error())

		return
	}

	defer func() {
		log.Info(ctx, "dependencies gracefully shutdown...")
		deps.Close()
		log.Info(ctx, "dependencies have been gracefully shutdown.")
	}()

	err = app.New().
		WithDeps(deps.Deps).
		WithGRPCServer(
			cfg.GRPC.ListenAddr,
			func(srv *grpc.Server) {
				buyerpb.RegisterInvoiceBuyerServiceServer(srv, deps.controllers.buyer)
				adminpb.RegisterMerchantAdminServiceServer(srv, deps.controllers.admin)
				merchantpb.RegisterInvoiceMerchantServiceServer(srv, deps.controllers.merchant)
				coinpb.RegisterCoinsServiceServer(srv, deps.controllers.coin)

				grpc_health_v1.RegisterHealthServer(srv, app.NewGrpcHealth(deps.HealthCheckers()))
			},
			// interceptors to convert business errors to grpc
			grpc.ChainUnaryInterceptor(
				grpckit.UnaryServerInterceptor(controller.NewErrorConverter()),
				grpckit.ValidationUnaryInterceptor(),
			),
			grpc.ChainStreamInterceptor(grpckit.StreamServerInterceptor(controller.NewErrorConverter())),
		).
		WithWorker(
			app.NewHTTPSrvWorker(
				deps.Infrastructure.HTTPSrv,
				deps.Config.HTTP.ListenAddr,
				deps.Config.HttpMetricShutdownTimeout,
			),
		).
		WithWorker(deps.Workers...).
		Run(ctx)
	if err != nil {
		go func() {
			log.Panic(ctx, err.Error())
		}()
	}

	log.Info(ctx, "main tail is start gracefully shutdown...")
}
