package main

import (
	"context"
	"errors"
	"math"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	sdkErr "code.emcdtech.com/emcd/sdk/error"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"code.emcdtech.com/emcd/service/referral/protocol/reward"
	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"code.emcdtech.com/emcd/service/accounting/internal/config"
	"code.emcdtech.com/emcd/service/accounting/internal/handler"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/internal/worker"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/pkg/kafka"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	referralPb "code.emcdtech.com/emcd/service/accounting/protocol/referral"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
)

const redisTimeout = 15 * time.Second
const outboxBatchSize = 1000
const serviceName = "accounting"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if errors.Is(ctx.Err(), context.Canceled) {
			sdkLog.Info(ctx, "service main is gracefully shutdown.")

		} else {
			cancel()

		}
	}()

	ctx = sdkLog.CreateClientNamedContext(ctx, serviceName)

	cfg, err := config.NewConfig()
	if err != nil {

		sdkLog.Panic(ctx, "can't read config: %v", err)

	}

	pool, err := pgxpool.New(ctx, cfg.PostgresConnectionString)
	if err != nil {

		sdkLog.Panic(ctx, "can't connect to postgres: %v", err)

	}
	if err = pool.Ping(ctx); err != nil {

		sdkLog.Panic(ctx, "can't ping to postgres: %v", err)

	}

	slaveDBPool, err := pgxpool.New(ctx, cfg.SlaveDBConnectionString)
	if err != nil {

		sdkLog.Panic(ctx, "can't connect to slave DB: %v", err)

	}

	if err = pool.Ping(ctx); err != nil {

		sdkLog.Panic(ctx, "can't ping to slave DB: %v", err)

	}

	rds := newRedisClient(ctx, cfg)

	grpcOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			sdkLog.ClientUnaryInterceptor,
			sdkErr.ClientUnaryInterceptor,
		),
	}

	rewardConn, err := grpc.NewClient(cfg.RewardAddress, grpcOptions...)
	if err != nil {

		sdkLog.Panic(ctx, "can't connect to reward: %v", err)

	}
	defer func(rewardConn *grpc.ClientConn) {
		if err := rewardConn.Close(); err != nil {
			sdkLog.Error(ctx, "can't close reward connection: %v", err)
		}
	}(rewardConn)

	rewardCli := reward.NewRewardServiceClient(rewardConn)
	rewardRepo := repository.NewReward(rewardCli)

	whiteLabelConn, err := grpc.NewClient(cfg.WhiteLabelAddress, grpcOptions...)
	if err != nil {

		sdkLog.Panic(ctx, "can't connect to whitelabel: %v", err)

	}
	defer func(whiteLabelConn *grpc.ClientConn) {
		if err := whiteLabelConn.Close(); err != nil {
			sdkLog.Error(ctx, "can't close whitelabel connection: %v", err)

		}
	}(whiteLabelConn)
	whiteLabelCli := whitelabel.NewWhitelabelServiceClient(whiteLabelConn)
	coinStrIDs := make(map[int]string)
	for i, coinIntID := range cfg.WalletCoinsIntIDs {
		coinStrIDs[coinIntID] = cfg.WalletCoinsStrIDs[i]
	}

	coinConn, err := grpc.NewClient(cfg.CoinAddress, grpcOptions...)
	if err != nil {

		sdkLog.Panic(ctx, "can't connect to coin: %v", err)

	}
	defer func(coinConn *grpc.ClientConn) {
		if err := coinConn.Close(); err != nil {
			sdkLog.Error(ctx, "can't close coin: %v", err)

		}
	}(coinConn)

	coinClient := coin.NewCoinServiceClient(coinConn)
	coinValidator := coinValidatorRepo.NewCoinValidatorRepository(coinClient)

	// coinRepository := repository.NewCoin(coinClient)

	whiteLabelRepo := repository.NewWhiteLabel(whiteLabelCli)
	balanceRepository := repository.NewBalance(pool, cfg.WhiteListBalanceUsers)
	referralsStatisticRepository := repository.NewReferralStatistic(pool)

	userAccountRepository := repository.NewUserAccountRepo(pool)
	historyWalletsRepository := repository.NewWalletsHistory(pool, coinStrIDs)
	historyMiningIncomesRepository := repository.NewIncomesHistory(pool, slaveDBPool, coinValidator)
	historyMiningPayoutsRepository := repository.NewPayoutsHistory(pool, slaveDBPool, coinValidator)
	limitPayoutRepository := repository.NewLimitPayouts(pool, rds, cfg.CoinsLimits)
	transactionRepository := repository.NewTransactionStore(pool)
	payoutRepository := repository.NewPayouts(pool)
	checkerRepository := repository.NewChecker(pool)

	balanceService, err := service.NewBalance(cfg.WalletCoinsIntIDs, coinStrIDs, balanceRepository, pool, cfg.ServiceData, rewardRepo, whiteLabelRepo, toMap(cfg.IgnoreReferralPaymentByUserID))
	if err != nil {
		sdkLog.Panic(ctx, "service.NewBalance: %v", err)
	}
	historyWalletsService := service.NewWalletsHistory(historyWalletsRepository)
	historyMiningIncomesService := service.NewIncomesHistory(historyMiningIncomesRepository)
	historyMiningPayoutsService := service.NewPayoutsHistory(historyMiningPayoutsRepository)
	limitPayoutService := service.NewLimitPayouts(limitPayoutRepository, cfg.SlackUrl)
	referralService := service.NewReferral(referralsStatisticRepository, userAccountRepository)
	payoutService := service.NewPayout(payoutRepository)
	checkerService := service.NewChecker(checkerRepository)
	transactionService := service.NewTransaction(transactionRepository)

	accountingUserAccountService := service.NewAccountingUserAccountService(userAccountRepository)
	accountingService := handler.NewAccountingHandler(
		balanceService, historyWalletsService, historyMiningIncomesService, historyMiningPayoutsService, limitPayoutService,
		accountingUserAccountService, transactionService, payoutService, checkerService, coinValidator,
	)
	referralStatisticService := handler.NewReferralHandler(referralService)
	userAccountService := service.NewUserAccountService(userAccountRepository)
	userAccountHandler := handler.NewUserAccountHandler(userAccountService, coinValidator)

	writerKafka := kafka.NewWriter(
		ctx,
		cfg.KafkaBrokers,
		cfg.KafkaSaslUser,
		cfg.KafkaSaslPassword,
		cfg.KafkaSaslEnable,
		outboxBatchSize)
	outboxTransactionsRepository := repository.NewOutboxTransactions()
	outboxTransactionService := service.NewOutboxTransaction(
		pool,
		balanceRepository,
		writerKafka,
		outboxTransactionsRepository,
		outboxBatchSize,
	)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	cronWorker := worker.NewCronWorker(time.Second, outboxTransactionService.ExportToKafka)
	go func() {
		defer wg.Done()
		if err := cronWorker.Run(ctx); err != nil {
			sdkLog.Error(ctx, err.Error())

			return
		}
	}()

	listen, err := net.Listen("tcp", cfg.Port)
	if err != nil {

		sdkLog.Panic(ctx, "can't listen: %v", err)

	}

	server := grpc.NewServer(grpc.MaxSendMsgSize(math.MaxInt64),
		grpc.ChainUnaryInterceptor(
			sdkLog.ServerUnaryNamedInterceptor(serviceName),
			sdkErr.ServerUnaryInterceptor,
		),
	)

	accountingPb.RegisterAccountingServiceServer(server, accountingService)
	referralPb.RegisterAccountingReferralServiceServer(server, referralStatisticService)
	userAccountPb.RegisterUserAccountServiceServer(server, userAccountHandler)
	reflection.Register(server)

	sdkLog.Info(ctx, "accounting service is running")
	go func() {
		if err := server.Serve(listen); err != nil {
			sdkLog.Error(ctx, "can't start grpc server: %v", err)
		}
	}()

	wg.Add(1)
	go coinValidator.Serve(ctx, wg, cfg.CoinCacheUpdateInterval)

	if cfg.PostgresMigrate {
		wg.Add(1)
		go migrate(ctx, wg, userAccountRepository, cfg.PostgresMigrateDelay)

	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sdkLog.Info(ctx, "main service is started.")
	<-quit
	sdkLog.Info(ctx, "signal received and start main gracefully shutdown...")

	server.GracefulStop()
	cancel()
	wg.Wait()
	sdkLog.Info(ctx, "wait group released")

}

