package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/jackc/pgx/v5"
)

func (r *postgresRepository) UpdateUser(ctx context.Context, dto models.UpdateUserDTO) (*models.User, error) {
	query := `
		UPDATE users
		SET username = $2, first_name = $3, second_name = $4, age = $5, salary = $6, work_sphere = $7
		WHERE uid = $1
		RETURNING uid, username, password, first_name, second_name, age, salary, work_sphere
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query,
		dto.UID,
		dto.Username,
		dto.FirstName,
		dto.SecondName,
		dto.Age,
		dto.Salary,
		dto.WorkSphereID,
	).Scan(
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
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
