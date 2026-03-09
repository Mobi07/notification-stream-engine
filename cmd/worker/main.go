package main

import (
	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	logger.Log.Info("Worker service started")
	err := rabbitmq.StartConsumer("amqp://guest:guest@localhost:5672/", "events_queue")
	if err != nil {
		logger.Log.Fatal("failed to start consumer", zap.Error(err))
	}
}
