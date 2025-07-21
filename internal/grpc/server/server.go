package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/prometheus"
	uploadpb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func NewServer() *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(prometheusUnaryInterceptor),
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)

	uploadpb.RegisterUploadServiceServer(server, &uploadServer{
		s3Storage: aws.NewS3Service(),
	})
	healthpb.RegisterHealthServer(server, &healthServer{})

	return server
}

func Serve(server *grpc.Server) error {
	c := config.Conf.GRPCServer
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	listener, err := net.Listen(constant.ProtocolTCP, address)
	if err != nil {
		logger.Fatal("Failed to create tcp listner on %q: %v", address, err)
		return err
	}

	logger.Info("gRPC server started on %q", address)

	if err := server.Serve(listener); err != nil {
		logger.Error("gRPC server failed to start %v", err)
		return err
	}

	return nil
}

func prometheusUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	response, err := handler(ctx, req)

	method := info.FullMethod
	statusCode := status.Code(err).String()
	prometheus.GRPCRequestCounter.WithLabelValues(method, statusCode).Inc()
	prometheus.GRPCRequestLatency.WithLabelValues(method).Observe(time.Since(start).Seconds())

	return response, err
}
