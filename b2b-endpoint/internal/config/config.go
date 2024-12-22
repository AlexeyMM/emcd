package config

import (
	"code.emcdtech.com/emcd/sdk/config"
	pg "code.emcdtech.com/emcd/sdk/pg"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Environment       config.Environment
	Log               config.Log
	Tracing           config.APMTracing
	HTTP              config.HTTPServer
	GRPC              config.GRPCServer
	PGXPool           pg.PGXPool
	SwapAddress       string `env:"SWAP_GRPC_ADDRESS"`
	EmailAddress      string `env:"EMAIL_GRPC_ADDRESS"`
	ProcessingAddress string `env:"PROCESSING_GRPC_ADDRESS"`
	EncryptorKey      string `env:"ENCRYPTOR_KEY"`
}
