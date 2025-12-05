package service

import (
	"context"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
)

func (s *userService) UpdateUser(ctx context.Context, dto models.UpdateUserDTO) (*models.User, error) {
	_, err := s.repo.GetUserByID(ctx, dto.UID)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.UpdateUser(ctx, dto)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
