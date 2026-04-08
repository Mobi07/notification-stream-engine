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
		conn.Close() // prevent leak
		return nil, err
	}

	return &Producer{
		connection: conn,
		channel:    ch,
		queue:      amqp.Queue{Name: queueName},
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

func (p *Producer) Close() {
	if err := p.channel.Close(); err != nil {
		logger.Log.Error("failed to close channel", zap.Error(err))
	}
	if err := p.connection.Close(); err != nil {
		logger.Log.Error("failed to close connection", zap.Error(err))
	}
}