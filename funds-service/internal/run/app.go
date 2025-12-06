package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cg-2025-crutch/backend/funds-service/internal/config"
	"github.com/cg-2025-crutch/backend/funds-service/internal/funds/repository"
	"github.com/cg-2025-crutch/backend/funds-service/internal/funds/service"
	grpcserver "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/server"
	"github.com/cg-2025-crutch/backend/funds-service/internal/infrastructure/kafka"
	"github.com/cg-2025-crutch/backend/funds-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/funds-service/internal/infrastructure/postgres"
	"go.uber.org/zap"
)

var defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func Run(mainCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(mainCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	l := log.New(defaultLevel, os.Stdout)
	ctx = log.WithLogger(ctx, l)

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to import config: %w", err)
	}

	db, err := postgres.NewPostgresDB(ctx, conf.Postgres)
	if err != nil {
		l.Fatal("Failed to connect to database", zap.Error(err))
		return err
	}

	migr, err := postgres.NewMigrator(conf.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	if err := migr.Up(ctx, conf.Postgres.MigrateTimeout); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	producer, err := kafka.InitKafkaProducer(ctx, "fs-service", conf.Kafka)
	if err != nil {
		l.Errorf("appRun: failed to init kafka producer: %v", err)
		return fmt.Errorf("appRun: %w", err)
	}

	repo := repository.NewRepository(db.Pool)
	svc := service.NewService(repo, producer)

	grpcServer := grpcserver.NewGrpcServer(conf.Server)

	grpcServer.RegisterGRPC(svc)

	if err = grpcServer.Run(ctx); err != nil {
		return fmt.Errorf("appRun: %w", err)
	}
	return nil

	return nil
}
