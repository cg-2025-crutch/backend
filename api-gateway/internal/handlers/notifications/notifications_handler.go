package notifications

import (
	"github.com/cg-2025-crutch/backend/api-gateway/internal/clients"
	"github.com/gofiber/fiber/v2"
)

type NotificationsHandler struct {
	clients *clients.GRPCClients
}

func NewNotificationsHandler(clients *clients.GRPCClients) *NotificationsHandler {
	return &NotificationsHandler{
		clients: clients,
	}
}

func (h *NotificationsHandler) RegisterRoutes(router fiber.Router) {
	notifications := router.Group("/notifications")

	// All routes require authentication
	notifications.Get("/vapid-key", h.GetVapidKey)
	notifications.Post("/subscribe", h.Subscribe)
}
