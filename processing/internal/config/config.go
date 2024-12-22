package config

import (
	"time"

	"code.emcdtech.com/emcd/sdk/app"
	"github.com/google/uuid"
)

type InvoiceServiceConfig struct {
	InvoiceTTL time.Duration `env:"INVOICE_TTL"`
}

type CoinConfig struct {
	FetchFrequency time.Duration `env:"COIN_FETCHER_FREQUENCY"`
	AvailableCoins []string      `env:"COIN_AVAILABLE_COINS"`
}

type GRPCClients struct {
	AddressServiceAddr     string `env:"ADDRESS_GRPC_ADDR"`
	ProfileServiceAddr     string `env:"PROFILE_GRPC_ADDR"`
	UserAccountServiceAddr string `env:"USER_ACCOUNT_GRPC_ADDR"`
	CoinServiceAddr        string `env:"COIN_GRPC_ADDR"`
	AccountingServiceAddr  string `env:"ACCOUNTING_GRPC_ADDR"`
}

type APIConfig struct {
	app.DepsConfig
	InvoiceService InvoiceServiceConfig
	Coin           CoinConfig
	GRPCClients    GRPCClients
}

type RabbitMQConfig struct {
	URL string `env:"RABBITMQ_URL"`
}

type CoinwatchClientConfig struct {
	app.DepsConfig
	RabbitMQ                    RabbitMQConfig
	GRPCClients                 GRPCClients
	CoinwatchProcessingExchange string    `env:"COINWATCH_PROCESSING_EXCHANGE"`
	FeeCollectorUserID          uuid.UUID `env:"FEE_COLLECTOR_USER_ID"`
}
