package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	consumers "github.com/cg-2025-crutch/backend/notification-service/internal/adapters/consumer"
	"github.com/cg-2025-crutch/backend/notification-service/internal/config"
	"github.com/cg-2025-crutch/backend/notification-service/internal/grpc/server"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/kafka"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/redis"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/repository"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/service"
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

	redisClient, err := redis.NewRedisClient(ctx, conf.Redis)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer (*redisClient).Close()

	repo := repository.NewRedisRepo(*redisClient)

	notificationService := service.NewNotificationService(conf.Notif, repo)

	kafkaConsumer, err := kafka.InitKafkaConsumer(ctx, "notification-consumer", conf.Kafka)
	if err != nil {
		logger.Fatal("failed to initialize Kafka consumer", zap.Error(err))
	}
	defer kafkaConsumer.Close()

	consumerHandler := consumers.NewNotificationConsumer(notificationService)

	consumers.StartConsuming(ctx, kafkaConsumer, conf.Kafka.ConsTopic, consumerHandler)
	logger.Infof("Kafka consumer started, listening to topics: %v", conf.Kafka.ConsTopic)

	grpcServer := server.NewGrpcServer(conf.Server)

	grpcServer.RegisterGRPC(*notificationService)

	if err = grpcServer.Run(ctx); err != nil {
		return fmt.Errorf("appRun: %w", err)
	}
	return nil
}
