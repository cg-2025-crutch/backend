package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cg-2025-crutch/backend/funds-service/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db  *sql.DB
	dir string
}

func NewMigrator(dbConf config.PostgresConfig) (*Migrator, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name, dbConf.SSLMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB for migrator: %w", err)
	}

	return &Migrator{db: db, dir: dbConf.MigrationsDir}, nil
}

func (m *Migrator) Up(parentCtx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	err := goose.UpContext(ctx, m.db, m.dir)
	if err != nil {
		return fmt.Errorf("migrator: failed to up: %w", err)
	}

	return nil
}

func (m *Migrator) Down() error {
	err := goose.Down(m.db, m.dir)
	if err != nil {
		return fmt.Errorf("migrator: failed to rollback: %w", err)
	}

	return nil
}
