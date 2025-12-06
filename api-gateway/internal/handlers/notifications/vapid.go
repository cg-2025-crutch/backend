package notifications

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

// GetVapidKey godoc
// @Summary Получить VAPID ключ
// @Description Получает публичный VAPID ключ для push-уведомлений
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "VAPID ключ"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /notifications/vapid-key [get]
func (h *NotificationsHandler) GetVapidKey(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.NotifService.GetVapidKey(ctx, &emptypb.Empty{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"vapid_key": resp.VapidKey,
	})
}
