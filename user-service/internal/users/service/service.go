package service

import (
	"context"
	"errors"

	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/jwt"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = repository.ErrUserNotFound
	ErrUsernameExists     = repository.ErrUsernameExists
)

// UserService defines the business logic interface
type UserService interface {
	CreateUser(ctx context.Context, dto models.CreateUserDTO) (*models.User, error)
	UpdateUser(ctx context.Context, dto models.UpdateUserDTO) (*models.User, error)
	Login(ctx context.Context, username, password string) (*models.TokenPair, *models.User, error)
	ValidateToken(ctx context.Context, token string) (*models.TokenClaims, error)
	GetUserByID(ctx context.Context, uid uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// userService implements UserService
type userService struct {
	repo       repository.UserRepository
	jwtManager *jwt.JWTManager
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, jwtManager *jwt.JWTManager) UserService {
	return &userService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}
