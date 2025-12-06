package user

import (
	"context"
	"time"

	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type UpdateUserRequest struct {
	Username     string  `json:"username"`
	FirstName    string  `json:"first_name"`
	SecondName   string  `json:"second_name"`
	Age          int32   `json:"age"`
	Salary       float64 `json:"salary"`
	WorkSphereId int64   `json:"work_sphere_id"`
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is required",
		})
	}

	// Check if the authenticated user is updating their own profile
	authUserID := c.Locals(middleware.UserIDKey).(string)
	if authUserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you can only update your own profile",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.UserService.UpdateUser(ctx, &user_pb.UpdateUserRequest{
		Id:           userID,
		Username:     req.Username,
		FirstName:    req.FirstName,
		SecondName:   req.SecondName,
		Age:          req.Age,
		Salary:       req.Salary,
		WorkSphereId: req.WorkSphereId,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": resp.User,
	})
}
