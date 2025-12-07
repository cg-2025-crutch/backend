package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/models"
	"github.com/redis/rueidis"
)

const (
	recommendationKeyPrefix = "analytics:recommendations:"
)

// Repository интерфейс для работы с хранилищем рекомендаций
type Repository interface {
	SaveRecommendations(ctx context.Context, userUID string, recommendations *models.AnalyticsResult, ttl time.Duration) error
	GetRecommendations(ctx context.Context, userUID string) (*models.AnalyticsResult, error)
	DeleteRecommendations(ctx context.Context, userUID string) error
}

// RedisRepository реализация репозитория для Redis
type RedisRepository struct {
	client rueidis.Client
}

// NewRedisRepository создает новый экземпляр Redis репозитория
func NewRedisRepository(client rueidis.Client) Repository {
	return &RedisRepository{
		client: client,
	}
}

// SaveRecommendations сохраняет рекомендации в Redis
func (r *RedisRepository) SaveRecommendations(ctx context.Context, userUID string, recommendations *models.AnalyticsResult, ttl time.Duration) error {
	l := log.FromContext(ctx)

	key := recommendationKeyPrefix + userUID

	data, err := json.Marshal(recommendations)
	if err != nil {
		l.Errorf("Failed to marshal recommendations: %v", err)
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	cmd := r.client.B().Set().Key(key).Value(rueidis.BinaryString(data)).ExSeconds(int64(ttl.Seconds())).Build()
	err = r.client.Do(ctx, cmd).Error()
	if err != nil {
		l.Errorf("Failed to save recommendations to Redis: %v", err)
		return fmt.Errorf("failed to save recommendations to Redis: %w", err)
	}

	l.Infof("Successfully saved recommendations for user %s to Redis (TTL: %s)", userUID, ttl)
	return nil
}

// GetRecommendations получает рекомендации из Redis
func (r *RedisRepository) GetRecommendations(ctx context.Context, userUID string) (*models.AnalyticsResult, error) {
	l := log.FromContext(ctx)

	key := recommendationKeyPrefix + userUID

	cmd := r.client.B().Get().Key(key).Build()
	val, err := r.client.Do(ctx, cmd).AsBytes()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			l.Infof("No recommendations found in Redis for user %s", userUID)
			return nil, nil
		}
		l.Errorf("Failed to get recommendations from Redis: %v", err)
		return nil, fmt.Errorf("failed to get recommendations from Redis: %w", err)
	}

	var recommendations models.AnalyticsResult
	err = json.Unmarshal(val, &recommendations)
	if err != nil {
		l.Errorf("Failed to unmarshal recommendations: %v", err)
		return nil, fmt.Errorf("failed to unmarshal recommendations: %w", err)
	}

	l.Infof("Successfully retrieved recommendations for user %s from Redis", userUID)
	return &recommendations, nil
}

// DeleteRecommendations удаляет рекомендации из Redis
func (r *RedisRepository) DeleteRecommendations(ctx context.Context, userUID string) error {
	l := log.FromContext(ctx)

	key := recommendationKeyPrefix + userUID

	cmd := r.client.B().Del().Key(key).Build()
	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		l.Errorf("Failed to delete recommendations from Redis: %v", err)
		return fmt.Errorf("failed to delete recommendations from Redis: %w", err)
	}

	l.Infof("Successfully deleted recommendations for user %s from Redis", userUID)
	return nil
}
