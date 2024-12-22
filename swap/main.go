package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.emcdtech.com/emcd/sdk/app"
	cfgSdk "code.emcdtech.com/emcd/sdk/config"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/coin/protocol/coin"
	"code.emcdtech.com/emcd/service/email/protocol/email"
	goGrpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"code.emcdtech.com/b2b/swap/internal/client"
	clientMetrics "code.emcdtech.com/b2b/swap/internal/client/metrics"
	grpcController "code.emcdtech.com/b2b/swap/internal/controller"
	"code.emcdtech.com/b2b/swap/internal/metrics"
	"code.emcdtech.com/b2b/swap/internal/repository"
	repositoryMetrics "code.emcdtech.com/b2b/swap/internal/repository/metrics"
	"code.emcdtech.com/b2b/swap/internal/repository/pg"
	"code.emcdtech.com/b2b/swap/internal/repository/status_history"
	"code.emcdtech.com/b2b/swap/internal/slack"
	"code.emcdtech.com/b2b/swap/internal/worker"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/protocol/swap"
	"code.emcdtech.com/b2b/swap/protocol/swapAdmin"
	"code.emcdtech.com/b2b/swap/protocol/swapCoin"
	"code.emcdtech.com/b2b/swap/protocol/swapWithdraw"

	"code.emcdtech.com/b2b/swap/internal/client/bybit"
	"code.emcdtech.com/b2b/swap/internal/config"
	"code.emcdtech.com/b2b/swap/internal/repository/local"
	"code.emcdtech.com/b2b/swap/internal/service"
)

const shutdownTimeout = 2 * time.Second

var (
	serviceName    = "swap-service"
	serviceVersion = "local-build"
)

