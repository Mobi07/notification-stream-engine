package service

import (
	"context"
	"time"
)

type IdempotencyStore interface {
	CheckAndSet(ctx context.Context, key string, ttl time.Duration) (bool, error)
}