package main

import (
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	grpcsrv "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/amqp"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

func main() {
	log.Init()
	config.Init()

	amqp.Connect()
	defer amqp.Channel.Close()
	defer amqp.Conn.Close()

	grpcsrv.Connect()
}
