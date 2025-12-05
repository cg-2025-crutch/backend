package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
	"github.com/redis/rueidis"
)

func (r *RedisRepo) InsertSubscription(ctx context.Context, userUID string, sub models.StoredSubscription) error {
	key := fmt.Sprintf("webpush:%s", userUID)

	// Use standard json instead of easyjson to avoid base64url encoding issues
	data, err := json.Marshal(sub)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}

	err = r.client.Do(ctx, r.client.B().Set().Key(key).Value(rueidis.BinaryString(data)).Build()).Error()
	if err != nil {
		return fmt.Errorf("failed to set subscription in redis: %w", err)
	}

	return nil
}
