package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
	"github.com/mailru/easyjson"
)

func (r *RedisRepo) GetSubscription(ctx context.Context, userUID string) (*models.StoredSubscription, error) {
	cmd := r.client.B().Get().Key(userUID).Build()
	result, err := r.client.Do(ctx, cmd).ToString()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription from redis: %w", err)
	}

	var sub models.StoredSubscription
	err = easyjson.Unmarshal([]byte(result), &sub)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return &sub, nil
}
