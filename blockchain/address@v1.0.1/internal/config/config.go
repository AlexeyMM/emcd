package config

import (
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"code.emcdtech.com/emcd/sdk/config"
	transactor "code.emcdtech.com/emcd/sdk/pg"
)

type Config struct {
	Environment config.Environment
	Log         config.Log
	GRPC        config.GRPCServer
	Tracing     config.APMTracing
	PGXPool     transactor.PGXPool
	HTTP        config.HTTPServer
	GrpcClient
	RabbitMQ
	CustomCfg
}

type GrpcClient struct {
	AccountingAddress string `env:"ACCOUNTING_ADDRESS,required"`
	ProfileAddress    string `env:"PROFILE_ADDRESS,required"`
	CoinAddress       string `env:"COIN_ADDRESS,required"`
	NodeAddress       string `env:"NODE_ADDRESS,required"`
}

type RabbitMQ struct {
	RabbitmqUrl          string `env:"RABBITMQ_URL,required"`
	RabbitmqExchangeName string `env:"RABBITMQ_EXCHANGE_NAME,required"`
}

type CustomCfg struct {
	EthMasterKeys                     []string        `env:"ETH_MASTER_KEYS,required"`
	AlphMasterKeys                    []string        `env:"ALPH_MASTER_KEYS,required"`
	IsNetworkOldWay                   map[string]bool `env:"IS_NETWORK_OLD_WAY,required"`
	MigrationPostgresConnectionString string          `env:"MIGRATION_POSTGRES_CONNECTION_STRING"`
	MigrationDatetimeGte              EnvTime         `env:"MIGRATION_DATETIME_GTE"`
	MigrationDatetimeLte              EnvTime         `env:"MIGRATION_DATETIME_LTE"`
}

type EnvTime time.Time

func (t *EnvTime) UnmarshalText(text []byte) error {
	tt, err := time.Parse(time.RFC3339Nano, string(text))
	*t = EnvTime(tt)

	return err
}

func (c *Config) Validate() error {
	if err := c.CustomCfg.Validate(); err != nil {

		return fmt.Errorf("invalid address config: %w", err)
	}

	return nil
}

func (c *CustomCfg) Validate() error {
	for network := range c.IsNetworkOldWay {
		if err := nodeCommon.NewNetworkEnum(network).Validate(); err != nil {

			return fmt.Errorf("invalid network: %v, %w", network, err)
		}
	}

	return nil
}

func (c *CustomCfg) GetIsNetworkOldWayMap() map[nodeCommon.NetworkEnum]bool {
	dump := make(map[nodeCommon.NetworkEnum]bool)

	for network, v := range c.IsNetworkOldWay {
		if err := nodeCommon.NewNetworkEnum(network).Validate(); err != nil {

			panic("must validate config before use")
		}

		dump[nodeCommon.NewNetworkEnum(network)] = v

	}

	return dump
}