var wsHub *bybit.Hub

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Info(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	cfg, err := cfgSdk.New[config.Config]()
	if err != nil {
		log.Fatal(ctx, "new config: %s", err.Error())
		return
	}
	if cfg.ByBit.ApiKey == "" || cfg.ByBit.ApiSecret == "" {
		log.Fatal(ctx, "config error, please set api key and secret")
		return
	}

	collector := metrics.New()

	byBitCl := bybit.NewByBit(cfg.ByBit.MasterUid, cfg.ByBit.ApiUrl, cfg.ByBit.ApiKey, cfg.ByBit.ApiSecret)
	swapExecutorCh := make(chan *model.Swap)

	metricsByBit := clientMetrics.NewByBit(byBitCl, byBitCl, byBitCl, collector.ByBitRequestTimeHistogram)

	orderbookUpdateCh := make(chan model.OrderBookUpdateMessage, worker.OrderbookChanSize)
	brokenOrderBookWsCh := make(chan []*model.Symbol)

	deps, swapExecutor, err := NewDeps(ctx, metricsByBit, cfg, collector, swapExecutorCh, orderbookUpdateCh, brokenOrderBookWsCh)
	if err != nil {
		log.Fatal(ctx, "new deps: %s", err.Error())
		return
	}

	subscribeOnOrderbooks(ctx, &deps)

	go func() {
		// Время, подписаться на orderbooks прежде чем начать выполнять незавершённые свопы
		<-time.After(3 * time.Second)
		err = startSwaps(ctx, &deps, swapExecutor)
		if err != nil {
			log.Fatal(ctx, "start swaps: %s", err.Error())
		}
	}()

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
		WithMetrics(deps.Infrastructure.HTTPSrv, deps.Collectors...).
		WithGRPC(
			func(opts ...grpc.ServerOption) *grpc.Server {
				grpcSrv = deps.Infrastructure.GRPCSrv(opts...)
				return grpcSrv
			},
			//healthchecker.NewCommon(deps.Infrastructure.DB),
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
		WithWorker(
			app.WorkerFn(func(ctx context.Context) error {
				wsHub.Run(ctx)
				return nil
			})).
		Run(ctx)
	if err != nil {
		log.Fatal(ctx, "app.Run: %s", err.Error())
	}
}

type Deps struct {
	Infrastructure struct {
		Tracer  trace.TracerProvider
		HTTPSrv *echo.Echo
		GRPCSrv app.GRPCServerFactory
		DB      *pgxpool.Pool
	}
	Repository struct {
		Account           repository.Account
		ActiveSwap        repository.ActiveSwap
		Coin              repository.Coin
		ErrorCounter      repository.ErrorCounter
		Deposit           repository.Deposit
		Explorer          repository.Explorer
		OrderFee          repository.OrderFee
		Order             repository.Order
		Orderbook         repository.OrderBook
		StatusCache       repository.StatusCache
		Subscribers       repository.Subscribers
		Swap              repository.Swap
		SwapStatusHistory repository.SwapStatusHistory
		Symbol            repository.Symbol
		Transfer          repository.Transfer
		User              repository.User
		Withdraw          repository.Withdraw
		Slack             slack.Slack
	}
	ExternalClient struct {
		ExchangeAccount     client.ExchangeAccount
		ExchangeTransaction client.ExchangeTransaction
		Market              client.Market
		Subscriber          client.Subscriber
	}
	GrpcClients struct {
		Email email.EmailServiceClient
		Coin  coin.CoinServiceClient
	}
	Client struct {
		Email client.Email
		Coin  client.Coin
	}
	Service struct {
		Admin                service.Admin
		Coin                 service.Coin
		Order                service.Order
		OrderFee             service.OrderFee
		OrderBook            service.OrderBook
		Swap                 service.Swap
		SwapStatusUpdater    service.SwapStatusUpdater
		SwapStatusSubscriber service.SwapStatusSubscriber
		Symbol               service.Symbol
		Transfer             service.Transfer
		Withdraw             service.Withdraw
	}
	//GRPCClients struct {
	//	Skeleton pb.SkeletonServiceClient
	//}
	Controller struct {
		Admin    *grpcController.Admin
		Coin     *grpcController.Coin
		Swap     *grpcController.Swap
		Withdraw *grpcController.Withdraw
	}
	Workers       app.Workers
	Collectors    []prometheus.Collector
	DeferHandlers []func(ctx context.Context)
}

func NewDeps(
	ctx context.Context,
	metricsByBit *clientMetrics.ByBit,
	cfg config.Config,
	collector *metrics.AppMetrics,

	swapExecutorCh chan *model.Swap,
	orderbookUpdateCh chan model.OrderBookUpdateMessage,
	brokenOrderBookWsCh chan []*model.Symbol,
) (Deps, *worker.SwapExecutor, error) {
	var deps Deps
	if err := deps.makeInfrastructure(ctx, cfg, collector); err != nil {
		return deps, nil, fmt.Errorf("configurer infrastructure deps: %w", err)
	}
	deps.makeRepository(ctx, metricsByBit, cfg, collector)
	deps.makeClient(metricsByBit, cfg, collector, orderbookUpdateCh, brokenOrderBookWsCh)
	err := deps.makeGRPCClients(cfg)
	if err != nil {
		return deps, nil, fmt.Errorf("makeGRPCClients: %w", err)
	}
	deps.makeClients(cfg)
	err = deps.makeService(ctx, cfg, swapExecutorCh, metricsByBit)
	if err != nil {
		return deps, nil, fmt.Errorf("makeService: %w", err)
	}

	//if err := deps.makeGRPCClients(cfg); err != nil {
	//	return deps, fmt.Errorf("configurer clients deps: %w", err)
	//}

	deps.makeController()

	swapExecutor, err := deps.makeWorkers(ctx, cfg, collector, swapExecutorCh, orderbookUpdateCh, brokenOrderBookWsCh)
	if err != nil {
		return deps, nil, fmt.Errorf("makeWorkers: %w", err)
	}

	deps.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpc.UnaryInterceptor(sdkError.ServerUnaryInterceptor))
		grpcSrv := grpc.NewServer(opts...)
		swapCoin.RegisterSwapCoinServiceServer(grpcSrv, deps.Controller.Coin)
		swap.RegisterSwapServiceServer(grpcSrv, deps.Controller.Swap)
		swapWithdraw.RegisterSwapWithdrawServiceServer(grpcSrv, deps.Controller.Withdraw)
		swapAdmin.RegisterAdminServiceServer(grpcSrv, deps.Controller.Admin)
		reflection.Register(grpcSrv)
		return grpcSrv
	}
	return deps, swapExecutor, nil
}

