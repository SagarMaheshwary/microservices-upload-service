package main

import (
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	encoderpc "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/client/encode"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/publisher"
)

func main() {
	logger.Init()
	config.Init()

	broker.Connect()
	defer broker.Conn.Close()

	publishChan, err := broker.NewChannel()

	if err != nil {
		logger.Fatal("Unable to create publish channel %v", err)
	}

	publisher.Init(publishChan)

	encoderpc.Connect()
	server.Connect()
}
