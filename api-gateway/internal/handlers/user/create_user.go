package user

import (
	"context"
	"time"

	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Username     string  `json:"username" validate:"required"`
	Password     string  `json:"password" validate:"required,min=6"`
	FirstName    string  `json:"first_name" validate:"required"`
	SecondName   string  `json:"second_name" validate:"required"`
	Age          int32   `json:"age" validate:"required,min=18"`
	Salary       float64 `json:"salary" validate:"required,min=0"`
	WorkSphereId int64   `json:"work_sphere_id" validate:"required"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.UserService.CreateUser(ctx, &user_pb.CreateUserRequest{
		Username:     req.Username,
		Password:     req.Password,
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": resp.User,
	})
}
