package config

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type RedisConfig struct {
	Host         string        `env:"REDIS_HOST"          envDefault:"localhost:6379"`
	Password     string        `env:"REDIS_PASSWORD"      envDefault:""`
	DB           int           `env:"REDIS_DB"            envDefault:"0"`
	PoolSize     int           `env:"REDIS_POOL_SIZE"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT"`
}

func (c RedisConfig) New(tp trace.TracerProvider) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         c.Host,
		Password:     c.Password,
		DB:           c.DB,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
	})
	// Enable tracing instrumentation.
	err := redisotel.InstrumentTracing(client,
		redisotel.WithDBStatement(true),
		redisotel.WithTracerProvider(tp),
	)
	if err != nil {
		return nil, fmt.Errorf("set instrument tracing: %w", err)
	}
	// when metrics are implemented
	// Enable metrics instrumentation.
	// err = redisotel.InstrumentMetrics(client, redisotel.WithMeterProvider(mp))
	// if err != nil {
	//	return nil, fmt.Errorf("set instrument metrics: %w", err)
	// }
	return client, nil
}
