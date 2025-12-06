package repository

import (
	"context"
	"errors"

	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameExists = errors.New("username already exists")
)

type UserRepository interface {
	CreateUser(ctx context.Context, dto models.CreateUserDTO) (*models.User, error)
	GetUserByID(ctx context.Context, uid uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, dto models.UpdateUserDTO) (*models.User, error)
	DeleteUser(ctx context.Context, uid uuid.UUID) error
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) UserRepository {
	return &postgresRepository{
		db: db,
	}
}
