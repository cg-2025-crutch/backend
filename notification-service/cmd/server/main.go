package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/cg-2025-crutch/backend/notification-service/internal/adapters/consumer"
	"github.com/cg-2025-crutch/backend/notification-service/internal/config"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/kafka"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/redis"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/repository"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/service"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize logger
	logger := log.NewStdOut(zap.NewAtomicLevelAt(zap.InfoLevel))
	if cfg.Logger.Development {
		logger = log.NewStdOut(zap.NewAtomicLevelAt(zap.DebugLevel))
	}
	if cfg.Logger.Caller {
		logger = logger.WithOptions(zap.AddCaller())
	}
	log.SetLogger(logger)
	ctx = log.WithLogger(ctx, logger)

	logger.Info("Starting notification service...")

	// Initialize Redis client
	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer (*redisClient).Close()

	// Initialize repository
	repo := repository.NewRedisRepo(*redisClient)

	// Initialize notification service
	notificationService := service.NewNotificationService(cfg.Notif, repo)

	// Initialize HTTP controller
	controller := notifications.NewController(notificationService)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/vapid-key", controller.GetVapidKeyHandler)
	mux.HandleFunc("/api/subscribe", controller.SubscribeHandler)

	// Add CORS middleware
	handler := corsMiddleware(mux)

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.HTTPHost, cfg.Server.HTTPPort)
	httpServer := &http.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()

	fmt.Println("VAPID_PRIVATE_KEY =", privateKey)
	fmt.Println("VAPID_PUBLIC_KEY  =", publicKey)

	// Start HTTP server in a goroutine
	go func() {
		logger.Infof("HTTP server listening on %s", serverAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start HTTP server", zap.Error(err))
		}
	}()

	// Initialize Kafka consumer
	kafkaConsumer, err := kafka.InitKafkaConsumer(ctx, "notification-consumer", cfg.Kafka)
	if err != nil {
		logger.Fatal("failed to initialize Kafka consumer", zap.Error(err))
	}
	defer kafkaConsumer.Close()

	// Create notification consumer handler
	consumerHandler := consumers.NewNotificationConsumer(notificationService)

	// Start consuming Kafka messages
	consumers.StartConsuming(ctx, kafkaConsumer, cfg.Kafka.ConsTopic, consumerHandler)
	logger.Infof("Kafka consumer started, listening to topics: %v", cfg.Kafka.ConsTopic)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down service...")

	// Gracefully shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	logger.Info("Service stopped")
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
