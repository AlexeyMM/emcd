package app

import (
	"context"
)

type HealthCheckServingStatus int

const (
	UnknownHealthCheckServingStatus        HealthCheckServingStatus = 0
	ServingHealthCheckServingStatus        HealthCheckServingStatus = 1
	NotServingHealthCheckServingStatus     HealthCheckServingStatus = 2
	ServiceUnknownHealthCheckServingStatus HealthCheckServingStatus = 3
)

type HealthChecker interface {
	Check(ctx context.Context) HealthCheckServingStatus
	ServiceName() string
}
