package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *postgresRepository) DeleteUser(ctx context.Context, uid uuid.UUID) error {
	query := `DELETE FROM users WHERE uid = $1`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
