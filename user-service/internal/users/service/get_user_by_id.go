package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/google/uuid"
)

func (s *userService) GetUserByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
