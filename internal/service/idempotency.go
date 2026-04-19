package service

import (
	"context"
	"time"
)

type IdempotencyStatus string

const (
	IdempotencyAcquired   IdempotencyStatus = "acquired"
	IdempotencyCompleted  IdempotencyStatus = "completed"
	IdempotencyInProgress IdempotencyStatus = "in_progress"
)

type IdempotencyStore interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (IdempotencyStatus, error)
	MarkCompleted(ctx context.Context, key string, ttl time.Duration) error
	Release(ctx context.Context, key string) error
}
