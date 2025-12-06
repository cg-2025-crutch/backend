package funds

import (
	"context"
	"time"

	funds_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/funds_service"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *FundsHandler) GetUserBalance(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetUserBalance(ctx, &funds_pb.GetUserBalanceRequest{
		UserUid: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"balance": resp.Balance,
	})
}
