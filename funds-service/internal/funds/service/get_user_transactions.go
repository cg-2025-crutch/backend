package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetUserTransactions(ctx context.Context, userUID string, limit, offset int32) ([]*models.Transaction, int64, error) {
	return s.repo.GetUserTransactions(ctx, userUID, limit, offset)
}
