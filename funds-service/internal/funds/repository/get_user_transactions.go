package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (r *FundsRepository) GetUserTransactions(ctx context.Context, userUID string, limit, offset int32) ([]*models.Transaction, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_uid = $1`
	err := r.db.QueryRow(ctx, countQuery, userUID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions count: %w", err)
	}

	query := `
		SELECT t.id, t.user_uid, t.category_id, t.type, t.amount, t.title, t.description, 
		       t.transaction_date, t.created_at, t.updated_at,
		       c.id, c.name, c.type, c.icon, c.created_at
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_uid = $1
		ORDER BY t.transaction_date DESC, t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userUID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]*models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		var category models.Category

		err := rows.Scan(
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
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}

		transaction.Category = &category
		transactions = append(transactions, &transaction)
	}

	return transactions, total, nil
}
