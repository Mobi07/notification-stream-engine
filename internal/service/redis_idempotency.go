package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	idempotencyStateProcessing = "processing"
	idempotencyStateCompleted  = "completed"
)

type redisIdempotencyStore struct {
	client *redis.Client
}

func NewRedisIdempotencyStore(client *redis.Client) IdempotencyStore {
	return &redisIdempotencyStore{client: client}
}

func (r *redisIdempotencyStore) Acquire(ctx context.Context, key string, ttl time.Duration) (IdempotencyStatus, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == nil {
		switch value {
		case idempotencyStateCompleted:
			return IdempotencyCompleted, nil
		default:
			return IdempotencyInProgress, nil
		}
	}

	if err != redis.Nil {
		return "", err
	}

	ok, err := r.client.SetNX(ctx, key, idempotencyStateProcessing, ttl).Result()
	if err != nil {
		return "", err
	}
	if ok {
		return IdempotencyAcquired, nil
	}

	value, err = r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			ok, setErr := r.client.SetNX(ctx, key, idempotencyStateProcessing, ttl).Result()
			if setErr != nil {
				return "", setErr
			}
			if ok {
				return IdempotencyAcquired, nil
			}
			return IdempotencyInProgress, nil
		}
		return "", err
	}

	if value == idempotencyStateCompleted {
		return IdempotencyCompleted, nil
	}

	return IdempotencyInProgress, nil
}

func (r *redisIdempotencyStore) MarkCompleted(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Set(ctx, key, idempotencyStateCompleted, ttl).Err()
}

func (r *redisIdempotencyStore) Release(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
