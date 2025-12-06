package user

import (
	"context"
	"time"

	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type UpdateUserRequest struct {
	Username     string  `json:"username" example:"john_doe"`
	FirstName    string  `json:"first_name" example:"John"`
	SecondName   string  `json:"second_name" example:"Doe"`
	Age          int32   `json:"age" example:"26"`
	Salary       float64 `json:"salary" example:"55000"`
	WorkSphereId int64   `json:"work_sphere_id" example:"2"`
}

// UpdateUser godoc
// @Summary Обновить профиль пользователя
// @Description Обновляет информацию о пользователе (только свой профиль)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param request body UpdateUserRequest true "Обновленные данные пользователя"
// @Success 200 {object} map[string]interface{} "Пользователь успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Неверный формат запроса"
// @Failure 403 {object} map[string]interface{} "Можно обновлять только свой профиль"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /users/{id} [put]
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
