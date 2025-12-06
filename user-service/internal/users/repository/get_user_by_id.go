package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *postgresRepository) GetUserByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	query := `
		SELECT uid, username, password, first_name, second_name, age, salary, work_sphere
		FROM users
		WHERE uid = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.SecondName,
		&user.Age,
		&user.Salary,
		&user.WorkSphereID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
