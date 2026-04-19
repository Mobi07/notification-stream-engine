package rabbitmq

import (
	"time"

	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type QueueConfig struct {
	MainQueue   string
	RetryQueue  string
	DLQ         string
	DLXExchange string
	RetryDelay  time.Duration
}

func SetUpQueues(ch *amqp.Channel, cfg QueueConfig) error {

	err := ch.ExchangeDeclare(
		cfg.DLXExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error("failed to declare exchange")
		return err
	}

	mainArgs := amqp.Table{
		"x-dead-letter-exchange":    cfg.DLXExchange,
		"x-dead-letter-routing-key": cfg.RetryQueue,
	}

	_, err = ch.QueueDeclare(
		cfg.MainQueue,
		true,
		false,
		false,
		false,
		mainArgs,
	)
	if err != nil {
		logger.Log.Error("failed to declare main queue", zap.Error(err))
		return err
	}

	_, err = ch.QueueDeclare(
		cfg.DLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error("failed to declare DLQ", zap.Error(err))
		return err
	}

	retryArgs := amqp.Table{
		"x-message-ttl":             int32(cfg.RetryDelay / time.Millisecond),
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": cfg.MainQueue,
	}

	_, err = ch.QueueDeclare(
		cfg.RetryQueue,
		true,
		false,
		false,
		false,
		retryArgs,
	)
	if err != nil {
		logger.Log.Error("failed to declare retry queue", zap.Error(err))
		return err
	}

	err = ch.QueueBind(
		cfg.RetryQueue,
		cfg.RetryQueue,
		cfg.DLXExchange,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error("failed to Bind Retry Queue", zap.Error(err))
		return err
	}

	err = ch.QueueBind(
		cfg.DLQ,
		cfg.DLQ,
		cfg.DLXExchange,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error("failed to Bind Queue", zap.Error(err))
		return err
	}

	logger.Log.Info("queues declared successfully")
	return nil
}
