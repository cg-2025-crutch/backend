package user

import (
	"context"
	"time"

	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.UserService.Login(ctx, &user_pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_at":    resp.ExpiresAt,
		"user":          resp.User,
	})
}
