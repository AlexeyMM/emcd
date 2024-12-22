package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.emcdtech.com/emcd/sdk/app"
	cfgSdk "code.emcdtech.com/emcd/sdk/config"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/b2b/processing/internal/config"
)

const shutdownTimeout = 30 * time.Second

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    = "processing-coinwatch-consumer"
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

	cfg, err := cfgSdk.New[config.CoinwatchClientConfig]()
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
		deps.Close(ctx)
		log.Info(ctx, "dependencies have been gracefully shutdown.")
	}()

	err = app.New().
		WithDeps(deps.Deps).
		WithWorker(
			app.NewHTTPSrvWorker(
				deps.Infrastructure.HTTPSrv,
				deps.Config.HTTP.ListenAddr,
				deps.Config.HttpMetricShutdownTimeout,
			),
			deps.consumer.coinwatch,
		).
		Run(ctx)
	if err != nil {
		go func() {
			log.Fatal(ctx, err.Error())
		}()
	}

	log.Info(ctx, "main tail is start gracefully shutdown...")
}
