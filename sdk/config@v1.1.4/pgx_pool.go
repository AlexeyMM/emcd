package config

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PGXPool struct {
	ConnectionString string `env:"POSTGRES_CONNECTION_STRING,required"`
}

func (cfg PGXPool) New(ctx context.Context, tracer trace.TracerProvider) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", errors.WithStack(err))
	}
	// подключение трассиров запросов
	pgxCfg.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithIncludeQueryParameters(),
		otelpgx.WithDisableQuerySpanNamePrefix(),
		otelpgx.WithTracerProvider(tracer),
	)
	pgxpool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", errors.WithStack(err))
	}

	if err = pgxpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("do ping: %w", errors.WithStack(err))
	}
	return pgxpool, nil
}

// NewPGXPoolFromDSN для случаев, где сервис имеет подключение к нескольким схемам и для каждой из
// них он будет переиспользовать эту функцию
func NewPGXPoolFromDSN(ctx context.Context, tracer trace.TracerProvider, dsn string) (*pgxpool.Pool, error) {
	pool := PGXPool{
		ConnectionString: dsn,
	}
	result, err := pool.New(ctx, tracer)
	if err != nil {
		return nil, fmt.Errorf("pool.New: %w", err)
	}

	return result, nil
}

// SqlxDBWrap создает sqlx.DB в виде обертки над pgxpool. Позволяет использовать эффективность управления пулом pgxpool,
// и возможностями sqlx.
func SqlxDBWrap(pool *pgxpool.Pool) *sqlx.DB {
	db := stdlib.OpenDBFromPool(pool)
	return sqlx.NewDb(db, "pgx")
}
