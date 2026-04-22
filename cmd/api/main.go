package main

import (
	"github.com/Mobi07/notification-stream-engine.git/config"
	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/internal/handler"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg, err := config.Load("../../config/config.yaml")
	if err != nil {
		logger.Log.Fatal("failed to load config", zap.Error(err))
	}

	producer, err := rabbitmq.NewProducer(cfg.RabbitMQ.URL, cfg.RabbitMQ.MainQueue)
	if err != nil {
		logger.Log.Fatal("failed to create producer", zap.Error(err))
	}
	defer producer.Close()

	r := gin.Default()

	eventHandler := handler.NewEventHandler(producer)

	r.POST("/events", eventHandler.PublishEvent)

	logger.Log.Info("API server running", zap.String("port", cfg.API.Port))

	if err := r.Run(cfg.API.Port); err != nil {
		logger.Log.Fatal("API server failed", zap.Error(err))
	}
}
