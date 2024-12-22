package main

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/sdk/app"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	accountingpb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	useraccountpb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	accounting "code.emcdtech.com/emcd/service/accounting/repository"
	"code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinrepository "code.emcdtech.com/emcd/service/coin/repository"
	profilepb "code.emcdtech.com/emcd/service/profile/protocol/profile"
	goGrpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"code.emcdtech.com/b2b/processing/internal/config"
	"code.emcdtech.com/b2b/processing/internal/consumer/coinwatch"
	"code.emcdtech.com/b2b/processing/internal/metrics"
	"code.emcdtech.com/b2b/processing/internal/repository"
	repositoryMetrics "code.emcdtech.com/b2b/processing/internal/repository/metrics"
	"code.emcdtech.com/b2b/processing/internal/repository/pg"
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/internal/service/fee"
	"code.emcdtech.com/b2b/processing/internal/service/transaction"
	"code.emcdtech.com/b2b/processing/pkg/rabbitmqkit"
)

type closer struct {
	name  string
	close func() error
}

type dependencies struct {
	*app.Deps
	cfg      config.CoinwatchClientConfig
	rabbitmq struct {
		channel *amqp.Channel
	}
	repository struct {
		invoice     repository.Invoice
		transaction repository.Transaction
	}
	service struct {
		transaction service.Transaction
		fee         service.Fee
	}
	consumer struct {
		coinwatch *coinwatch.Consumer
	}
	grpcClient struct {
		profileService     profilepb.ProfileServiceClient
		userAccountService useraccountpb.UserAccountServiceClient
		coinService        coin.CoinServiceClient
		accountingService  accountingpb.AccountingServiceClient
	}
	externalServiceRepositories struct {
		accounting struct {
			userAccount accounting.UserAccountRepository
			accounting  accounting.AccountingRepository // operations with balance, etc
		}
		coin coinrepository.CoinValidatorRepository
	}
	closers []closer
}

func newDeps(ctx context.Context, cfg config.CoinwatchClientConfig) (*dependencies, error) {
	appDeps, err := app.NewDeps(ctx, cfg.DepsConfig, serviceName, serviceVersion)
	if err != nil {
		return nil, fmt.Errorf("newDeps: %w", err)
	}

	deps := &dependencies{
		Deps: &appDeps,
		cfg:  cfg,
	}

	if err := deps.setupRabbitMQ(ctx); err != nil {
		return nil, fmt.Errorf("setupRabbitMQ: %w", err)
	}

	if err := deps.setupGRPCClients(); err != nil {
		return nil, fmt.Errorf("setupGRPCClients: %w", err)
	}

	collector := metrics.New()
	appDeps.Collectors = append(appDeps.Collectors, collector)

	deps.setupRepository(collector)
	deps.setupExternalServiceRepository()
	deps.setupService(collector)
	deps.setupConsumer()

	return deps, nil
}

func (d *dependencies) HealthCheckers() []app.HealthChecker {
	return []app.HealthChecker{
		&d.Deps.Infrastructure.DB,
		rabbitmqkit.NewChannelHealthChecker(d.rabbitmq.channel, "coinwatch-rabbitmq"),
	}
}

func (d *dependencies) setupRepository(collector *metrics.AppMetrics) {
	d.repository.invoice = repositoryMetrics.NewInvoice(
		pg.NewInvoice(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
	d.repository.transaction = repositoryMetrics.NewTransaction(
		pg.NewTransaction(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
}

func (d *dependencies) setupExternalServiceRepository() {
	d.externalServiceRepositories.coin = coinrepository.NewCoinValidatorRepository(d.grpcClient.coinService)
	d.externalServiceRepositories.accounting.accounting = accounting.NewAccountingRepository(
		d.grpcClient.accountingService,
	)
	d.externalServiceRepositories.accounting.userAccount = accounting.NewUserAccountRepository(
		d.grpcClient.userAccountService,
		d.externalServiceRepositories.coin,
	)
}

func (d *dependencies) setupService(collector *metrics.AppMetrics) {
	d.service.fee = fee.NewService(
		d.grpcClient.profileService,
		d.externalServiceRepositories.accounting.userAccount,
		d.externalServiceRepositories.accounting.accounting,
		d.cfg.FeeCollectorUserID,
	)

	d.service.transaction = transaction.NewService(
		d.repository.transaction,
		d.repository.invoice,
		d.service.fee,
		collector.InvoiceExecutionHistogram,
	)
}

func (d *dependencies) setupConsumer() {
	d.consumer.coinwatch = coinwatch.NewConsumer(
		d.rabbitmq.channel,
		d.cfg.CoinwatchProcessingExchange,
		d.service.transaction,
	)
}

func (d *dependencies) setupRabbitMQ(ctx context.Context) error {
	conn, err := amqp.Dial(d.cfg.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	d.rabbitmq.channel = channel

	closedChan := make(chan *amqp.Error, 1)
	d.rabbitmq.channel.NotifyClose(closedChan)

	go func() {
		err := <-closedChan
		log.SError(ctx, "rabbitmq channel was closed", map[string]any{"err": err})
	}()

	d.closers = append(d.closers, closer{name: "rabbitmq", close: conn.Close})

	return nil
}

func (d *dependencies) setupGRPCClients() error {
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

	// Profile service
	profileConn, err := grpc.NewClient(d.cfg.GRPCClients.ProfileServiceAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to profile service: %w", err)
	}

	d.grpcClient.profileService = profilepb.NewProfileServiceClient(profileConn)

	// User Account service
	userAccountConn, err := grpc.NewClient(d.cfg.GRPCClients.UserAccountServiceAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to user account service: %w", err)
	}

	d.grpcClient.userAccountService = useraccountpb.NewUserAccountServiceClient(userAccountConn)

	// Coin service
	coinConn, err := grpc.NewClient(d.cfg.GRPCClients.CoinServiceAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to coin service: %w", err)
	}

	d.grpcClient.coinService = coin.NewCoinServiceClient(coinConn)

	// Accounting service
	accountingConn, err := grpc.NewClient(d.cfg.GRPCClients.AccountingServiceAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to accounting service: %w", err)
	}

	d.grpcClient.accountingService = accountingpb.NewAccountingServiceClient(accountingConn)

	d.closers = append(
		d.closers,
		closer{name: "profileConn", close: profileConn.Close},
		closer{name: "coinConn", close: coinConn.Close},
		closer{name: "accountingConn", close: accountingConn.Close},
		closer{name: "userAccountConn", close: userAccountConn.Close},
	)

	return nil
}

func (d *dependencies) Close(ctx context.Context) {
	for _, c := range d.closers {
		if err := c.close(); err != nil {
			log.SError(ctx, "closer failed", map[string]any{"err": err, "closer": c.name})
		}
	}
}
