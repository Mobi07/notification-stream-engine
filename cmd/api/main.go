package main

import (
	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	producer, err := rabbitmq.NewProducer("amqp://guest:guest@localhost:5672/", "events_queue")
	if err != nil {
		logger.Log.Fatal("failed to create producer", zap.Error(err))
	}

	err = producer.Publish([]byte(`{"id":"1","type":"user_signup","timestamp":1627847261,"payload":{"user_id":"123"}}`))
	if err != nil {
		logger.Log.Error("failed to publish event", zap.Error(err))
		return
	}

	logger.Log.Info("API service started")
}
