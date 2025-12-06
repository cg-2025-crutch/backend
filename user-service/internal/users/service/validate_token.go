package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/jwt"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/repository"
)

func (s *userService) ValidateToken(ctx context.Context, token string) (*models.TokenClaims, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, jwt.ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	return claims, nil
}
