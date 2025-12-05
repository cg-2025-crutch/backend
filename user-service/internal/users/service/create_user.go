package service

import (
	"context"
	"fmt"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) CreateUser(ctx context.Context, dto models.CreateUserDTO) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	dto.Password = string(hashedPassword)

	user, err := s.repo.CreateUser(ctx, dto)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
