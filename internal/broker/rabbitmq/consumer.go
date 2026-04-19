package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Mobi07/notification-stream-engine.git/internal/constants"
	"github.com/Mobi07/notification-stream-engine.git/internal/events"
	"github.com/Mobi07/notification-stream-engine.git/internal/service"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func StartConsumer(ch *amqp.Channel, queueName string, svc service.NotificationService) error {
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

			// Direct DLQ
			if err := PublishToDLQ(ch, msg); err != nil {
				logger.Log.Error("failed to publish to DLQ", zap.Error(err))
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
			continue
		}

		retryCount := GetRetryCount(msg)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err := svc.ProcessEvent(ctx, event)
		cancel()

		if err != nil {
			logger.Log.Error("event processing failed", zap.String("event_type", event.Type), zap.Error(err))

			if retryCount >= constants.MaxRetryCount {
				logger.Log.Warn("max retry exceeded, sending to DLQ", zap.Int("retry_count", retryCount))
				if err := PublishToDLQ(ch, msg); err != nil {
					logger.Log.Error("failed to publish to DLQ", zap.Error(err))
					continue
				}
				msg.Ack(false)
				continue
			}

			if err := PublishToRetryQueue(ch, msg, retryCount+1); err != nil {
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

func StartDLQConsumer(ch *amqp.Channel) error {

	msgs, err := ch.Consume(
		"notification_dlq",
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
