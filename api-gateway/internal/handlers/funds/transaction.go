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
	CategoryId      int32   `json:"category_id" validate:"required" example:"1"`
	Type            string  `json:"type" validate:"required,oneof=income expense" example:"income"`
	Amount          float64 `json:"amount" validate:"required,gt=0" example:"1000.50"`
	Title           string  `json:"title" validate:"required" example:"Зарплата"`
	Description     string  `json:"description" example:"Месячная зарплата"`
	TransactionDate string  `json:"transaction_date" validate:"required" example:"2025-12-06"` // YYYY-MM-DD
}

// CreateTransaction godoc
// @Summary Создать транзакцию
// @Description Создает новую финансовую транзакцию для пользователя
// @Tags funds
// @Accept json
// @Produce json
// @Param request body CreateTransactionRequest true "Данные транзакции"
// @Success 201 {object} map[string]interface{} "Транзакция успешно создана"
// @Failure 400 {object} map[string]interface{} "Неверный формат запроса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/transactions [post]
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

// GetTransactionById godoc
// @Summary Получить транзакцию по ID
// @Description Получает информацию о транзакции по ее ID
// @Tags funds
// @Accept json
// @Produce json
// @Param id path int true "ID транзакции"
// @Success 200 {object} map[string]interface{} "Информация о транзакции"
// @Failure 400 {object} map[string]interface{} "Неверный ID транзакции"
// @Failure 404 {object} map[string]interface{} "Транзакция не найдена"
// @Security BearerAuth
// @Router /funds/transactions/{id} [get]
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

// GetUserTransactions godoc
// @Summary Получить транзакции пользователя
// @Description Получает список всех транзакций пользователя с пагинацией
// @Tags funds
// @Accept json
// @Produce json
// @Param limit query int false "Лимит транзакций" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} map[string]interface{} "Список транзакций"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/transactions [get]
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

// GetUserTransactionsByPeriod godoc
// @Summary Получить транзакции за период
// @Description Получает транзакции пользователя за указанное количество дней
// @Tags funds
// @Accept json
// @Produce json
// @Param days query int false "Количество дней" default(30)
// @Param limit query int false "Лимит транзакций" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} map[string]interface{} "Список транзакций за период"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/transactions/period [get]
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
	CategoryId      int32   `json:"category_id" example:"2"`
	Type            string  `json:"type" example:"expense"`
	Amount          float64 `json:"amount" example:"500.00"`
	Title           string  `json:"title" example:"Продукты"`
	Description     string  `json:"description" example:"Покупка продуктов"`
	TransactionDate string  `json:"transaction_date" example:"2025-12-06"`
}

// UpdateTransaction godoc
// @Summary Обновить транзакцию
// @Description Обновляет информацию о транзакции
// @Tags funds
// @Accept json
// @Produce json
// @Param id path int true "ID транзакции"
// @Param request body UpdateTransactionRequest true "Обновленные данные транзакции"
// @Success 200 {object} map[string]interface{} "Транзакция успешно обновлена"
// @Failure 400 {object} map[string]interface{} "Неверный формат запроса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/transactions/{id} [put]
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

// DeleteTransaction godoc
// @Summary Удалить транзакцию
// @Description Удаляет транзакцию по ее ID
// @Tags funds
// @Accept json
// @Produce json
// @Param id path int true "ID транзакции"
// @Success 200 {object} map[string]interface{} "Транзакция успешно удалена"
// @Failure 400 {object} map[string]interface{} "Неверный ID транзакции"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /funds/transactions/{id} [delete]
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
