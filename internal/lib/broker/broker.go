package broker

import (
	"fmt"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
)

var Conn *amqplib.Connection

type MessageType struct {
	Key  string `json:"key"`
	Data any    `json:"data"`
}

func Connect() {
	c := config.Conf.AMQP

	address := fmt.Sprintf("%s://%s:%s@%s:%d", constant.ProtocolAMQP, c.Username, c.Password, c.Host, c.Port)

	var err error

	Conn, err = amqplib.Dial(address)

	if err != nil {
		logger.Fatal("Broker connection error %v", err)
	}

	logger.Info("Broker connected on %q", address)
}

func NewChannel() (*amqplib.Channel, error) {
	c, err := Conn.Channel()

	if err != nil {
		logger.Error("Broker channel error %v", err)

		return nil, err
	}

	return c, nil
}

func HealthCheck() bool {
	if Conn == nil || Conn.IsClosed() {
		logger.Info("AMQP health check failed!")

		return false
	}

	return true
}
