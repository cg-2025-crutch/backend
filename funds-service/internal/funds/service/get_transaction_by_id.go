package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) GetTransactionById(ctx context.Context, id int64) (*models.Transaction, error) {
	return s.repo.GetTransactionById(ctx, id)
}
