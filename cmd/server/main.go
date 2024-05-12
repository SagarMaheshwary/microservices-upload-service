package main

import (
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	grpcsrv "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

func main() {
	log.Init()
	config.Init()

	grpcsrv.Connect()
}
