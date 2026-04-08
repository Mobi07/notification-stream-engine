package main

import (
	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/internal/constants"
	"github.com/Mobi07/notification-stream-engine.git/internal/delivery"
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

	logger.Log.Info("Worker service started")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	if err = rabbitmq.SetUpQueues(ch); err != nil {
		logger.Log.Fatal("failed to start consumer", zap.Error(err))
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rateLimiter := service.NewRedisRateLimiter(redisClient)
	emailSender := delivery.NewEmailSender()
	eventHandlers := map[string]service.EventHandler{
		"UserRegistration": handlers.NewUserRegistrationHandler(emailSender, rateLimiter),
	}

	notificationService := service.NewNotificationService(eventHandlers)

	if err := rabbitmq.StartDLQConsumer(ch); err != nil {
		logger.Log.Fatal("failed to start DLQ consumer", zap.Error(err))
	}

	if err = rabbitmq.StartConsumer(ch, constants.MainQueueName, notificationService); err != nil {
		logger.Log.Fatal("failed to start consumer", zap.Error(err))
	}

}
