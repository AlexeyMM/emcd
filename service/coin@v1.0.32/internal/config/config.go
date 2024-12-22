package config

import (
	_ "github.com/joho/godotenv/autoload"

	"code.emcdtech.com/emcd/sdk/config"
)

type Config struct {
	Environment        config.Environment
	Log                config.Log
	GRPC               config.GRPCServer
	Tracing            config.APMTracing
	PGXPool            config.PGXPool
	HTTP               config.HTTPServer
	NodeAddress        string  `env:"NODE_ADDRESS"`
	NodeAuthKey        string  `env:"NODE_AUTH_KEY"`
	WhiteLabelAddress  string  `env:"WHITELABEL_ADDRESS"`
	RateAddress        string  `env:"RATE_ADDRESS"`
	FeeMultiplierErc20 float64 `env:"FEE_MULTIPLIER_ERC20"`
}
