package funds

import (
	"context"
	"strconv"
	"time"

	funds_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/funds_service"
	"github.com/gofiber/fiber/v2"
)

// GetAllCategories godoc
// @Summary Получить все категории
// @Description Получает список всех категорий транзакций
// @Tags funds
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Список категорий"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/categories [get]
func (h *FundsHandler) GetAllCategories(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetAllCategories(ctx, &funds_pb.GetAllCategoriesRequest{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": resp.Categories,
	})
}

// GetCategoriesByType godoc
// @Summary Получить категории по типу
// @Description Получает список категорий по типу (income или expense)
// @Tags funds
// @Accept json
// @Produce json
// @Param type path string true "Тип категории (income или expense)"
// @Success 200 {object} map[string]interface{} "Список категорий"
// @Failure 400 {object} map[string]interface{} "Неверный тип категории"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/categories/type/{type} [get]
func (h *FundsHandler) GetCategoriesByType(c *fiber.Ctx) error {
	categoryType := c.Params("type")
	if categoryType != "income" && categoryType != "expense" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type must be 'income' or 'expense'",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetCategoriesByType(ctx, &funds_pb.GetCategoriesByTypeRequest{
		Type: categoryType,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": resp.Categories,
	})
}

// GetCategoryById godoc
// @Summary Получить категорию по ID
// @Description Получает информацию о категории по ее ID
// @Tags funds
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} map[string]interface{} "Информация о категории"
// @Failure 400 {object} map[string]interface{} "Неверный ID категории"
// @Failure 404 {object} map[string]interface{} "Категория не найдена"
// @Security BearerAuth
// @Router /funds/categories/{id} [get]
func (h *FundsHandler) GetCategoryById(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid category id",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetCategoryById(ctx, &funds_pb.GetCategoryByIdRequest{
		Id: int32(categoryID),
	})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "category not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"category": resp.Category,
	})
}
