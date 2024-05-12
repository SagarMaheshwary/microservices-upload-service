package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"google.golang.org/grpc"
)

func Connect() {
	c := config.GetgrpcServer()

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterUploadServer(grpcServer, &uploadServer{})

	log.Info("gRPC server started on %q", address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Error("gRPC server failed to start %v", err)
	}
}
