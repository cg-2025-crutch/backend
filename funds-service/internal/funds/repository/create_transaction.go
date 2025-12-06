package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (r *FundsRepository) CreateTransaction(ctx context.Context, input models.CreateTransactionInput) (*models.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO transactions (user_uid, category_id, type, amount, title, description, transaction_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_uid, category_id, type, amount, title, description, transaction_date, created_at, updated_at
	`

	var transaction models.Transaction
	err = tx.QueryRow(ctx, query,
		input.UserUID,
		input.CategoryID,
		input.Type,
		input.Amount,
		input.Title,
		input.Description,
		input.TransactionDate,
	).Scan(
		&transaction.ID,
		&transaction.UserUID,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Title,
		&transaction.Description,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := r.updateUserBalanceInTx(ctx, tx, input.UserUID, input.Type, input.Amount, true); err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Get category details
	category, err := r.GetCategoryById(ctx, input.CategoryID)
	if err == nil {
		transaction.Category = category
	}

	return &transaction, nil
}
