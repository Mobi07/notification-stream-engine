package main

import (
	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/internal/constants"
	"github.com/Mobi07/notification-stream-engine.git/internal/handler"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	producer, err := rabbitmq.NewProducer("amqp://guest:guest@localhost:5672/", constants.MainQueueName)
	if err != nil {
		logger.Log.Fatal("failed to create producer", zap.Error(err))
	}

	r := gin.Default()

	eventHandler := handler.NewEventHandler(producer)

	r.POST("/events", eventHandler.PublishEvent)

	logger.Log.Info("API server running on :8080")

	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("API server failed", zap.Error(err))
	}
}
