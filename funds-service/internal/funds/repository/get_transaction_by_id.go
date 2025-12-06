package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *FundsRepository) GetTransactionById(ctx context.Context, id int64) (*models.Transaction, error) {
	query := `
		SELECT t.id, t.user_uid, t.category_id, t.type, t.amount, t.title, t.description, 
		       t.transaction_date, t.created_at, t.updated_at,
		       c.id, c.name, c.type, c.icon, c.created_at
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1
	`

	var transaction models.Transaction
	var category models.Category

	err := r.db.QueryRow(ctx, query, id).Scan(
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
		&category.ID,
		&category.Name,
		&category.Type,
		&category.Icon,
		&category.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	transaction.Category = &category
	return &transaction, nil
}