func (d *Deps) makeInfrastructure(ctx context.Context, cfg config.Config, collector *metrics.AppMetrics) (err error) {
	tracer, err := cfg.Tracing.New(serviceName, serviceVersion, cfg.Environment.Name)
	if err != nil {
		return fmt.Errorf("new apm tracer: %w", err)
	}
	d.Infrastructure.Tracer, err = apmotel.NewTracerProvider(apmotel.WithAPMTracer(tracer))
	if err != nil {
		return fmt.Errorf("new tracer provider: %w", err)
	}
	otel.SetTracerProvider(d.Infrastructure.Tracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	d.Infrastructure.DB, err = cfg.PGXPool.New(ctx, d.Infrastructure.Tracer)
	if err != nil {
		return fmt.Errorf("pool to referral statistic: %w", err)
	}

	d.Collectors = append(d.Collectors, collector)

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

func (d *Deps) makeRepository(ctx context.Context, metricsByBit *clientMetrics.ByBit, cfg config.Config, collector *metrics.AppMetrics) {
	d.Repository.Account = repositoryMetrics.NewAccount(pg.NewAccount(d.Infrastructure.DB), collector.RepositoryRequestTimeHistogram)
	d.Repository.ActiveSwap = local.NewActiveSwap(cfg.SwapExecutorWorkerGroup)
	d.Repository.Coin = local.NewCoins()
	d.Repository.ErrorCounter = local.NewErrorCounter()
	d.Repository.Deposit = repositoryMetrics.NewDeposit(pg.NewDeposit(d.Infrastructure.DB), collector.RepositoryRequestTimeHistogram)
	d.Repository.Explorer = local.NewExplorer()

	d.Repository.OrderFee = local.NewFee()
	feeMap, err := metricsByBit.GetAllFeeRate(ctx)
	if err != nil {
		log.Fatal(ctx, "byBitCl.GetFeeRate: %s", err.Error())
		return
	}
	err = d.Repository.OrderFee.UpdateAll(ctx, feeMap)
	if err != nil {
		log.Fatal(ctx, "orderFee.UpdateAll: %s", err.Error())
		return
	}

	d.Repository.Order = repositoryMetrics.NewOrder(pg.NewOrder(d.Infrastructure.DB), collector.RepositoryRequestTimeHistogram)
	d.Repository.Orderbook = local.NewOrderBookStore()
	d.Repository.StatusCache = local.NewStatusCache()
	d.Repository.Subscribers = local.NewSubscribers()
	d.Repository.SwapStatusHistory = pg.NewSwapStatusHistory(d.Infrastructure.DB)
	d.Repository.Swap = repositoryMetrics.NewSwap(
		status_history.NewHistory(pg.NewSwap(d.Infrastructure.DB), d.Repository.SwapStatusHistory),
		collector.RepositoryRequestTimeHistogram,
	)

	symbols, err := metricsByBit.GetInstrumentsInfo(ctx)
	if err != nil {
		log.Fatal(ctx, "byBitCl.GetInstrumentsInfo: %s", err.Error())
		return
	}
	for _, sym := range symbols {
		sym.Accuracy = &model.Accuracy{
			BaseAccuracy:  service.CountDecimalPlaces(sym.BasePrecision.String()),
			QuoteAccuracy: service.CountDecimalPlaces(sym.QuotePrecision.String()),
		}
	}
	d.Repository.Symbol = local.NewSymbol(symbols)

	d.Repository.Transfer = repositoryMetrics.NewInternalTransfer(pg.NewTransfer(d.Infrastructure.DB), collector.RepositoryRequestTimeHistogram)
	d.Repository.User = repositoryMetrics.NewUser(pg.NewUser(d.Infrastructure.DB), collector.RepositoryRequestTimeHistogram)
	d.Repository.Withdraw = repositoryMetrics.NewWithdraw(pg.NewWithdraw(d.Infrastructure.DB), collector.RepositoryRequestTimeHistogram)
	d.Repository.Slack = slack.NewSlack(cfg.SlackWebhookUrl)
}

func (d *Deps) makeClient(metricsByBit *clientMetrics.ByBit, cfg config.Config, collector *metrics.AppMetrics,
	updateCh chan model.OrderBookUpdateMessage, brokenOrderBookWs chan []*model.Symbol) {
	d.Service.Order = service.NewOrder(d.Repository.Order)

	swapStatus := service.NewSwapStatus(d.Repository.Swap, d.Repository.SwapStatusHistory, d.Repository.Subscribers, d.Repository.StatusCache)
	d.Service.SwapStatusUpdater = swapStatus
	d.Service.SwapStatusSubscriber = swapStatus

	d.Service.Symbol = service.NewSymbol(d.Repository.Symbol, metricsByBit)

	wsHub = bybit.NewHub(d.Service.SwapStatusUpdater, d.Service.Order, cfg.ByBit.WsExpiredTime, &model.Secrets{
		ApiKey:    cfg.ByBit.ApiKey,
		ApiSecret: cfg.ByBit.ApiSecret,
	}, updateCh, brokenOrderBookWs, cfg.DepositWaitingPeriod, d.Repository.Slack, collector.ByBitOrderBookWebsocketGauge,
		collector.ByBitReconnectWebsocketCounter)

	d.ExternalClient.ExchangeAccount = metricsByBit
	d.ExternalClient.ExchangeTransaction = metricsByBit
	d.ExternalClient.Market = metricsByBit
	d.ExternalClient.Subscriber = wsHub

	d.Service.OrderFee = service.NewOrderFee(d.ExternalClient.Market, d.Repository.OrderFee)

	d.Service.OrderBook = service.NewOrderBook(d.Repository.Orderbook, d.ExternalClient.Subscriber)
}

func (d *Deps) makeGRPCClients(cfg config.Config) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(d.Infrastructure.Tracer),
			),
		),
		grpc.WithChainUnaryInterceptor(
			apmgrpc.NewUnaryClientInterceptor(),
			goGrpcPrometheus.UnaryClientInterceptor,
			sdkError.ClientUnaryInterceptor,
			log.ClientUnaryInterceptor,
		),
	}

	emailConn, err := grpc.NewClient(cfg.EmailAddress, opts...)
	if err != nil {
		return fmt.Errorf("new grpc client skeleton: %w", err)
	}
	d.GrpcClients.Email = email.NewEmailServiceClient(emailConn)

	coinConn, err := grpc.NewClient(cfg.CoinAddress, opts...)
	if err != nil {
		return fmt.Errorf("new grpc client skeleton: %w", err)
	}
	d.GrpcClients.Coin = coin.NewCoinServiceClient(coinConn)

	d.DeferHandlers = append(d.DeferHandlers,
		func(ctx context.Context) {
			_ = emailConn.Close()
			_ = coinConn.Close()
		},
	)
	return nil
}

