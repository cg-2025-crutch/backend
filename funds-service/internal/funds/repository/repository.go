package repository

import (
	"context"

	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FundsRepositorer interface {
	// Transaction methods
	CreateTransaction(ctx context.Context, input models.CreateTransactionInput) (*models.Transaction, error)
	GetTransactionById(ctx context.Context, id int64) (*models.Transaction, error)
	GetUserTransactions(ctx context.Context, userUID string, limit, offset int32) ([]*models.Transaction, int64, error)
	GetUserTransactionsByPeriod(ctx context.Context, userUID string, days int32, limit, offset int32) ([]*models.Transaction, int64, error)
	UpdateTransaction(ctx context.Context, input models.UpdateTransactionInput) (*models.Transaction, error)
	DeleteTransaction(ctx context.Context, id int64, userUID string) error

	// Category methods
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoriesByType(ctx context.Context, categoryType string) ([]*models.Category, error)
	GetCategoryById(ctx context.Context, id int32) (*models.Category, error)

	// Balance methods
	GetUserBalance(ctx context.Context, userUID string) (*models.UserBalance, error)
}

type FundsRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) FundsRepositorer {
	return &FundsRepository{db: db}
}
