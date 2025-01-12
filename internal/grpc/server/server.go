package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	uploadpb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect() {
	c := config.Conf.GRPCServer

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	listener, err := net.Listen(constant.ProtocolTCP, address)

	if err != nil {
		logger.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var options []grpc.ServerOption

	server := grpc.NewServer(options...)

	uploadpb.RegisterUploadServiceServer(server, &uploadServer{})
	healthpb.RegisterHealthServer(server, &healthServer{})

	logger.Info("gRPC server started on %q", address)

	if err := server.Serve(listener); err != nil {
		logger.Error("gRPC server failed to start %v", err)
	}
}
