package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *FundsRepository) GetUserBalance(ctx context.Context, userUID string) (*models.UserBalance, error) {
	query := `
		SELECT user_uid, total_balance, total_income, total_expense, last_transaction_at, updated_at
		FROM user_balances
		WHERE user_uid = $1
	`

	var balance models.UserBalance
	err := r.db.QueryRow(ctx, query, userUID).Scan(
		&balance.UserUID,
		&balance.TotalBalance,
		&balance.TotalIncome,
		&balance.TotalExpense,
		&balance.LastTransactionAt,
		&balance.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &models.UserBalance{
				UserUID:      userUID,
				TotalBalance: 0,
				TotalIncome:  0,
				TotalExpense: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	return &balance, nil
}
