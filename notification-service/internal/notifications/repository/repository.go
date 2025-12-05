package repository

import (
	"context"

	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
)

type RedisRepo interface {
	InsertSubscription(ctx context.Context, userUID string, sub models.StoredSubscription) error
	GetSubscription(ctx context.Context, userUID string) (*models.StoredSubscription, error)
}
