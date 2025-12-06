package funds

import (
	"context"
	"strconv"
	"time"

	funds_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/funds_service"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type CreateTransactionRequest struct {
	CategoryId      int32   `json:"category_id" validate:"required"`
	Type            string  `json:"type" validate:"required,oneof=income expense"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Title           string  `json:"title" validate:"required"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date" validate:"required"` // YYYY-MM-DD
}

func (h *FundsHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.CreateTransaction(ctx, &funds_pb.CreateTransactionRequest{
		UserUid:         userID,
		CategoryId:      req.CategoryId,
		Type:            req.Type,
		Amount:          req.Amount,
		Title:           req.Title,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"transaction": resp.Transaction,
	})
}

func (h *FundsHandler) GetTransactionById(c *fiber.Ctx) error {
	transactionIDStr := c.Params("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transaction id",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetTransactionById(ctx, &funds_pb.GetTransactionByIdRequest{
		Id: transactionID,
	})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transaction not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transaction": resp.Transaction,
	})
}

func (h *FundsHandler) GetUserTransactions(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetUserTransactions(ctx, &funds_pb.GetUserTransactionsRequest{
		UserUid: userID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": resp.Transactions,
		"total":        resp.Total,
	})
}

func (h *FundsHandler) GetUserTransactionsByPeriod(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	days := c.QueryInt("days", 30)
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.GetUserTransactionsByPeriod(ctx, &funds_pb.GetUserTransactionsByPeriodRequest{
		UserUid: userID,
		Days:    int32(days),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": resp.Transactions,
		"total":        resp.Total,
	})
}

type UpdateTransactionRequest struct {
	CategoryId      int32   `json:"category_id"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
}

func (h *FundsHandler) UpdateTransaction(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	transactionIDStr := c.Params("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transaction id",
		})
	}

	var req UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.UpdateTransaction(ctx, &funds_pb.UpdateTransactionRequest{
		Id:              transactionID,
		CategoryId:      req.CategoryId,
		Type:            req.Type,
		Amount:          req.Amount,
		Title:           req.Title,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		UserUid:         userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transaction": resp.Transaction,
	})
}

func (h *FundsHandler) DeleteTransaction(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	transactionIDStr := c.Params("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transaction id",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.FundsService.DeleteTransaction(ctx, &funds_pb.DeleteTransactionRequest{
		Id:      transactionID,
		UserUid: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": resp.Success,
	})
}
