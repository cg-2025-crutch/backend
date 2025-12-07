package analytics

import (
	"github.com/cg-2025-crutch/backend/api-gateway/internal/clients"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	clients *clients.GRPCClients
}

func NewAnalyticsHandler(clients *clients.GRPCClients) *AnalyticsHandler {
	return &AnalyticsHandler{
		clients: clients,
	}
}

func (h *AnalyticsHandler) RegisterRoutes(router fiber.Router) {
	analytics := router.Group("/analytics")

	// All routes require authentication
	analytics.Get("/recommendations", h.GetRecommendations)
}
