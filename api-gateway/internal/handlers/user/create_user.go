package user

import (
	"context"
	"time"

	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Username     string  `json:"username" validate:"required" example:"john_doe"`
	Password     string  `json:"password" validate:"required,min=6" example:"password123"`
	FirstName    string  `json:"first_name" validate:"required" example:"John"`
	SecondName   string  `json:"second_name" validate:"required" example:"Doe"`
	Age          int32   `json:"age" validate:"required,min=18" example:"25"`
	Salary       float64 `json:"salary" validate:"required,min=0" example:"50000"`
	WorkSphereId int64   `json:"work_sphere_id" validate:"required" example:"1"`
}

// CreateUser godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Данные нового пользователя"
// @Success 201 {object} map[string]interface{} "Пользователь успешно создан"
// @Failure 400 {object} map[string]interface{} "Неверный формат запроса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/register [post]
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
