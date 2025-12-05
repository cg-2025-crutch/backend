package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
	"github.com/mailru/easyjson"
)

func (s *NotificationService) SendNotification(ctx context.Context, userUID string, not models.Notification) error {
	l := log.FromContext(ctx)

	payload, err := easyjson.Marshal(not)
	if err != nil {
		l.Error("failed to marshal notification", err)
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	key := fmt.Sprintf("webpush:%s", userUID)

	sub, err := s.repo.GetSubscription(ctx, key)
	if err != nil {
		l.Error("failed to fetch subscription from redis", err)
		return fmt.Errorf("failed to fetch subscription from redis: %w", err)
	}
	err = s.sendToSubscriber(ctx, *sub, payload)
	if err != nil {
		l.Error("failed to send notification to subscriber", err)
		return fmt.Errorf("failed to send notification to subscriber: %w", err)
	}

	return nil
}

func (s *NotificationService) sendToSubscriber(ctx context.Context, sb models.StoredSubscription, payload []byte) error {
	l := log.FromContext(ctx)

	l.Infof("stored subscriber endpoint: %s", sb.Endpoint)
	l.Infof("stored subscriber p256 length: %d, value: %s", len(sb.P256dh), sb.P256dh)
	l.Infof("stored subscriber auth length: %d, value: %s", len(sb.Auth), sb.Auth)

	// Validate base64url keys
	if _, err := base64.RawURLEncoding.DecodeString(sb.P256dh); err != nil {
		l.Errorf("Invalid p256dh key (not valid base64url): %v", err)
		return fmt.Errorf("invalid p256dh key: %w", err)
	}

	if _, err := base64.RawURLEncoding.DecodeString(sb.Auth); err != nil {
		l.Errorf("Invalid auth key (not valid base64url): %v", err)
		return fmt.Errorf("invalid auth key: %w", err)
	}

	sub := &webpush.Subscription{
		Endpoint: sb.Endpoint,
		Keys: webpush.Keys{
			Auth:   sb.Auth,
			P256dh: sb.P256dh,
		},
	}

	opts := &webpush.Options{
		Subscriber:      s.subscriber,
		VAPIDPublicKey:  s.VapidPublicKey,
		VAPIDPrivateKey: s.VapidPrivateKey,
		TTL:             60,
	}

	resp, err := webpush.SendNotification(payload, sub, opts)
	if err != nil {
		l.Errorf("webpush send failed: %v", err)
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	l.Infof("Notification sent successfully: status=%d", resp.StatusCode)

	return nil
}
