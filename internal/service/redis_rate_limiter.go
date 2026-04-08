package service

import (
	"context"
	"time"

	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) RateLimiter {
	return &redisRateLimiter{client: client}
}

func (r *redisRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		logger.Log.Error("failed to increment rate limiter key", zap.String("key", key), zap.Error(err))
		return false, err
	}

	if count == 1 {
		err = r.client.Expire(ctx, key, window).Err()
		if err != nil {
			logger.Log.Error("failed to set expiration for rate limiter key", zap.String("key", key), zap.Error(err))
			return false, err
		}
	}

	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}
