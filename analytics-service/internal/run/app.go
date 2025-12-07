package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cg-2025-crutch/backend/analytics-service/internal/adapters/consumers"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/handler"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/repository"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/service"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/clients"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/config"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/server"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/kafka"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/redis"
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

	// Подключаемся к Redis
	redisClient, err := redis.NewRedisClient(ctx, conf.Redis)
	if err != nil {
		logger.Fatal("failed to connect to Redis", zap.Error(err))
		return err
	}
	defer redisClient.Close()
	logger.Info("Redis client initialized")

	// Создаем repository
	repo := repository.NewRedisRepository(redisClient)

	// Создаем gRPC клиенты для других сервисов
	grpcClients, err := clients.NewClients(conf.UserServiceAddr, conf.FundsServiceAddr)
	if err != nil {
		logger.Fatal("failed to create gRPC clients", zap.Error(err))
		return err
	}
	defer grpcClients.Close()
	logger.Info("gRPC clients initialized")

	// Создаем сервис аналитики
	analyticsService := service.NewAnalyticsService(grpcClients, repo, conf.Redis)

	// Инициализируем Kafka consumer
	kafkaConsumer, err := kafka.InitKafkaConsumer(ctx, "analytics-consumer", conf.Kafka)
	if err != nil {
		logger.Fatal("failed to initialize Kafka consumer", zap.Error(err))
		return err
	}
	defer kafkaConsumer.Close()

	// Создаем обработчик для Kafka
	consumerHandler := consumers.NewAnalyticsConsumer(analyticsService)

	// Запускаем consumer
	consumers.StartConsuming(ctx, kafkaConsumer, conf.Kafka.ConsTopic, consumerHandler)
	logger.Infof("Kafka consumer started, listening to topics: %v", conf.Kafka.ConsTopic)

	// Создаем gRPC handler
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Создаем и запускаем gRPC server
	grpcServer := server.NewGrpcServer(conf.Server)
	grpcServer.RegisterGRPC(analyticsHandler)

	logger.Infof("Starting gRPC server on %s:%s", conf.Server.GRPCHost, conf.Server.GRPCPort)

	if err = grpcServer.Run(ctx); err != nil {
		return fmt.Errorf("appRun: %w", err)
	}

	return nil
}
