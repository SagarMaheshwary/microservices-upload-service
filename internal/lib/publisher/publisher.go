package publisher

import (
	"context"
	"encoding/json"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
)

var P *Publisher

type Publisher struct {
	channel *amqplib.Channel
}

func (p *Publisher) Publish(queue string, message *broker.MessageType) error {
	c := config.Conf.AMQP
	q, err := p.declareQueue(queue)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.PublishTimeout)
	defer cancel()

	messageData, err := json.Marshal(&message)

	if err != nil {
		logger.Error("Unable to parse message! %v", err)

		return err
	}

	err = p.channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqplib.Publishing{
			ContentType: constant.ContentTypeJSON,
			Body:        messageData,
		},
	)

	if err != nil {
		logger.Error("Unable to publish message! %v", err)

		return err
	}

	logger.Info("Message %q Sent: %v", message.Key, message.Data)

	return nil
}

func (p *Publisher) declareQueue(queue string) (*amqplib.Queue, error) {
	q, err := p.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Error("Declare queue error %v", err)

		return nil, err
	}

	return &q, err
}

func Init(channel *amqplib.Channel) {
	P = &Publisher{channel: channel}
}
