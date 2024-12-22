package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
	"code.emcdtech.com/b2b/processing/protocol/coinpb"
	"code.emcdtech.com/b2b/swap/protocol/swap"
	"code.emcdtech.com/b2b/swap/protocol/swapCoin"
	"code.emcdtech.com/b2b/swap/protocol/swapWithdraw"
	"code.emcdtech.com/b2b/swap/repository"
	"code.emcdtech.com/emcd/sdk/app"
	cfgSdk "code.emcdtech.com/emcd/sdk/config"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/email/protocol/email"
	"github.com/go-playground/validator/v10"
	goGrpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "code.emcdtech.com/b2b/endpoint/docs"
	"code.emcdtech.com/b2b/endpoint/internal/config"
	"code.emcdtech.com/b2b/endpoint/internal/controller"
	myGrpc "code.emcdtech.com/b2b/endpoint/internal/controller/grpc"
	"code.emcdtech.com/b2b/endpoint/internal/encryptor"
	internalMiddleware "code.emcdtech.com/b2b/endpoint/internal/middleware"
	internalRepository "code.emcdtech.com/b2b/endpoint/internal/repository"
	"code.emcdtech.com/b2b/endpoint/internal/repository/pg"
	"code.emcdtech.com/b2b/endpoint/internal/service"
	"code.emcdtech.com/b2b/endpoint/protocol/b2bEndpointAdmin"
)

const shutdownTimeout = 2 * time.Second

// serviceName and serviceVersion is set during the build in CI/CD pipeline using ldflags
// (eg.: go build -ldflags="-X 'main.serviceVersion=<release id or tag>'").
var (
	serviceName    string = "b2b-endpoint"
	serviceVersion string = "local-build"
)

// @title b2b-endpoint
// @version 1.0

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

	deps, err := NewDeps(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "new deps: %s", err.Error())
		return
	}

	var grpcSrv *grpc.Server

	err = app.New().
		WithTracing(deps.Infrastructure.Tracer).
		WithPprof(deps.Infrastructure.HTTPSrv).
		WithGRPC(
			func(opts ...grpc.ServerOption) *grpc.Server {
				grpcSrv = deps.Infrastructure.GRPCSrv(opts...)
				return grpcSrv
			},
			//healthchecker.NewCommon(deps.Infrastructure.DB),
		).
		WithMetrics(deps.Infrastructure.HTTPSrv).
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
		go func() {
			log.Panic(ctx, err.Error())
		}()
	}
	log.Info(ctx, "main tail is start gracefully shutdown...")

}

// Deps dependencies
type Deps struct {
	Infrastructure struct {
		Tracer  trace.TracerProvider
		HTTPSrv *echo.Echo
		GRPCSrv app.GRPCServerFactory
		DB      *pgxpool.Pool
	}
	Tools struct {
		Encryptor encryptor.Encryptor
	}
	Repository struct {
		Swap         repository.Swap
		SwapCoin     repository.Coin
		SwapWithdraw repository.Withdraw
		Client       internalRepository.Client
		Email        internalRepository.Email
		IP           internalRepository.IP
		RequestLog   internalRepository.RequestLog
		Secret       internalRepository.Secret
	}
	Middleware struct {
		Auth internalMiddleware.Auth
	}
	Service struct {
		Client       service.Client
		IP           service.IP
		Secret       service.Secret
		Swap         service.Swap
		SwapCoin     service.SwapCoin
		SwapWithdraw service.SwapWithdraw
	}
	GRPCClients struct {
		Swap            swap.SwapServiceClient
		SwapCoin        swapCoin.SwapCoinServiceClient
		SwapWithdraw    swapWithdraw.SwapWithdrawServiceClient
		Email           email.EmailServiceClient
		Processing      buyerpb.InvoiceBuyerServiceClient
		ProcessingCoins coinpb.CoinsServiceClient
	}
	Controller struct {
		Auth            *controller.Auth
		Swap            *controller.Swap
		SwapCoin        *controller.SwapCoin
		SwapUser        *controller.SwapUser
		SwapWithdraw    *controller.SwapWithdraw
		Admin           *myGrpc.Admin
		ProcessingBuyer *controller.ProcessingBuyerController
		ProcessingCoins *controller.ProcessingCoinsController
	}
	//Collectors    []prometheus.Collector
	DeferHandlers []func(ctx context.Context)
}

