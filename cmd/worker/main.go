package main

import (
	"time"

	"github.com/Mobi07/notification-stream-engine.git/config"
	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/internal/delivery"
	"github.com/Mobi07/notification-stream-engine.git/internal/policy"
	"github.com/Mobi07/notification-stream-engine.git/internal/service"
	"github.com/Mobi07/notification-stream-engine.git/internal/service/handlers"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg, err := config.Load("../../config/config.yaml")
	if err != nil {
		logger.Log.Fatal("failed to load config", zap.Error(err))
	}

	logger.Log.Info("Worker service started")

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Log.Fatal("failed to connect rabbitmq", zap.Error(err))
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Log.Error("failed to close connection ", zap.Error(err))
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Fatal("failed to open channel", zap.Error(err))
	}
	defer func() {
		if err := ch.Close(); err != nil {
			logger.Log.Fatal("failed to close channel", zap.Error(err))
		}
	}()

	queueConfig := rabbitmq.QueueConfig{
		MainQueue:   cfg.RabbitMQ.MainQueue,
		RetryQueue:  cfg.RabbitMQ.RetryQueue,
		DLQ:         cfg.RabbitMQ.DLQ,
		DLXExchange: cfg.RabbitMQ.DLXExchange,
		RetryDelay:  time.Duration(cfg.Processing.RetryDelayMilliseconds) * time.Millisecond,
	}

	processingConfig := policy.ProcessingConfig{
		MaxRetryCount:            cfg.Processing.MaxRetryCount,
		RetryDelay:               time.Duration(cfg.Processing.RetryDelayMilliseconds) * time.Millisecond,
		ConsumerTimeout:          time.Duration(cfg.Processing.ConsumerTimeoutSeconds) * time.Second,
		IdempotencyProcessingTTL: time.Duration(cfg.Processing.IdempotencyProcessingTTL) * time.Second,
		IdempotencyCompletedTTL:  time.Duration(cfg.Processing.IdempotencyCompletedTTL) * time.Hour,
	}

	if err = rabbitmq.SetUpQueues(ch, queueConfig); err != nil {
		logger.Log.Fatal("failed to start consumer", zap.Error(err))
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	rateLimiter := service.NewRedisRateLimiter(redisClient)
	idempotencyStore := service.NewRedisIdempotencyStore(redisClient)
	emailSender := delivery.NewEmailSender()
	eventHandlers := map[string]service.EventHandler{
		"UserRegistration": handlers.NewUserRegistrationHandler(emailSender, rateLimiter),
	}

	notificationService := service.NewNotificationService(eventHandlers, idempotencyStore, processingConfig)

	if err := rabbitmq.StartDLQConsumer(ch, queueConfig.DLQ); err != nil {
		logger.Log.Fatal("failed to start DLQ consumer", zap.Error(err))
	}

	if err = rabbitmq.StartConsumer(ch, queueConfig.MainQueue, queueConfig, processingConfig, notificationService); err != nil {
		logger.Log.Fatal("failed to start consumer", zap.Error(err))
	}

}
