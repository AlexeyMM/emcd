package config

import (
	"time"

	"code.emcdtech.com/emcd/sdk/config"

	pg "code.emcdtech.com/emcd/sdk/pg"
)

type Config struct {
	Environment                  config.Environment
	Log                          config.Log
	GRPC                         config.GRPCServer
	Tracing                      config.APMTracing
	PGXPool                      pg.PGXPool
	HTTP                         config.HTTPServer
	ByBit                        ByBit
	DepositWaitingPeriod         time.Duration `env:"DEPOSIT_WAITING_PERIOD"`
	OurFee                       float64       `env:"OUR_FEE" envDefault:"0.02"` // В десятичном выражении
	MaxLimit                     float64       `env:"MAX_LIMIT"`                 // Максимальный лимит swap, в USDT
	MinLimit                     float64       `env:"MIN_LIMIT"`                 // Минимальный лимит swap, в USDT
	SlackWebhookUrl              string        `env:"SLACK_WEBHOOK_URL"`
	EmailAddress                 string        `env:"EMAIL_ADDRESS"`
	CoinAddress                  string        `env:"COIN_ADDRESS"`
	SwapExecutorWorkerGroup      int           `env:"SWAP_EXECUTOR_WORKER_GROUP"`
	BusyWorkersThresholdForAlert int           `env:"BUSY_WORKER_THRESHOLD_FOR_ALERT"`
}

type ByBit struct {
	MasterUid     int           `env:"MASTER_UID" envDefault:"134640868"`
	ApiUrl        string        `env:"API_URL"`
	ApiKey        string        `env:"API_KEY"`
	ApiSecret     string        `env:"API_SECRET"`
	WsExpiredTime time.Duration `env:"WS_EXPIRED_TIME" envDefault:"10m"`
}
