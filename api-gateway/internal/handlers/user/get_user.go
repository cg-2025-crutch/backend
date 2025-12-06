package user

import (
	"context"
	"time"

	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
	"github.com/gofiber/fiber/v2"
)

// GetUserById godoc
// @Summary Получить пользователя по ID
// @Description Получает информацию о пользователе по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} map[string]interface{} "Информация о пользователе"
// @Failure 400 {object} map[string]interface{} "ID пользователя обязателен"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.UserService.GetUserById(ctx, &user_pb.GetUserByIdRequest{
		Id: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": resp.User,
	})
}

// GetUserByUsername godoc
// @Summary Получить пользователя по имени
// @Description Получает информацию о пользователе по его username
// @Tags users
// @Accept json
// @Produce json
// @Param username path string true "Имя пользователя"
// @Success 200 {object} map[string]interface{} "Информация о пользователе"
// @Failure 400 {object} map[string]interface{} "Username обязателен"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Security BearerAuth
// @Router /users/username/{username} [get]
func (h *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.UserService.GetUserByUsername(ctx, &user_pb.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": resp.User,
	})
}
