package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/cg-2025-crutch/backend/api-gateway/docs"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/clients"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/config"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/handlers/analytics"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/handlers/funds"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/handlers/notifications"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/handlers/user"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"go.uber.org/zap"
)

var defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func Run(mainCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(mainCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zapLogger := log.New(defaultLevel, os.Stdout)
	ctx = log.WithLogger(ctx, zapLogger)

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to import config: %w", err)
	}

	// Initialize gRPC clients
	grpcClients, err := clients.NewGRPCClients(conf)
	if err != nil {
		return fmt.Errorf("failed to initialize grpc clients: %w", err)
	}
	defer grpcClients.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "API Gateway",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Initialize handlers
	userHandler := user.NewUserHandler(grpcClients)
	fundsHandler := funds.NewFundsHandler(grpcClients)
	notificationsHandler := notifications.NewNotificationsHandler(grpcClients)
	analyticsHandler := analytics.NewAnalyticsHandler(grpcClients)

	// Register public routes (user registration and login)
	userHandler.RegisterPublicRoutes(api)

	// Protected routes - apply auth middleware
	authMiddleware := middleware.AuthMiddleware(grpcClients)

	// Create protected group for secured routes
	protected := api.Group("")
	protected.Use(authMiddleware)

	// Register protected routes
	userHandler.RegisterSecuredRoutes(protected)
	fundsHandler.RegisterRoutes(protected)
	notificationsHandler.RegisterRoutes(protected)
	analyticsHandler.RegisterRoutes(protected)

	// Start server
	addr := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	zapLogger.Info("Starting API Gateway", zap.String("address", addr))

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		zapLogger.Info("Shutting down API Gateway...")
		if err := app.Shutdown(); err != nil {
			zapLogger.Error("Error during shutdown", zap.Error(err))
		}
	}()

	if err := app.Listen(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
