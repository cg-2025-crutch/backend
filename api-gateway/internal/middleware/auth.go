package middleware

import (
	"strings"

	"github.com/cg-2025-crutch/backend/api-gateway/internal/clients"
	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey = "user_id"
)

func AuthMiddleware(grpcClients *clients.GRPCClients) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token via user service
		userID, err := grpcClients.ValidateToken(c.UserContext(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Store user ID in context
		c.Locals(UserIDKey, userID)

		return c.Next()
	}
}
