package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/repository"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Login(ctx context.Context, username, password string) (*models.TokenPair, *models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(user.UID, user.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	user.Password = ""

	return tokenPair, user, nil
}