func (d *Deps) makeClients(cfg config.Config) {
	d.Client.Email = client.NewEmail(d.GrpcClients.Email, cfg.Environment.Name)
	d.Client.Coin = client.NewCoin(d.GrpcClients.Coin)
}

func (d *Deps) makeService(ctx context.Context, cfg config.Config, swapExecutorCh chan *model.Swap, marketCl client.Market) error {
	d.Service.Admin = service.NewAdmin(
		d.ExternalClient.Market,
		d.ExternalClient.ExchangeAccount,
		d.ExternalClient.ExchangeTransaction,
		d.Client.Email,
		d.Repository.Swap,
		d.Repository.SwapStatusHistory,
		d.Repository.Transfer,
		d.Repository.Withdraw,
		d.Repository.Coin,
		d.Repository.Orderbook,
		d.Repository.Explorer,
		d.Repository.User,
		d.Service.SwapStatusUpdater,
		cfg.ByBit.MasterUid,
		cfg.ByBit.ApiKey,
		cfg.ByBit.ApiSecret,
	)
	d.Service.Coin = service.NewCoin(d.Repository.Coin, marketCl, d.Client.Coin)
	err := d.Service.Coin.SyncWithAPI(ctx)
	if err != nil {
		return err
	}
	d.Service.Withdraw = service.NewWithdraw(d.Repository.Withdraw, d.Repository.Explorer)
	d.Service.Transfer = service.NewTransfer(d.Repository.Transfer)

	// move in makeClient because of dependencies
	//d.Service.Order = service.NewOrder(d.Repository.Order)

	d.Service.Swap = service.NewSwap(
		d.ExternalClient.Market,
		d.ExternalClient.ExchangeAccount,
		d.ExternalClient.ExchangeTransaction,
		d.ExternalClient.Subscriber,
		d.Client.Email,
		d.Repository.Swap,
		d.Repository.ActiveSwap,
		d.Repository.Order,
		d.Repository.Account,
		d.Repository.Deposit,
		d.Repository.Transfer,
		d.Repository.Withdraw,
		d.Repository.Coin,
		d.Repository.OrderFee,
		d.Repository.Symbol,
		d.Repository.Orderbook,
		d.Repository.User,
		d.Service.SwapStatusUpdater,
		decimal.NewFromFloat(cfg.OurFee),
		decimal.NewFromFloat(cfg.MaxLimit),
		decimal.NewFromFloat(cfg.MinLimit),
		cfg.ByBit.MasterUid,
		ctx,
		d.Repository.Slack,
		cfg.DepositWaitingPeriod,
		swapExecutorCh)

	return nil
}

