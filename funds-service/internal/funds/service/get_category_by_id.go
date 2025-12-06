package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetCategoryById(ctx context.Context, id int32) (*models.Category, error) {
	return s.repo.GetCategoryById(ctx, id)
}
