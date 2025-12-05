package run

import (
	"context"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/cg-2025-crutch/backend/user-service/internal/config"
	"github.com/cg-2025-crutch/backend/user-service/internal/grpc/server"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/jwt"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/postgres"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/repository"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/service"
	"go.uber.org/zap"
)

var defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func Run(mainCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(mainCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := log.New(defaultLevel, os.Stdout)
	ctx = log.WithLogger(ctx, logger)

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to import config: %w", err)
	}

	pg, err := postgres.NewPostgresDB(ctx, conf.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	migr, err := postgres.NewMigrator(conf.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	if err := migr.Up(ctx, conf.Postgres.MigrateTimeout); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	jwtManager := jwt.NewJWTManager(
		conf.JWT.SecretKey,
		conf.JWT.AccessTokenTTL,
		conf.JWT.RefreshTokenTTL,
	)
	userRepo := repository.NewPostgresRepository(pg.Pool)

	userService := service.NewUserService(userRepo, jwtManager)

	grpcServer := server.NewGrpcServer(conf.Server)

	grpcServer.RegisterGRPC(userService)

	if err = grpcServer.Run(ctx); err != nil {
		return fmt.Errorf("appRun: %w", err)
	}
	return nil
}
