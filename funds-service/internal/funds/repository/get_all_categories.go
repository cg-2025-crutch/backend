package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

func (r *FundsRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	query := `
		SELECT id, name, type, icon, created_at
		FROM categories
		ORDER BY type, name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Type,
			&category.Icon,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}
