package healthchecker

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/sdk/app"
)

type Common struct {
	pool *pgxpool.Pool
}

func NewCommon(pool *pgxpool.Pool) *Common {
	return &Common{
		pool: pool,
	}
}

func (s *Common) Check(ctx context.Context) app.HealthCheckServingStatus {
	if err := s.pool.Ping(ctx); err != nil {
		return app.NotServingHealthCheckServingStatus
	}
	return app.ServingHealthCheckServingStatus
}

func (s *Common) ServiceName() string {
	return ""
}
