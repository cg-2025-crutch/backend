package notifications

import (
	"context"
	"time"

	notif_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/notification_service"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type SubscribeRequest struct {
	Endpoint string `json:"endpoint" validate:"required" example:"https://fcm.googleapis.com/fcm/send/..."`
	P256dh   string `json:"p256dh" validate:"required" example:"BPxY..."`
	Auth     string `json:"auth" validate:"required" example:"xyz123..."`
}

// Subscribe godoc
// @Summary Подписаться на уведомления
// @Description Подписывает пользователя на push-уведомления
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body SubscribeRequest true "Данные подписки"
// @Success 200 {object} map[string]interface{} "Успешная подписка"
// @Failure 400 {object} map[string]interface{} "Неверный формат запроса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /notifications/subscribe [post]
func (h *NotificationsHandler) Subscribe(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	var req SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.NotifService.Subscribe(ctx, &notif_pb.SubscribeReq{
		UserId:   userID,
		Endpoint: req.Endpoint,
		P256Dh:   req.P256dh,
		Auth:     req.Auth,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": resp.Message,
	})
}
