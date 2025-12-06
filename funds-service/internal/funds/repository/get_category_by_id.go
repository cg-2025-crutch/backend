package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *FundsRepository) GetCategoryById(ctx context.Context, id int32) (*models.Category, error) {
	query := `
		SELECT id, name, type, icon, created_at
		FROM categories
		WHERE id = $1
	`

	var category models.Category
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Type,
		&category.Icon,
		&category.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}
