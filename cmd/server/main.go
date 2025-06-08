package main

import (
	"context"
	"log"

	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	encoderpc "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/client/encode"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/publisher"
)

func main() {
	logger.Init()
	config.Init()

	ctx := context.Background()
	shutdown := jaeger.Init(ctx)

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown jaeger tracer: %v", err)
		}
	}()

	go func() {
		prometheus.Connect()
	}()

	broker.Connect()
	defer broker.Conn.Close()

	publishChan, err := broker.NewChannel()

	if err != nil {
		logger.Fatal("Unable to create publish channel %v", err)
	}

	publisher.Init(publishChan)

	encoderpc.Connect(ctx)

	server.Connect()
}