func newRedisClient(ctx context.Context, cfg *config.Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Network:            "",
		Addr:               cfg.RedisHost + ":" + cfg.RedisPort,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           cfg.RedisUsername,
		Password:           cfg.RedisPassword,
		DB:                 cfg.RedisBase,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        redisTimeout,
		ReadTimeout:        redisTimeout,
		WriteTimeout:       redisTimeout,
		PoolFIFO:           false,
		PoolSize:           cfg.RedisPool * runtime.GOMAXPROCS(0),
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {

		sdkLog.Panic(ctx, "redis is working on address %s", cfg.RedisHost)

	}

	return client
}

func toMap(arr []uuid.UUID) map[uuid.UUID]bool {
	res := make(map[uuid.UUID]bool, len(arr))
	for i := range arr {
		res[arr[i]] = true
	}
	return res
}

func migrate(ctx context.Context, wg *sync.WaitGroup, userAccountRepo repository.UserAccountRepo, delay time.Duration) {
	const batchSize = 250
	const delayMilliParts = 100

	defer wg.Done()
	sdkLog.Info(ctx, "migrator has been started with delay %v", delay)

	for i := 0; ; i++ {
		for j := 0; j < delayMilliParts; j++ {
			time.Sleep(delay / delayMilliParts)

			if ctx.Err() != nil {
				break
			}
		}

		if ctx.Err() != nil {
			sdkLog.Info(ctx, "migrator has gracefully shutdown.")

			break
		}

		emptyStruct := struct{}{}
		filterForUserIdNew := &model.UserAccountFilter{
			ID:            nil,
			UserID:        nil,
			AccountTypeID: nil,
			UserIDNew:     nil,
			CoinNew:       nil,
			IsActive:      nil,
			Pagination: &model.Pagination{
				Limit:  batchSize,
				Offset: 0,
			},
			UserIDNewIsNull: &emptyStruct,
			CoinNewIsNull:   nil,
		}

		if totalReturn, userAccounts, err := userAccountRepo.FindUserAccountByFilterMigrateOnly(ctx, filterForUserIdNew); err != nil {
			sdkLog.Error(ctx, "migrator failed find user account by filter: %s", err)

		} else if len(userAccounts) > 0 {
			if err := userAccountRepo.UpdateUserAccountForMigrateUserIdNew(ctx, userAccounts); err != nil {
				sdkLog.Error(ctx, "migrator failed update user account for user_id_new: %s", err)

			} else {
				sdkLog.Info(ctx, "migrator success update %d/%d user account for user_id_new", len(userAccounts), *totalReturn)

			}
		} else if i%int(time.Hour/delay) == 0 { // every 1 hour if nothing migrate
			sdkLog.Error(ctx, "migrator not found any migration for user_id_new")

		}

		filterForCoinNew := &model.UserAccountFilter{
			ID:            nil,
			UserID:        nil,
			AccountTypeID: nil,
			UserIDNew:     nil,
			CoinNew:       nil,
			IsActive:      nil,
			Pagination: &model.Pagination{
				Limit:  batchSize,
				Offset: 0,
			},
			UserIDNewIsNull: nil,
			CoinNewIsNull:   &emptyStruct,
		}

		if totalReturn, userAccounts, err := userAccountRepo.FindUserAccountByFilterMigrateOnly(ctx, filterForCoinNew); err != nil {
			sdkLog.Error(ctx, "migrator failed find user account by filter: %s", err)

		} else if len(userAccounts) > 0 {
			if err := userAccountRepo.UpdateUserAccountForMigrateCoinNew(ctx, userAccounts); err != nil {
				sdkLog.Error(ctx, "migrator failed update user account for coin_new: %s", err)

			} else {
				sdkLog.Info(ctx, "migrator success update %d/%d user account for coin_new", len(userAccounts), *totalReturn)

			}
		} else if i%int(time.Hour/delay) == 0 { // every 1 hour if nothing migrate
			sdkLog.Info(ctx, "migrator not found any migration for coin_new")

		}
	}
}
