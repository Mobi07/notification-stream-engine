package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mobi07/notification-stream-engine.git/internal/errors"
	"github.com/Mobi07/notification-stream-engine.git/internal/events"
	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"go.uber.org/zap"
)

const (
	idempotencyProcessingTTL = 30 * time.Second
	idempotencyCompletedTTL  = 24 * time.Hour
)

type EventHandler interface {
	Handle(ctx context.Context, payload interface{}) error
}

type NotificationService interface {
	ProcessEvent(ctx context.Context, event events.Event) error
}

type notificationService struct {
	handlers         map[string]EventHandler
	idempotencyStore IdempotencyStore
}

func NewNotificationService(handlers map[string]EventHandler, idempotencyStore IdempotencyStore) NotificationService {
	return &notificationService{
		handlers:         handlers,
		idempotencyStore: idempotencyStore,
	}
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

	if event.ID == "" {
		return errors.AppError{
			Err:       fmt.Errorf("missing event id"),
			Retryable: false,
		}
	}

	idempotencyKey := fmt.Sprintf("idempotency:event:%s", event.ID)
	status, err := s.idempotencyStore.Acquire(ctx, idempotencyKey, idempotencyProcessingTTL)
	if err != nil {
		return errors.AppError{
			Err:       fmt.Errorf("failed to acquire idempotency key: %w", err),
			Retryable: true,
		}
	}

	if status == IdempotencyCompleted || status == IdempotencyInProgress {
		logger.Log.Info(
			"duplicate event skipped",
			zap.String("event_id", event.ID),
			zap.String("event_type", event.Type),
			zap.String("idempotency_status", string(status)),
		)
		return nil
	}

	err = handler.Handle(ctx, event.Payload)
	if err != nil {
		if releaseErr := s.idempotencyStore.Release(ctx, idempotencyKey); releaseErr != nil {
			logger.Log.Error(
				"failed to release idempotency key after handler failure",
				zap.String("event_id", event.ID),
				zap.Error(releaseErr),
			)
		}
		return err
	}

	if err := s.idempotencyStore.MarkCompleted(ctx, idempotencyKey, idempotencyCompletedTTL); err != nil {
		return errors.AppError{
			Err:       fmt.Errorf("failed to mark event as completed: %w", err),
			Retryable: true,
		}
	}

	return nil
}
