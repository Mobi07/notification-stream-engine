package rabbitmq

import (
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func StartConsumer(url, queueName string) error {

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Log.Error("rabbitmq connection failed", zap.Error(err))
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Error("channel creation failed", zap.Error(err))
		return err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error("consumer start failed", zap.Error(err))
		return nil
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			logger.Log.Info(
				"message received",
				zap.ByteString("body", msg.Body),
			)
		}
	}()

	logger.Log.Info("worker started", zap.String("queue", queueName))

	<-forever
	return nil
}
