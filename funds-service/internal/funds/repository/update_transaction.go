package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *FundsRepository) UpdateTransaction(ctx context.Context, input models.UpdateTransactionInput) (*models.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldTransaction models.Transaction
	getQuery := `
		SELECT user_uid, type, amount 
		FROM transactions 
		WHERE id = $1
	`
	err = tx.QueryRow(ctx, getQuery, input.ID).Scan(
		&oldTransaction.UserUID,
		&oldTransaction.Type,
		&oldTransaction.Amount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get old transaction: %w", err)
	}

	// Verify that the transaction belongs to the user
	if oldTransaction.UserUID != input.UserUID {
		return nil, fmt.Errorf("transaction does not belong to user")
	}

	if err := r.updateUserBalanceInTx(ctx, tx, oldTransaction.UserUID, oldTransaction.Type, oldTransaction.Amount, false); err != nil {
		return nil, fmt.Errorf("failed to revert old balance: %w", err)
	}

	updateQuery := `
		UPDATE transactions 
		SET category_id = $1, type = $2, amount = $3, title = $4, description = $5, 
		    transaction_date = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, user_uid, category_id, type, amount, title, description, transaction_date, created_at, updated_at
	`

	var transaction models.Transaction
	err = tx.QueryRow(ctx, updateQuery,
		input.CategoryID,
		input.Type,
		input.Amount,
		input.Title,
		input.Description,
		input.TransactionDate,
		input.ID,
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
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	if err := r.updateUserBalanceInTx(ctx, tx, transaction.UserUID, transaction.Type, transaction.Amount, true); err != nil {
		return nil, fmt.Errorf("failed to apply new balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	category, err := r.GetCategoryById(ctx, input.CategoryID)
	if err == nil {
		transaction.Category = category
	}

	return &transaction, nil
}
