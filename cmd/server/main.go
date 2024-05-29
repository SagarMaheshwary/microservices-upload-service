package main

import (
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	grpcsrv "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/publisher"
)

func main() {
	log.Init()
	config.Init()

	broker.Connect()
	defer broker.Conn.Close()

	publishChan, err := broker.NewChannel()

	if err != nil {
		log.Fatal("Unable to create publish channel %v", err)
	}

	publisher.Init(publishChan)

	grpcsrv.Connect()
}
