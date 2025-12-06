package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *FundsRepository) DeleteTransaction(ctx context.Context, id int64, userUID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var transaction models.Transaction
	getQuery := `
		SELECT user_uid, type, amount 
		FROM transactions 
		WHERE id = $1
	`
	err = tx.QueryRow(ctx, getQuery, id).Scan(
		&transaction.UserUID,
		&transaction.Type,
		&transaction.Amount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("transaction not found")
		}
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Verify that the transaction belongs to the user
	if transaction.UserUID != userUID {
		return fmt.Errorf("transaction does not belong to user")
	}

	deleteQuery := `DELETE FROM transactions WHERE id = $1`
	_, err = tx.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	if err := r.updateUserBalanceInTx(ctx, tx, transaction.UserUID, transaction.Type, transaction.Amount, false); err != nil {
		return fmt.Errorf("failed to revert balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
