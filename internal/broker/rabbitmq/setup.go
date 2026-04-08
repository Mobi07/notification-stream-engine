package rabbitmq

import (
	"github.com/Mobi07/notification-stream-engine.git/internal/constants"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func SetUpQueues(ch *amqp.Channel) error {

	err := ch.ExchangeDeclare(
		constants.DLXExchange,
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
		"x-dead-letter-exchange":    constants.DLXExchange,
		"x-dead-letter-routing-key": constants.RetryQueueName,
	}

	_, err = ch.QueueDeclare(
		constants.MainQueueName,
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
		constants.DLQName,
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
		"x-message-ttl":             int32(5000), // 5 sec delay
		"x-dead-letter-exchange":    "",          // default exchange
		"x-dead-letter-routing-key": constants.MainQueueName,
	}

	_, err = ch.QueueDeclare(
		constants.RetryQueueName,
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
		constants.RetryQueueName,
		constants.RetryQueueName,
		constants.DLXExchange,
		false,
		nil,
	) 
	if err != nil {
		logger.Log.Error("failed to Bind Retry Queue", zap.Error(err))
		return err
	}

	err = ch.QueueBind(
		constants.DLQName,
		constants.DLQName,
		constants.DLXExchange,
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
