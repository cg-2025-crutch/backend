package service

import (
	"github.com/cg-2025-crutch/backend/notification-service/internal/config"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/repository"
)

type NotificationService struct {
	repo            repository.RedisRepo
	VapidPublicKey  string
	VapidPrivateKey string
	subscriber      string
}

func NewNotificationService(cfg config.NotificationsConfig, repo repository.RedisRepo) NotificationService {
	return NotificationService{
		repo:            repo,
		VapidPublicKey:  cfg.VapidPublic,
		VapidPrivateKey: cfg.VapidPrivate,
	}
}
