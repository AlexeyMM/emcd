# Config

Пакет предназначен для стандартизации переменных окружения используемых для настройки стандартных инфраструктурных 
зависимостей (например: amp tracing, pg, grpc server и т.д.)

В файле [.env.example](./.env.example) по разделам моэно увидеть перечень всех доступных env для каждой из зависимостей.


Пример использования:
```go
package main

import (
	...
	"code.emcdtech.com/emcd/sdk/config"
	...
)

type Config struct {
	Environment config.Environment
	Log         config.Log
	GRPC        config.GRPCServer
	Tracing     config.APMTracing
	PGXPool     config.PGXPool
	HTTP        config.HTTPServer
}

func New() (*Config, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Caller().Stack().Logger()
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	logLevel, err := zerolog.ParseLevel(cfg.Environment.Name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	zerolog.SetGlobalLevel(logLevel)
	return &cfg, nil
}


func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	pool, err := cfg.PGXPool.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	....
}
```
