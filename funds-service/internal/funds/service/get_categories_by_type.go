package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetCategoriesByType(ctx context.Context, categoryType string) ([]*models.Category, error) {
	return s.repo.GetCategoriesByType(ctx, categoryType)
}
