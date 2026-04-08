package service

import (
	"context"
	"fmt"

	"github.com/Mobi07/notification-stream-engine.git/internal/errors"
	"github.com/Mobi07/notification-stream-engine.git/internal/events"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"go.uber.org/zap"
)

type EventHandler interface {
	Handle(ctx context.Context, payload interface{}) error
}

type NotificationService interface {
	ProcessEvent(ctx context.Context, event events.Event) error
}

type notificationService struct {
	handlers map[string]EventHandler
}

func NewNotificationService(handlers map[string]EventHandler) NotificationService {
	return &notificationService{handlers: handlers}
}

func (s *notificationService) ProcessEvent(ctx context.Context, event events.Event) error {
	handler, ok := s.handlers[event.Type]
	if !ok {
		logger.Log.Warn("no handler for event type", zap.String("type", event.Type))
		return errors.AppError{
			Err:       fmt.Errorf("unsupported event type: %s", event.Type),
			Retryable: false,
		}
	}

	return handler.Handle(ctx, event.Payload)
}
