package user

import (
	"github.com/cg-2025-crutch/backend/api-gateway/internal/clients"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	clients *clients.GRPCClients
}

func NewUserHandler(clients *clients.GRPCClients) *UserHandler {
	return &UserHandler{
		clients: clients,
	}
}

// RegisterPublicRoutes registers public user routes (register, login)
func (h *UserHandler) RegisterPublicRoutes(router fiber.Router) {
	users := router.Group("/users")

	users.Post("/register", h.CreateUser)
	users.Post("/login", h.Login)
}

// RegisterSecuredRoutes registers secured user routes
func (h *UserHandler) RegisterSecuredRoutes(router fiber.Router) {
	users := router.Group("/users")

	users.Put("/:id", h.UpdateUser)
	users.Get("/:id", h.GetUserById)
	users.Get("/username/:username", h.GetUserByUsername)
}