func (d *Deps) makeController() {
	d.Controller.Admin = grpcController.NewAdmin(d.Service.Admin, d.Service.Swap)
	d.Controller.Coin = grpcController.NewCoin(d.Service.Coin)
	d.Controller.Swap = grpcController.NewSwap(d.Service.Swap, d.Service.SwapStatusSubscriber)
	d.Controller.Withdraw = grpcController.NewWithdraw(d.Service.Withdraw)
}

func (d *Deps) makeWorkers(ctx context.Context, cfg config.Config, collector *metrics.AppMetrics, exchangeExecutorCh chan *model.Swap,
	updateCh chan model.OrderBookUpdateMessage, brokenOrderBookWsCh chan []*model.Symbol) (*worker.SwapExecutor, error) {
	newSymbolsCh := make(chan []*model.Symbol)

	swapExecutor := worker.NewSwapExecutor(
		d.Repository.ActiveSwap,
		d.Repository.ErrorCounter,
		d.Service.Swap,
		d.Service.SwapStatusUpdater,
		d.Service.Order,
		d.Service.Withdraw,
		d.Service.Transfer,
		exchangeExecutorCh,
		d.Repository.Slack,
		cfg.DepositWaitingPeriod,
		cfg.SwapExecutorWorkerGroup,
		int32(cfg.BusyWorkersThresholdForAlert),
	)
	d.Workers = append(d.Workers,
		swapExecutor,
		worker.NewOrderBookManager(d.Service.OrderBook, updateCh, newSymbolsCh),
		worker.NewSymbolUpdater(d.Service.Symbol, newSymbolsCh),
		worker.NewOrderFeeUpdater(d.Service.OrderFee),
		worker.NewCoinsUpdater(d.Service.Coin),
		worker.NewWsRecover(d.ExternalClient.Subscriber, brokenOrderBookWsCh),
		worker.NewMetrics(d.Repository.Swap, d.ExternalClient.ExchangeAccount, collector.SubAccountsGauge, collector.SwapStatusesGauge),
	)

	return swapExecutor, nil
}

