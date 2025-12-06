package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	return s.repo.GetAllCategories(ctx)
}
