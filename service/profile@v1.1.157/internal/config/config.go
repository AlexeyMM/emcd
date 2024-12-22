package config

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/sdk/config"
)

type Config struct {
	Environment   config.Environment
	Log           config.Log
	PGXPool       config.PGXPool
	Tracing       config.APMTracing
	HTTP          config.HTTPServer
	GRPC          config.GRPCServer
	GRPCClient    GRPCClient
	APIKeyConfig  APIKeyConfig
	APIJobsConfig APIJobsConfig

	EmcdPostgresConnectionString string           `env:"EMCD_POSTGRES_CONNECTION_STRING,required"`
	SegmentKey                   string           `env:"SEGMENT_KEY,required"`
	AccessSecret                 string           `env:"ACCESS_SECRET,required"`
	IdenfyRetryDelayMinutes      time.Duration    `env:"IDENFY_RETRY_DELAY,required"`
	MinGetAllReferralsTake       int              `env:"MIN_GET_ALL_REFERRALS_TAKE"               envDefault:"20"`
	DefaultDonations             DefaultDonations `env:"DEFAULT_DONATIONS_JSON,required"`
}

type GRPCClient struct {
	WhitelabelAddress string `env:"WHITELABEL_ADDRESS,required"`
	ReferralAddress   string `env:"REFERRAL_ADDRESS,required"`
	EmailAddress      string `env:"EMAIL_ADDRESS,required"`
	CoinAddress       string `env:"COIN_ADDRESS,required"`
	AccountingAddress string `env:"ACCOUNTING_ADDRESS,required"`
	WalletAddress     string `env:"WALLET_ADDRESS,required"`
}

type APIKeyConfig struct {
	Salt   string `env:"API_KEY_SALT,required"`
	Secret string `env:"API_KEY_SECRET,required"`
}

type APIJobsConfig struct {
	Host  string `env:"API_JOBS_HOST,required"`
	Path  string `env:"API_JOBS_PATH,required"`
	Token string `env:"API_JOBS_TOKEN,required"`
}

type DefaultDonations map[uuid.UUID]decimal.Decimal

func (t *DefaultDonations) UnmarshalText(text []byte) error {
	res := make(map[uuid.UUID]decimal.Decimal)
	err := json.Unmarshal(text, &res)
	if err != nil {
		return err
	}
	*t = res
	return nil
}
