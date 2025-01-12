package server

import (
	"context"

	encoderpc "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/client/encode"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthServer struct {
	healthpb.HealthServer
}

func (h *healthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	status := getServicesHealthStatus()

	logger.Info("Overall health status: %q", status)

	return &healthpb.HealthCheckResponse{
		Status: status,
	}, nil
}

func getServicesHealthStatus() healthpb.HealthCheckResponse_ServingStatus {
	if !broker.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	if !encoderpc.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	return healthpb.HealthCheckResponse_SERVING
}
