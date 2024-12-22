package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var _HealthCheckServingStatusToHealthCheckResponse = map[HealthCheckServingStatus]grpc_health_v1.HealthCheckResponse_ServingStatus{
	UnknownHealthCheckServingStatus:        grpc_health_v1.HealthCheckResponse_UNKNOWN,
	ServingHealthCheckServingStatus:        grpc_health_v1.HealthCheckResponse_SERVING,
	NotServingHealthCheckServingStatus:     grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	ServiceUnknownHealthCheckServingStatus: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
}

type GrpcHealth struct {
	grpc_health_v1.UnimplementedHealthServer
	healthCheckers map[string]HealthChecker
}

func NewGrpcHealth(healthCheckers []HealthChecker) *GrpcHealth {
	hc := make(map[string]HealthChecker, len(healthCheckers))
	for i := range healthCheckers {
		hc[healthCheckers[i].ServiceName()] = healthCheckers[i]
	}
	return &GrpcHealth{
		healthCheckers: hc,
	}
}

func (s *GrpcHealth) Check(
	ctx context.Context,
	request *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	healthChecker, ok := s.healthCheckers[request.Service]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("service `%s` not found", request.Service))
	}
	st := healthChecker.Check(ctx)
	return &grpc_health_v1.HealthCheckResponse{
		Status: _HealthCheckServingStatusToHealthCheckResponse[st],
	}, nil
}
