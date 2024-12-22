package config

import (
	_ "github.com/joho/godotenv/autoload"

	"code.emcdtech.com/emcd/sdk/config"
)

type Config struct {
	Environment       config.Environment
	Log               config.Log
	GRPC              config.GRPCServer
	Tracing           config.APMTracing
	PGXPool           config.PGXPool
	HTTP              config.HTTPServer
	NodeAddress       string `env:"NODE_ADDRESS"`
	WhiteLabelAddress string `env:"WHITELABEL_ADDRESS"`
}
