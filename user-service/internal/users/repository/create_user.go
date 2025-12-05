package repository

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/google/uuid"
)

func (r *postgresRepository) CreateUser(ctx context.Context, dto models.CreateUserDTO) (*models.User, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", dto.Username).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, ErrUsernameExists
	}

	uid := uuid.New()

	query := `
		INSERT INTO users (uid, username, password, first_name, second_name, age, salary, work_sphere)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING uid, username, password, first_name, second_name, age, salary, work_sphere
	`

	user := &models.User{}
	err = r.db.QueryRow(ctx, query,
		uid,
		dto.Username,
		dto.Password,
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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
