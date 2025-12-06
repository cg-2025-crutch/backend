package funds

import (
	"github.com/cg-2025-crutch/backend/api-gateway/internal/clients"
	"github.com/gofiber/fiber/v2"
)

type FundsHandler struct {
	clients *clients.GRPCClients
}

func NewFundsHandler(clients *clients.GRPCClients) *FundsHandler {
	return &FundsHandler{
		clients: clients,
	}
}

func (h *FundsHandler) RegisterRoutes(router fiber.Router) {
	funds := router.Group("/funds")

	// All routes require authentication
	// Transactions
	funds.Post("/transactions", h.CreateTransaction)
	funds.Get("/transactions/:id", h.GetTransactionById)
	funds.Get("/transactions", h.GetUserTransactions)
	funds.Get("/transactions/period", h.GetUserTransactionsByPeriod)
	funds.Put("/transactions/:id", h.UpdateTransaction)
	funds.Delete("/transactions/:id", h.DeleteTransaction)

	// Categories
	funds.Get("/categories", h.GetAllCategories)
	funds.Get("/categories/type/:type", h.GetCategoriesByType)
	funds.Get("/categories/:id", h.GetCategoryById)

	// Balance
	funds.Get("/balance", h.GetUserBalance)
}
