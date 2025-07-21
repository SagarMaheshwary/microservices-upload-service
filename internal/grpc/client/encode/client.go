package encode

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func NewClient(ctx context.Context) error {
	var opts []grpc.DialOption

	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)
	address := config.Conf.GRPCClient.EncodeServiceURL

	connection, err := grpc.NewClient(address, opts...)
	if err != nil {
		logger.Error("Encode gRPC failed to connect on %q: %v", address, err)
		return err
	}

	Encode = &encodeClient{
		health: healthpb.NewHealthClient(connection),
	}

	if err := HealthCheck(ctx); err != nil {
		return err
	}

	logger.Info("Encode gRPC client connected on %q", address)

	return nil
}

func HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := Encode.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		logger.Error("Encode gRPC health check failed! %v", err)
		return err
	}

	if response.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("Encode gRPC health check failed!")
		return errors.New("Encode gRPC health check failed")
	}

	return nil
}
