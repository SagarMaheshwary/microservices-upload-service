package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

var Conn *amqplib.Connection
var Channel *amqplib.Channel

type MessageType struct {
	Key  string `json:"key"`
	Data any    `json:"data"`
}

func Connect() {
	c := config.Getamqp()

	address := fmt.Sprintf("amqp://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)

	var err error

	Conn, err = amqplib.Dial(address)

	if err != nil {
		log.Error("AMQP connection error %v", err)
	}

	Channel, err = Conn.Channel()

	if err != nil {
		log.Error("AMQP channel error %v", err)
	}

	log.Info("AMQP connected on %q", address)
}

func PublishMessage(queue string, message *MessageType) error {
	q, err := Channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("AMQP queue error %v", err)

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageData, err := json.Marshal(&message)

	if err != nil {
		log.Error("Unable to parse message %v", message)

		return err
	}

	err = Channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqplib.Publishing{
			ContentType: "application/json",
			Body:        messageData,
		},
	)

	if err != nil {
		log.Error("AMQP Unable to publish message %v", err)

		return err
	}

	log.Info("Message Sent")

	return nil
}
