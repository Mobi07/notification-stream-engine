package policy

import (
	"errors"
	"time"

	appErrors "github.com/Mobi07/notification-stream-engine.git/internal/errors"
)

type Decision string

const (
	DecisionAck   Decision = "ack"
	DecisionRetry Decision = "retry"
	DecisionDLQ   Decision = "dlq"
)

type ProcessingConfig struct {
	MaxRetryCount            int
	RetryDelay               time.Duration
	ConsumerTimeout          time.Duration
	IdempotencyProcessingTTL time.Duration
	IdempotencyCompletedTTL  time.Duration
}

func Decide(err error, retryCount int, cfg ProcessingConfig) Decision {
	var appErr appErrors.AppError
	if errors.As(err, &appErr) && !appErr.Retryable {
		return DecisionDLQ
	}

	if retryCount >= cfg.MaxRetryCount {
		return DecisionDLQ
	}

	return DecisionRetry
}
