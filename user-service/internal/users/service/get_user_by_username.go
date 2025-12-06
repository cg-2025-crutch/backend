package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
)

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
