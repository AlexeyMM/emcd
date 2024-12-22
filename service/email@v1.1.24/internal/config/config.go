package config

import (
	"code.emcdtech.com/emcd/sdk/config"
)

type Config struct {
	Environment             config.Environment
	Log                     config.Log
	GRPC                    config.GRPCServer
	Tracing                 config.APMTracing
	PGXPool                 config.PGXPool
	HTTP                    config.HTTPServer
	WhiteLabelAddress       string            `env:"WHITELABEL_ADDRESS,required"`
	ProfileAddress          string            `env:"PROFILE_ADDRESS,required"`
	ChangeWalletAddressLink string            `env:"CHANGE_WALLET_ADDRESS_LINK,required"`
	Domains                 map[string]string `env:"DOMAINS,required"                    envSeparator:","`
}
