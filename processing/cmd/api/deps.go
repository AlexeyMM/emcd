package main

import (
	"context"
	"fmt"

	addresspb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/sdk/app"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	coinspb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	goGrpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"code.emcdtech.com/b2b/processing/internal/client"
	coinClient "code.emcdtech.com/b2b/processing/internal/client/coin"
	"code.emcdtech.com/b2b/processing/internal/config"
	"code.emcdtech.com/b2b/processing/internal/controller/admin"
	"code.emcdtech.com/b2b/processing/internal/controller/buyer"
	"code.emcdtech.com/b2b/processing/internal/controller/coin"
	"code.emcdtech.com/b2b/processing/internal/controller/merchant"
	"code.emcdtech.com/b2b/processing/internal/metrics"
	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/internal/repository/inmemory"
	repositoryMetrics "code.emcdtech.com/b2b/processing/internal/repository/metrics"
	"code.emcdtech.com/b2b/processing/internal/repository/pg"
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/internal/service/address"
	adminService "code.emcdtech.com/b2b/processing/internal/service/admin"
	coinService "code.emcdtech.com/b2b/processing/internal/service/coin"
	invoiceService "code.emcdtech.com/b2b/processing/internal/service/invoice"
	"code.emcdtech.com/b2b/processing/internal/worker"
)

type dependencies struct {
	*app.Deps
	cfg        config.APIConfig
	repository struct {
		merchantAdmin      repository.MerchantAdmin
		merchant           repository.Merchant
		depositAddressPool repository.DepositAddressPool
		invoice            repository.Invoice
		transaction        repository.Transaction
		coin               repository.Coin
	}
	service struct {
		merchantAdmin *adminService.MerchantService
		addressPool   *address.PoolService
		invoice       *invoiceService.Service
		coin          service.CoinService
	}
	controllers struct {
		buyer    *buyer.Controller
		admin    *admin.Controller
		merchant *merchant.Controller
		coin     *coin.Controller
	}
	serviceConfig struct {
		merchantInvoiceService invoiceService.ServiceConfig
	}
	grpcClient struct {
		addressService addresspb.AddressServiceClient
		coinService    coinspb.CoinServiceClient
	}
	workers struct {
		coinFetcher *worker.CoinFetcher
	}
	clients struct {
		coin client.CoinClient
	}
}

func newDeps(ctx context.Context, cfg config.APIConfig) (*dependencies, error) {
	appDeps, err := app.NewDeps(ctx, cfg.DepsConfig, serviceName, serviceVersion)
	if err != nil {
		return nil, fmt.Errorf("newDeps: %w", err)
	}

	deps := &dependencies{
		Deps: &appDeps,
		cfg:  cfg,
	}

	if err := deps.setupGRPClients(); err != nil {
		return nil, err
	}

	collector := metrics.New()
	appDeps.Collectors = append(appDeps.Collectors, collector)

	deps.setupClients()
	deps.setupServiceConfig()
	deps.setupRepository(collector)
	deps.setupService()
	deps.setupController()
	deps.setupWorkers()

	return deps, nil
}

func (d *dependencies) HealthCheckers() []app.HealthChecker {
	return []app.HealthChecker{
		&d.Deps.Infrastructure.DB,
	}
}

func (d *dependencies) setupGRPClients() error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(otelgrpc.WithTracerProvider(d.Infrastructure.Tracer))),
		grpc.WithChainUnaryInterceptor(
			apmgrpc.NewUnaryClientInterceptor(),
			goGrpcPrometheus.UnaryClientInterceptor,
			sdkError.ClientUnaryInterceptor,
			log.ClientUnaryInterceptor,
		),
	}

	addressConn, err := grpc.NewClient(d.cfg.GRPCClients.AddressServiceAddr, opts...)
	if err != nil {
		return err
	}

	d.grpcClient.addressService = addresspb.NewAddressServiceClient(addressConn)

	coinConn, err := grpc.NewClient(d.cfg.GRPCClients.CoinServiceAddr, opts...)
	if err != nil {
		return err
	}

	d.grpcClient.coinService = coinspb.NewCoinServiceClient(coinConn)

	return nil
}

func (d *dependencies) setupServiceConfig() {
	d.serviceConfig.merchantInvoiceService = invoiceService.ServiceConfig{InvoiceTTL: d.cfg.InvoiceService.InvoiceTTL}
}

func (d *dependencies) setupRepository(collector *metrics.AppMetrics) {
	d.repository.merchantAdmin = repositoryMetrics.NewMerchantAdmin(
		pg.NewMerchantAdmin(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
	d.repository.merchant = repositoryMetrics.NewMerchant(
		pg.NewMerchant(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
	d.repository.depositAddressPool = repositoryMetrics.NewDepositAddressPool(
		pg.NewDepositAddressPool(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
	d.repository.invoice = repositoryMetrics.NewInvoice(
		pg.NewInvoice(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
	d.repository.transaction = repositoryMetrics.NewTransaction(
		pg.NewTransaction(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
	d.repository.coin = inmemory.NewCoin()
}

func (d *dependencies) setupService() {
	d.service.merchantAdmin = adminService.NewMerchantService(d.repository.merchantAdmin)
	d.service.addressPool = address.NewPoolService(
		d.grpcClient.addressService,
		d.repository.depositAddressPool,
	)
	d.service.invoice = invoiceService.NewService(
		&invoiceService.ServiceConfig{
			InvoiceTTL: d.serviceConfig.merchantInvoiceService.InvoiceTTL,
		},
		d.repository.merchant,
		d.repository.invoice,
		d.service.addressPool,
		d.repository.transaction,
		d.repository.coin,
	)

	d.service.coin = coinService.NewCoinService(d.repository.coin, d.clients.coin, d.cfg.Coin.AvailableCoins)
}

func (d *dependencies) setupController() {
	// TODO: use actual constructors
	d.controllers.buyer = buyer.NewController(d.service.invoice)
	d.controllers.admin = admin.NewController(d.service.merchantAdmin)
	d.controllers.merchant = merchant.NewController(d.service.invoice)
	d.controllers.coin = coin.NewController(d.service.coin)
}

func (d *dependencies) setupWorkers() {
	d.workers.coinFetcher = worker.NewCoinFetcher(d.service.coin, d.cfg.Coin.FetchFrequency)
	d.Workers = append(d.Workers, d.workers.coinFetcher)
}

func (d *dependencies) Close() {
}

func (d *dependencies) setupClients() {
	d.clients.coin = coinClient.NewCoinClient(d.grpcClient.coinService)
}