func NewDeps(ctx context.Context, cfg config.Config) (Deps, error) {
	var deps Deps
	if err := deps.makeInfrastructure(ctx, cfg); err != nil {
		return deps, fmt.Errorf("configurer infrastructure deps: %w", err)
	}

	if err := deps.makeClients(cfg); err != nil {
		return deps, fmt.Errorf("configurer clients deps: %w", err)
	}

	deps.makeTools(&cfg)
	deps.makeRepository()
	deps.MakeMiddleware()

	deps.makeService()
	deps.makeController()

	deps.setAuthRoutes()
	deps.setSwapRoutes()
	deps.setTestRoutes()
	deps.setProcessingRoutes()
	deps.setTechRoutes(cfg.Environment.IsProduction())

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
	otel.SetTracerProvider(d.Infrastructure.Tracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	d.Infrastructure.DB, err = cfg.PGXPool.New(ctx, d.Infrastructure.Tracer)
	if err != nil {
		return fmt.Errorf("pool to referral statistic: %w", err)
	}

	d.Infrastructure.HTTPSrv = echo.New()
	d.Infrastructure.HTTPSrv.HideBanner = false
	d.Infrastructure.HTTPSrv.HidePort = false
	d.Infrastructure.HTTPSrv.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{
				"https://emcd.io",
				"https://stage.mytstnv.site",
				"http://localhost:8080",
				"https://b2b-pay-form-web.index.dev.emcd.io",
				"https://b2b-pay-form-web.stage.mytstnv.site",
			},
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
	d.Infrastructure.HTTPSrv.Validator = controller.NewCustomValidator(validator.New())

	d.Infrastructure.GRPCSrv = func(opts ...grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpc.UnaryInterceptor(sdkError.ServerUnaryInterceptor))
		grpcSrv := grpc.NewServer(opts...)
		b2bEndpointAdmin.RegisterEndpointAdminServiceServer(grpcSrv, d.Controller.Admin)
		reflection.Register(grpcSrv)
		return grpcSrv
	}
	return nil
}

func (d *Deps) makeTools(cfg *config.Config) {
	d.Tools.Encryptor = encryptor.NewEncryptor(cfg.EncryptorKey)
}

func (d *Deps) makeRepository() {
	d.Repository.Swap = repository.NewSwap(d.GRPCClients.Swap)
	d.Repository.SwapCoin = repository.NewCoin(d.GRPCClients.SwapCoin)
	d.Repository.SwapWithdraw = repository.NewWithdraw(d.GRPCClients.SwapWithdraw)

	d.Repository.Client = pg.NewClient(d.Infrastructure.DB)
	d.Repository.Email = internalRepository.NewEmailImp(d.GRPCClients.Email)
	d.Repository.IP = pg.NewIP(d.Infrastructure.DB)
	d.Repository.RequestLog = pg.NewRequestLog(d.Infrastructure.DB)
	d.Repository.Secret = pg.NewSecret(d.Infrastructure.DB, d.Tools.Encryptor)
}

func (d *Deps) MakeMiddleware() {
	d.Middleware.Auth = internalMiddleware.NewAuth(d.Repository.Secret, d.Repository.IP, d.Repository.RequestLog)
}

func (d *Deps) makeService() {
	d.Service.Client = service.NewClient(d.Repository.Client)
	d.Service.IP = service.NewIP(d.Repository.IP)
	d.Service.Secret = service.NewSecret(d.Repository.Secret, d.Repository.IP)
	d.Service.Swap = service.NewSwap(d.Repository.Swap, d.Repository.Email)
	d.Service.SwapCoin = service.NewSwapCoin(d.Repository.SwapCoin)
	d.Service.SwapWithdraw = service.NewSwapWithdraw(d.Repository.SwapWithdraw)
}

