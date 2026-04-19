package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/Mobi07/notification-stream-engine.git/internal/events"
	"github.com/Mobi07/notification-stream-engine.git/internal/policy"
	"github.com/Mobi07/notification-stream-engine.git/internal/service"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func StartConsumer(ch *amqp.Channel, queueName string, queueConfig QueueConfig, processingConfig policy.ProcessingConfig, svc service.NotificationService) error {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error("consumer start failed", zap.Error(err))
		return err
	}

	logger.Log.Info("worker started", zap.String("queue", queueName))

	for msg := range msgs {

		logger.Log.Info("StartConsumer: message received", zap.ByteString("body", msg.Body))

		var event events.Event

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			logger.Log.Error("failed to unmarshal event", zap.Error(err))

			if err := PublishToDLQ(ch, queueConfig.DLXExchange, queueConfig.DLQ, msg); err != nil {
				logger.Log.Error("failed to publish to DLQ", zap.Error(err))
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
			continue
		}

		retryCount := GetRetryCount(msg)

		ctx, cancel := context.WithTimeout(context.Background(), processingConfig.ConsumerTimeout)

		err := svc.ProcessEvent(ctx, event)
		cancel()

		if err != nil {
			logger.Log.Error("event processing failed", zap.String("event_type", event.Type), zap.Error(err))

			decision := policy.Decide(err, retryCount, processingConfig)
			if decision == policy.DecisionDLQ {
				logger.Log.Warn("message routed to DLQ", zap.String("event_type", event.Type), zap.Int("retry_count", retryCount))
				if err := PublishToDLQ(ch, queueConfig.DLXExchange, queueConfig.DLQ, msg); err != nil {
					logger.Log.Error("failed to publish to DLQ", zap.Error(err))
					msg.Nack(false, true)
					continue
				}
				msg.Ack(false)
				continue
			}

			if err := PublishToRetryQueue(ch, queueConfig.RetryQueue, msg, retryCount+1); err != nil {
				logger.Log.Error("failed to publish to target", zap.Error(err))
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
			continue
		}

		if err := msg.Ack(false); err != nil {
			logger.Log.Error("failed to ack message", zap.Error(err))
		}
	}

	return nil
}

func StartDLQConsumer(ch *amqp.Channel, dlqName string) error {

	msgs, err := ch.Consume(
		dlqName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			logger.Log.Info("StartDLQConsumer: message received in DLQ", zap.String("body", string(msg.Body)))
			msg.Ack(false)
		}
	}()

	return nil
}
