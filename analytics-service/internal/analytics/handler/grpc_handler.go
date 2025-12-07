package handler

import (
	"context"

	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/service"
	pb "github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/gen/analytics_service"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticsHandler обрабатывает gRPC запросы для аналитики
type AnalyticsHandler struct {
	pb.UnimplementedAnalyticsServiceServer
	service *service.AnalyticsService
}

// NewAnalyticsHandler создает новый обработчик
func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: svc,
	}
}

// GetRecommendations возвращает список рекомендаций для пользователя из Redis
func (h *AnalyticsHandler) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsReq) (*pb.GetRecommendationsResp, error) {
	l := log.FromContext(ctx)

	if req.UserUid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uid is required")
	}

	l.Infof("GetRecommendations called for user: %s", req.UserUid)

	// Получаем рекомендации из Redis
	result, err := h.service.GetRecommendations(ctx, req.UserUid)
	if err != nil {
		l.Errorf("Failed to get recommendations: %v", err)
		return nil, status.Error(codes.NotFound, "recommendations not found, please wait for calculation")
	}

	// Конвертируем в proto формат
	var pbRecommendations []*pb.CategoryRecommendation
	for _, r := range result.Recommendations {
		pbRecommendations = append(pbRecommendations, &pb.CategoryRecommendation{
			CategoryName:     r.CategoryName,
			ActualAmount:     r.ActualAmount,
			ActualPercentage: r.ActualPercentage,
			RecommendedMin:   r.RecommendedMin,
			RecommendedMax:   r.RecommendedMax,
			Status:           r.Status,
			Message:          r.Message,
			Deviation:        r.Deviation,
		})
	}

	return &pb.GetRecommendationsResp{
		UserUid:         result.UserUID,
		Salary:          result.Salary,
		SalaryBracket:   result.SalaryBracket,
		TotalCategories: int32(result.TotalCategories),
		ExcellentCount:  int32(result.ExcellentCount),
		NormalCount:     int32(result.NormalCount),
		WarningCount:    int32(result.WarningCount),
		CriticalCount:   int32(result.CriticalCount),
		OverallStatus:   result.OverallStatus,
		OverallMessage:  result.OverallMessage,
		Recommendations: pbRecommendations,
		CalculatedAt:    result.CalculatedAt,
	}, nil
}
