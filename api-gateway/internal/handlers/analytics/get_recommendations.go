package analytics

import (
	"context"
	"time"

	analytics_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/analytics_service"
	"github.com/cg-2025-crutch/backend/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// GetRecommendations godoc
// @Summary Получить рекомендации пользователя
// @Description Получает финансовые рекомендации для пользователя на основе его транзакций и категорий расходов
// @Tags analytics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Рекомендации пользователя"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /analytics/recommendations [get]
func (h *AnalyticsHandler) GetRecommendations(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	resp, err := h.clients.AnalyticsService.GetRecommendations(ctx, &analytics_pb.GetRecommendationsReq{
		UserUid: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_uid":         resp.UserUid,
		"salary":           resp.Salary,
		"salary_bracket":   resp.SalaryBracket,
		"total_categories": resp.TotalCategories,
		"excellent_count":  resp.ExcellentCount,
		"normal_count":     resp.NormalCount,
		"warning_count":    resp.WarningCount,
		"critical_count":   resp.CriticalCount,
		"overall_status":   resp.OverallStatus,
		"overall_message":  resp.OverallMessage,
		"recommendations":  resp.Recommendations,
		"calculated_at":    resp.CalculatedAt,
	})
}
