package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/repository"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/clients"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/config"
	fundspb "github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/gen/funds_service"
	userpb "github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/gen/user_service"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/models"
)

// AnalyticsService предоставляет методы для работы с аналитикой
type AnalyticsService struct {
	clients *clients.Clients
	repo    repository.Repository
	ttl     time.Duration
}

// NewAnalyticsService создает новый экземпляр сервиса аналитики
func NewAnalyticsService(clients *clients.Clients, repo repository.Repository, cfg config.RedisConfig) *AnalyticsService {
	return &AnalyticsService{
		clients: clients,
		repo:    repo,
		ttl:     cfg.TTL,
	}
}

// GetRecommendations возвращает рекомендации из Redis для пользователя
func (s *AnalyticsService) GetRecommendations(ctx context.Context, userUID string) (*models.AnalyticsResult, error) {
	l := log.FromContext(ctx)

	l.Infof("Getting recommendations for user %s from Redis", userUID)

	// Получаем из Redis
	result, err := s.repo.GetRecommendations(ctx, userUID)
	if err != nil {
		l.Errorf("Failed to get recommendations from Redis: %v", err)
		return nil, fmt.Errorf("failed to get recommendations from Redis: %w", err)
	}

	if result == nil {
		l.Infof("No recommendations found in Redis for user %s", userUID)
		return nil, fmt.Errorf("no recommendations found for user %s, please wait for calculation", userUID)
	}

	l.Infof("Successfully retrieved recommendations for user %s from Redis", userUID)
	return result, nil
}

