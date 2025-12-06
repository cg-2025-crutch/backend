package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetUserBalance(ctx context.Context, userUID string) (*models.UserBalance, error) {
	return s.repo.GetUserBalance(ctx, userUID)
}
