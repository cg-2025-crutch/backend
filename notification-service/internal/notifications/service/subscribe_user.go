package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
)

func (s *NotificationService) SubscribeUser(ctx context.Context, userId, endpoint, p256dh, auth string) error {
	l := log.FromContext(ctx)
	l.Info(userId, endpoint, p256dh, auth)

	sub := models.StoredSubscription{
		Endpoint:  endpoint,
		P256dh:    p256dh,
		Auth:      auth,
		CreatedAt: time.Now().UTC(),
	}

	err := s.repo.InsertSubscription(ctx, userId, sub)
	if err != nil {
		l.Errorf("failed to insert subscription: %s", err)
		return fmt.Errorf("failed to insert subscription: %s", err)
	}

	return nil
}