func startSwaps(ctx context.Context, deps *Deps, swapExecutor *worker.SwapExecutor) error {
	// Если ждём депозит, то подписываемся на wallet
	// Если будем размещать ордер, подписываемся на order (подписка сервисом на предыдущем шаге)
	completed := []model.Status{model.Unknown, model.Completed, model.ManualCompleted, model.Cancel, model.DepositError, model.Error}
	activeSwaps, err := deps.Repository.Swap.Find(ctx, &model.SwapFilter{
		NotEqStatus: completed,
	})
	if err != nil {
		return fmt.Errorf("getAllActive: %w", err)
	}
	for i := range activeSwaps {
		var account *model.Account
		if activeSwaps[i].Status == model.WaitDeposit || activeSwaps[i].Status == model.PlaceOrder {
			account, err = deps.Repository.Account.FindOne(ctx, &model.AccountFilter{
				ID: &activeSwaps[i].AccountFrom,
			})
			if err != nil {
				return fmt.Errorf("getSubAccount: %w", err)
			}
		}

		// Актуально только когда ожидаем получения депозитов по websocket, не по REST
		//if activeSwaps[i].Status == model.WaitDeposit {
		//	err = deps.ExternalClient.Subscriber.SubscribeOnWallet(ctx, activeSwaps[i].ID, activeSwaps[i].StartTime, account, activeSwaps[i].CoinFrom, activeSwaps[i].AmountFrom)
		//	if err != nil {
		//		return fmt.Errorf("subscribeOnWallet: %w", err)
		//	}
		//
		//}

		if activeSwaps[i].Status == model.PlaceOrder {
			orders, err := deps.Repository.Order.Find(ctx, &model.OrderFilter{
				AccountID: &activeSwaps[i].AccountFrom,
			})
			if err != nil {
				return fmt.Errorf("getByAccountID: %w", err)
			}

			err = deps.ExternalClient.Subscriber.SubscribeOnOrders(ctx, account, orders, false)
			if err != nil {
				return fmt.Errorf("subscribeOnOrders: %w", err)
			}
		}
	}

	// Стартуем выполнение всех не завершённых swaps
	err = swapExecutor.UploadSwaps(ctx)
	if err != nil {
		return fmt.Errorf("uploadSwaps: %w", err)
	}

	return nil
}

// subscribeOnOrderbooks подписывается на ордербуки по всем символом
// для быстрого старта приложения, каждый конекшн обрабатывается в отдельной горутине
func subscribeOnOrderbooks(ctx context.Context, deps *Deps) {
	allSymbols, err := deps.Repository.Symbol.GetAll(ctx)
	if err != nil {
		log.Fatal(ctx, "repository.Symbol.GetAll: %s", err.Error())
		return
	}
	var (
		symbols        []*model.Symbol
		dividedSymbols [][]*model.Symbol
	)
	for i := range allSymbols {
		if (i+1)%service.SymbolsPerConnection == 0 {
			dividedSymbols = append(dividedSymbols, symbols)
			symbols = []*model.Symbol{}
		}
		symbols = append(symbols, allSymbols[i])
	}
	if len(symbols) > 0 {
		dividedSymbols = append(dividedSymbols, symbols)
	}
	for _, syms := range dividedSymbols {
		go func() {
			err = deps.Service.OrderBook.Subscribe(ctx, syms)
			if err != nil {
				log.Fatal(ctx, "subscribe orderbook: %s", err.Error())
				return
			}
		}()
	}
}
