package repository

import (
	"context"

	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
	"github.com/redis/rueidis"
)

type RedisReporer interface {
	InsertSubscription(ctx context.Context, userUID string, sub models.StoredSubscription) error
	GetSubscription(ctx context.Context, userUID string) (*models.StoredSubscription, error)
}

type RedisRepo struct {
	client rueidis.Client
}

func NewRedisRepo(client rueidis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}
