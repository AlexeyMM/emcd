package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"

	"code.emcdtech.com/emcd/sdk/config"
)

type Config struct {
	Environment config.Environment
	Log         config.Log
	GRPC        config.GRPCServer
	PGXPool     config.PGXPool
	HTTP        config.HTTPServer
	Languages   []string `env:"LANGUAGES" envSeparator:"," envDefault:"ru"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error with parsing env variables in config %w", err)
	}
	return &cfg, nil
}
