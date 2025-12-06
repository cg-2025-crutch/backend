package main

import (
	"context"
	"log"

	_ "github.com/cg-2025-crutch/backend/api-gateway/docs"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/run"
)

// @title API Gateway
// @version 1.0
// @description API Gateway для управления пользователями, финансами и уведомлениями
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	ctx := context.Background()

	if err := run.Run(ctx); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}
}