func (d *Deps) makeClients(cfg config.Config) error {
	grpcOpts := []grpc.DialOption{
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

	conn, err := grpc.NewClient(cfg.SwapAddress, grpcOpts...)
	if err != nil {
		return fmt.Errorf("new grpc client skeleton: %w", err)
	}

	connEmail, err := grpc.NewClient(cfg.EmailAddress, grpcOpts...)
	if err != nil {
		return fmt.Errorf("new email grpc client: %w", err)
	}

	connProcessing, err := grpc.NewClient(cfg.ProcessingAddress, grpcOpts...)
	if err != nil {
		return fmt.Errorf("new processing grpc client: %w", err)
	}

	d.DeferHandlers = append(d.DeferHandlers,
		func(ctx context.Context) {
			_ = conn.Close()
			_ = connEmail.Close()
			_ = connProcessing.Close()
		},
	)

	d.GRPCClients.Swap = swap.NewSwapServiceClient(conn)
	d.GRPCClients.SwapCoin = swapCoin.NewSwapCoinServiceClient(conn)
	d.GRPCClients.SwapWithdraw = swapWithdraw.NewSwapWithdrawServiceClient(conn)
	d.GRPCClients.Email = email.NewEmailServiceClient(connEmail)
	d.GRPCClients.Processing = buyerpb.NewInvoiceBuyerServiceClient(connProcessing)
	d.GRPCClients.ProcessingCoins = coinpb.NewCoinsServiceClient(connProcessing)
	return nil
}

func (d *Deps) makeController() {
	d.Controller.Auth = controller.NewAuth(d.Service.Secret)
	d.Controller.Swap = controller.NewSwap(d.Service.Swap)
	d.Controller.SwapCoin = controller.NewSwapCoin(d.Service.SwapCoin)
	d.Controller.SwapUser = controller.NewSwapUser(d.Service.Swap)
	d.Controller.SwapWithdraw = controller.NewSwapWithdraw(d.Service.SwapWithdraw)
	d.Controller.Admin = myGrpc.NewAdmin(d.Service.Client, d.Service.Secret, d.Service.IP)
	d.Controller.ProcessingBuyer = controller.NewProcessingBuyerController(d.GRPCClients.Processing)
	d.Controller.ProcessingCoins = controller.NewProcessingCoinsController(d.GRPCClients.ProcessingCoins)
}

func (d *Deps) setAuthRoutes() {
	d.Infrastructure.HTTPSrv.POST("/auth/rotate-keys", d.Controller.Auth.RotateKeys, d.Middleware.Auth.ValidateRequest)
}

func (d *Deps) setSwapRoutes() {
	d.Infrastructure.HTTPSrv.GET("/swap/estimate", d.Controller.Swap.Estimate)
	d.Infrastructure.HTTPSrv.POST("/swap", d.Controller.Swap.Swap)
	d.Infrastructure.HTTPSrv.GET("/swap/status/:swapID", d.Controller.Swap.Status)
	d.Infrastructure.HTTPSrv.GET("/swap/:swapID", d.Controller.Swap.GetSwapByID)
	d.Infrastructure.HTTPSrv.POST("/swap/support/message", d.Controller.Swap.SendSupportMessage)
	d.Infrastructure.HTTPSrv.POST("/swap/email", d.Controller.Swap.SendSwapInfoEmail)

	d.Infrastructure.HTTPSrv.GET("/swap/coins", d.Controller.SwapCoin.GetAllV1)
	d.Infrastructure.HTTPSrv.GET("/v2/swap/coins", d.Controller.SwapCoin.GetAllV2)

	d.Infrastructure.HTTPSrv.POST("/swap/user", d.Controller.SwapUser.CreateUserBySwapID)

	d.Infrastructure.HTTPSrv.GET("/swap/withdraw/tx_link/:swap_id", d.Controller.SwapWithdraw.GetTransactionLink)
}

func (d *Deps) setTestRoutes() {
	d.Infrastructure.HTTPSrv.GET("/test/swap/estimate", d.Controller.Swap.Estimate, d.Middleware.Auth.ValidateRequest)
	d.Infrastructure.HTTPSrv.POST("/test/swap", d.Controller.Swap.Swap, d.Middleware.Auth.ValidateRequest)
}

func (d *Deps) setProcessingRoutes() {
	d.Infrastructure.HTTPSrv.GET("/processing/coins", d.Controller.ProcessingCoins.GetCoins)
	d.Infrastructure.HTTPSrv.GET("/processing/invoice/:invoice_id", d.Controller.ProcessingBuyer.GetInvoice)
	d.Infrastructure.HTTPSrv.GET("/processing/invoice_form/:form_id", d.Controller.ProcessingBuyer.GetInvoiceForm)
	d.Infrastructure.HTTPSrv.PUT("/processing/invoice_form/:form_id", d.Controller.ProcessingBuyer.SubmitInvoiceForm)
	d.Infrastructure.HTTPSrv.GET("/processing/merchant/:id/calculate_invoice_payment", d.Controller.ProcessingBuyer.CalculateInvoicePayment)
}

func (d *Deps) setTechRoutes(isProduction bool) {
	if !isProduction {
		d.Infrastructure.HTTPSrv.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
