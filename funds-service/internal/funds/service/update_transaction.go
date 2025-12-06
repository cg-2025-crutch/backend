package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (s *FundsService) UpdateTransaction(ctx context.Context, input models.UpdateTransactionInput) (*models.Transaction, error) {
	err := s.prod.Produce(ctx, []byte(input.UserUID), []byte("update"))
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateTransaction(ctx, input)
}
