package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/google/uuid"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	Port                     string `env:"PORT" envDefault:":8080"`
	PostgresConnectionString string `env:"POSTGRES_CONNECTION_STRING" envDefault:""`
	SlaveDBConnectionString  string `env:"SLAVE_DB_CONNECTION_STRING" envDefault:""`

	RedisHost     string `env:"REDIS_HOST" envDefault:""`
	RedisPort     string `env:"REDIS_PORT" envDefault:""`
	RedisUsername string `env:"REDIS_USERNAME" envDefault:""`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisBase     int    `env:"REDIS_BASE" envDefault:""`
	RedisPool     int    `env:"REDIS_POOL" envDefault:""`

	SlackUrl string `env:"SLACK_WEBHOOK_URL" envDefault:""`

	WalletCoinsIntIDs []int    `env:"WALLET_COINS_INT_IDS" envDefault:"1,2,4,5,6,7,8,10,11,13,14" envSeparator:","`
	WalletCoinsStrIDs []string `env:"WALLET_COINS_STR_IDS" envDefault:"btc,bch,ltc,dash,eth,etc,doge,usdt,usdc,ton,kas" envSeparator:","`

	CoinsLimits map[string]float64 `env:"COIN_LIMITS" envDefault:""`

	ServiceData ServiceData `env:"SERVICE_DATA" envDefault:""`

	WhiteListBalanceUsers []string `env:"WHITE_LIST_BALANCE_USERS" envDefault:"" envSeparator:""` // TODO: check it linter say set empty envSeparator

	RewardAddress                 string        `env:"REWARD_ADDRESS" envDefault:""`
	WhiteLabelAddress             string        `env:"WHITELABEL_ADDRESS"`
	CoinAddress                   string        `env:"COIN_ADDRESS"`
	ProfileAddress                string        `env:"PROFILE_ADDRESS"`
	IgnoreReferralPaymentByUserID []uuid.UUID   `env:"IGNORE_REFERRAL_PAYMENT_BY_USER_ID" envSeparator:","`
	CoinCacheUpdateInterval       time.Duration `env:"COIN_CACHE_UPDATE_INTERVAL" envDefault:"15m"`

	KafkaBrokers      []string `env:"KAFKA_BROKERS" envSeparator:","`
	KafkaSaslUser     string   `env:"KAFKA_SASL_USER"`
	KafkaSaslPassword string   `env:"KAFKA_SASL_PASSWORD" `
	KafkaSaslEnable   bool     `env:"KAFKA_SASL_ENABLE" envDefault:"false"`

	PostgresMigrate      bool          `env:"POSTGRES_MIGRATE" envDefault:"false"`
	PostgresMigrateDelay time.Duration `env:"POSTGRES_MIGRATE_DELAY" envDefault:"1s"`
}

type ServiceData struct {
	ExchangeUsername    string `env:"EXCHANGE_USERNAME" envDefault:""`
	PayoutsNodeUsername string `env:"PAYOUTS_NODE_USERNAME" envDefault:"payouts_node"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("can't parse env: %w", err)
	}
	return &cfg, nil
}
