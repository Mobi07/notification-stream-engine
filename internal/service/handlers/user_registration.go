package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/Mobi07/notification-stream-engine.git/internal/delivery"
	"github.com/Mobi07/notification-stream-engine.git/internal/errors"
	"github.com/Mobi07/notification-stream-engine.git/internal/events"
	"github.com/Mobi07/notification-stream-engine.git/internal/service"
	"github.com/go-viper/mapstructure/v2"
)

type UserRegistrationHandler struct {
	emailSender delivery.EmailSender
	rateLimiter service.RateLimiter
}

func NewUserRegistrationHandler(emailSender delivery.EmailSender, rateLimiter service.RateLimiter) *UserRegistrationHandler {
	return &UserRegistrationHandler{
		emailSender: emailSender,
		rateLimiter: rateLimiter,
	}
}

func (h *UserRegistrationHandler) Handle(ctx context.Context, payload interface{}) error {
	var data events.UserRegistration

	if err := mapstructure.Decode(payload, &data); err != nil {
		return errors.AppError{
			Err:       err,
			Retryable: false,
		}
	}

	key := fmt.Sprintf("rate_limit:%s:usr_reg", data.UserID)

	allowed, err := h.rateLimiter.Allow(ctx, key, 5, time.Minute)
	if err != nil {
		return errors.AppError{Err: err, Retryable: true}
	}

	if !allowed {
		return errors.AppError{Err: fmt.Errorf("rate limit exceeded"), Retryable: false}
	}

	err = h.emailSender.SendWelcomeEmail(ctx, data.Email)
	if err != nil {
		return errors.AppError{Err: err, Retryable: true}
	}

	return nil
}