// CalculateAndSaveRecommendations рассчитывает и сохраняет рекомендации в Redis
func (s *AnalyticsService) CalculateAndSaveRecommendations(ctx context.Context, userUID string) error {
	l := log.FromContext(ctx)

	l.Infof("Calculating recommendations for user %s", userUID)

	// Получаем данные пользователя (зарплату)
	userResp, err := s.clients.UserClient.GetUserById(ctx, &userpb.GetUserByIdRequest{
		Id: userUID,
	})
	if err != nil {
		l.Errorf("Failed to get user data: %v", err)
		return fmt.Errorf("failed to get user data: %w", err)
	}

	salary := userResp.User.Salary
	l.Infof("User %s salary: %.2f", userUID, salary)

	// Получаем транзакции за последний месяц (30 дней)
	transactionsResp, err := s.clients.FundsClient.GetUserTransactionsByPeriod(ctx, &fundspb.GetUserTransactionsByPeriodRequest{
		UserUid: userUID,
		Days:    30,
		Limit:   1000,
		Offset:  0,
	})
	if err != nil {
		l.Errorf("Failed to get transactions: %v", err)
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	l.Infof("Found %d transactions for user %s", len(transactionsResp.Transactions), userUID)

	// Получаем все категории
	categoriesResp, err := s.clients.FundsClient.GetAllCategories(ctx, &fundspb.GetAllCategoriesRequest{})
	if err != nil {
		l.Errorf("Failed to get categories: %v", err)
		return fmt.Errorf("failed to get categories: %w", err)
	}

	// Создаем карту категорий: id -> имя
	categoryMap := make(map[int32]string)
	for _, cat := range categoriesResp.Categories {
		categoryMap[cat.Id] = cat.Name
	}

	// Группируем траты по категориям
	categorySpending := s.calculateCategorySpending(transactionsResp.Transactions, categoryMap, salary)

	// Генерируем рекомендации
	result := s.generateRecommendations(userUID, salary, categorySpending)

	// Сохраняем в Redis
	err = s.repo.SaveRecommendations(ctx, userUID, result, s.ttl)
	if err != nil {
		l.Errorf("Failed to save recommendations to Redis: %v", err)
		return fmt.Errorf("failed to save recommendations to Redis: %w", err)
	}

	l.Infof("Successfully calculated and saved recommendations for user %s", userUID)
	return nil
}

// calculateCategorySpending группирует траты по категориям и вычисляет проценты от зарплаты
func (s *AnalyticsService) calculateCategorySpending(transactions []*fundspb.Transaction, categoryMap map[int32]string, salary float64) []models.CategorySpending {
	spendingMap := make(map[string]float64)

	// Суммируем расходы по категориям
	for _, tx := range transactions {
		if tx.Type == "expense" {
			categoryName := categoryMap[tx.CategoryId]
			if categoryName != "" {
				spendingMap[categoryName] += tx.Amount
			}
		}
	}

	// Конвертируем в массив с процентами
	var result []models.CategorySpending
	for catName, amount := range spendingMap {
		percentage := 0.0
		if salary > 0 {
			percentage = (amount / salary) * 100
		}

		result = append(result, models.CategorySpending{
			CategoryName: catName,
			Amount:       amount,
			Percentage:   percentage,
		})
	}

	return result
}

// generateRecommendations генерирует список рекомендаций по категориям
func (s *AnalyticsService) generateRecommendations(userUID string, salary float64, spending []models.CategorySpending) *models.AnalyticsResult {
	bracket := models.GetBracketForSalary(salary)
	if bracket == nil {
		return &models.AnalyticsResult{
			UserUID:         userUID,
			Salary:          salary,
			SalaryBracket:   "unknown",
			TotalCategories: 0,
			OverallStatus:   models.OverallStatusCritical,
			OverallMessage:  "Не удалось определить диапазон зарплаты для формирования рекомендаций.",
			Recommendations: []models.CategoryRecommendation{},
			CalculatedAt:    time.Now().Unix(),
		}
	}

	// Создаем карту рекомендуемых диапазонов
	rangeMap := make(map[string]models.CategoryRange)
	for _, cr := range bracket.Categories {
		rangeMap[cr.Name] = cr
	}

	var recommendations []models.CategoryRecommendation
	var excellentCount, normalCount, warningCount, criticalCount int

	// Анализируем каждую категорию трат
	for _, cs := range spending {
		rec, exists := rangeMap[cs.CategoryName]
		if !exists {
			continue
		}

		// Определяем статус и сообщение
		var status string
		var message string
		var deviation float64

		// Очень плохо - превышение на 20% от верхней границы
		criticalThreshold := rec.MaxPerc * 1.2

		if cs.Percentage > criticalThreshold {
			status = models.StatusCritical
			deviation = cs.Percentage - rec.MaxPerc
			message = fmt.Sprintf("Критическое превышение на %.1f%%. Необходимо срочно сократить расходы!", deviation)
			criticalCount++
		} else if cs.Percentage > rec.MaxPerc {
			status = models.StatusWarning
			deviation = cs.Percentage - rec.MaxPerc
			message = fmt.Sprintf("Превышение на %.1f%%. Рекомендуем снизить расходы.", deviation)
			warningCount++
		} else if cs.Percentage < rec.MinPerc && rec.MinPerc > 0 {
			status = models.StatusExcellent
			deviation = rec.MinPerc - cs.Percentage
			message = fmt.Sprintf("Вы очень экономны! Экономия %.1f%% от минимальной нормы.", deviation)
			excellentCount++
		} else {
			status = models.StatusNormal
			deviation = 0
			message = "Расходы в пределах нормы."
			normalCount++
		}

		recommendations = append(recommendations, models.CategoryRecommendation{
			CategoryName:     cs.CategoryName,
			ActualAmount:     cs.Amount,
			ActualPercentage: cs.Percentage,
			RecommendedMin:   rec.MinPerc,
			RecommendedMax:   rec.MaxPerc,
			Status:           status,
			Message:          message,
			Deviation:        deviation,
		})
	}

	// Определяем общий статус
	overallStatus := s.determineOverallStatus(excellentCount, normalCount, warningCount, criticalCount)
	overallMessage := s.generateOverallMessage(overallStatus, excellentCount, normalCount, warningCount, criticalCount)

	// Формируем диапазон зарплаты
	salaryBracket := fmt.Sprintf("%.0f-%.0f", bracket.MinSalary, bracket.MaxSalary)
	if bracket.MaxSalary >= 999999999 {
		salaryBracket = fmt.Sprintf("%.0f+", bracket.MinSalary)
	}

	return &models.AnalyticsResult{
		UserUID:         userUID,
		Salary:          salary,
		SalaryBracket:   salaryBracket,
		TotalCategories: len(recommendations),
		ExcellentCount:  excellentCount,
		NormalCount:     normalCount,
		WarningCount:    warningCount,
		CriticalCount:   criticalCount,
		OverallStatus:   overallStatus,
		OverallMessage:  overallMessage,
		Recommendations: recommendations,
		CalculatedAt:    time.Now().Unix(),
	}
}

// determineOverallStatus определяет общий статус на основе статистики категорий
func (s *AnalyticsService) determineOverallStatus(excellent, normal, warning, critical int) string {
	total := excellent + normal + warning + critical
	if total == 0 {
		return models.OverallStatusGood
	}

	if critical > 0 && critical >= total/3 {
		return models.OverallStatusCritical
	}
	if warning+critical > total/2 {
		return models.OverallStatusAttentionRequired
	}
	if excellent > total/2 {
		return models.OverallStatusExcellent
	}
	return models.OverallStatusGood
}

// generateOverallMessage генерирует общее сообщение
func (s *AnalyticsService) generateOverallMessage(status string, excellent, normal, warning, critical int) string {
	switch status {
	case models.OverallStatusExcellent:
		return "Отличное управление финансами! Вы контролируете большинство расходов."
	case models.OverallStatusGood:
		return "Ваши финансы в порядке. Продолжайте в том же духе."
	case models.OverallStatusAttentionRequired:
		return "Требуется внимание к расходам. Рекомендуем оптимизировать бюджет."
	case models.OverallStatusCritical:
		return "Критическая ситуация с финансами. Необходимо срочно пересмотреть расходы."
	default:
		return "Анализ расходов завершен."
	}
}

// ProcessAnalyticsEvent обрабатывает событие аналитики из Kafka
func (s *AnalyticsService) ProcessAnalyticsEvent(ctx context.Context, userUID string) error {
	l := log.FromContext(ctx)

	l.Infof("Processing analytics event for user: %s", userUID)

	// Рассчитываем и сохраняем рекомендации в Redis
	err := s.CalculateAndSaveRecommendations(ctx, userUID)
	if err != nil {
		l.Errorf("Failed to calculate and save recommendations: %v", err)
		return fmt.Errorf("failed to calculate and save recommendations: %w", err)
	}

	l.Infof("Successfully processed analytics event and saved recommendations for user %s", userUID)
	return nil
}
