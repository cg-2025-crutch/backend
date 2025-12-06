package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetUserTransactionsByPeriod(ctx context.Context, userUID string, days int32, limit, offset int32) ([]*models.Transaction, int64, error) {
	return s.repo.GetUserTransactionsByPeriod(ctx, userUID, days, limit, offset)
}
