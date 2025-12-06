package funds

import (
	"context"
	"strconv"
	"time"

	funds_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/funds_service"
	"github.com/gofiber/fiber/v2"
)

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
