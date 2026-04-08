package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mobi07/notification-stream-engine.git/internal/broker/rabbitmq"
	"github.com/Mobi07/notification-stream-engine.git/internal/events"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EventHandler struct {
	producer *rabbitmq.Producer
}

func NewEventHandler(p *rabbitmq.Producer) *EventHandler {
	return &EventHandler{producer: p}
}

func (h *EventHandler) PublishEvent(c *gin.Context) {
	var event events.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		logger.Log.Error("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	event.ID = uuid.New().String()
	event.Timestamp = time.Now().Unix()

	body, err := json.Marshal(event)
	if err != nil {
		logger.Log.Error("failed to marshal event", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	err = h.producer.Publish(body)
	if err != nil {
		logger.Log.Error("failed to publish event", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish"})
		return
	}

	logger.Log.Info("event published", zap.String("event_id", event.ID), zap.String("type", event.Type))
	c.JSON(http.StatusOK, gin.H{"status": "event published", "event_id": event.ID})

}
