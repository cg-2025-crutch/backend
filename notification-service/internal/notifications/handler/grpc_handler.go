package handler

import (
	"github.com/cg-2025-crutch/backend/notification-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/service"
)

type GRPCHandler struct {
	gen.UnimplementedNotificationServiceServer
	service service.NotificationService
}

func NewGRPCHandler(service service.NotificationService) *GRPCHandler {
	return &GRPCHandler{service: service}
}
