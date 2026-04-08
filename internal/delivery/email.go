package delivery

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mobi07/notification-stream-engine.git/pkg/logger"
	"go.uber.org/zap"
)

type EmailSender interface {
	SendWelcomeEmail(ctx context.Context, email string) error
}

type emailSender struct{}

func NewEmailSender() EmailSender {
	return &emailSender{}
}

func (e *emailSender) SendWelcomeEmail(ctx context.Context, email string) error {

	if strings.Contains(email, "fail") {
		return fmt.Errorf("simulated email failure")
	}

	logger.Log.Info("welcome email sent", zap.String("email", email))
	return nil
}
