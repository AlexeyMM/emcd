// Package config provides functions for handling configuration files.
package config

import (
	_ "github.com/joho/godotenv/autoload"

	"code.emcdtech.com/emcd/sdk/config"
)

type Config struct {
	Environment config.Environment
	Log         config.Log
	GRPC        config.GRPCServer
	Tracing     config.APMTracing
	PGXPool     config.PGXPool
	HTTP        config.HTTPServer
	GRPCClient  GRCPClient
}

type GRCPClient struct {
	Profile   string `env:"PROFILE_ADDR" required:"true"`
	PromoCode string `env:"PROMO_CODE_ADDR" required:"true"`
	Coin      string `env:"COIN_ADDR" required:"true"`
}
