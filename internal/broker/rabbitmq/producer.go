package rabbitmq

import (
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewProducer(url, queueName string) (*Producer, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (p *Producer) Publish(body []byte) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logger.Log.Error("message publish failed", zap.Error(err))
		return err
	}

	logger.Log.Info("event published", zap.String("queue", p.queue.Name))
	return nil
}
