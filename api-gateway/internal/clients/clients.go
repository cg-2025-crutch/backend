package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cg-2025-crutch/backend/api-gateway/internal/config"
	analytics_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/analytics_service"
	funds_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/funds_service"
	notif_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/notification_service"
	user_pb "github.com/cg-2025-crutch/backend/api-gateway/internal/grpc/gen/user_service"
)

type GRPCClients struct {
	UserService      user_pb.UserServiceClient
	FundsService     funds_pb.FundsServiceClient
	NotifService     notif_pb.NotificationServiceClient
	AnalyticsService analytics_pb.AnalyticsServiceClient

	userConn      *grpc.ClientConn
	fundsConn     *grpc.ClientConn
	notifConn     *grpc.ClientConn
	analyticsConn *grpc.ClientConn
}

func NewGRPCClients(cfg config.AppConfig) (*GRPCClients, error) {
	clients := &GRPCClients{}

	// Connect to User Service
	userAddr := fmt.Sprintf("%s:%s", cfg.UserServiceClient.Host, cfg.UserServiceClient.Port)
	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	clients.userConn = userConn
	clients.UserService = user_pb.NewUserServiceClient(userConn)

	// Connect to Funds Service
	fundsAddr := fmt.Sprintf("%s:%s", cfg.FundsServiceClient.Host, cfg.FundsServiceClient.Port)
	fundsConn, err := grpc.NewClient(
		fundsAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		clients.Close()
		return nil, fmt.Errorf("failed to connect to funds service: %w", err)
	}
	clients.fundsConn = fundsConn
	clients.FundsService = funds_pb.NewFundsServiceClient(fundsConn)

	// Connect to Notification Service
	notifAddr := fmt.Sprintf("%s:%s", cfg.NotifServiceClient.Host, cfg.NotifServiceClient.Port)
	notifConn, err := grpc.NewClient(
		notifAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		clients.Close()
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}
	clients.notifConn = notifConn
	clients.NotifService = notif_pb.NewNotificationServiceClient(notifConn)

	// Connect to Analytics Service
	analyticsAddr := fmt.Sprintf("%s:%s", cfg.AnalyticsServiceClient.Host, cfg.AnalyticsServiceClient.Port)
	analyticsConn, err := grpc.NewClient(
		analyticsAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		clients.Close()
		return nil, fmt.Errorf("failed to connect to analytics service: %w", err)
	}
	clients.analyticsConn = analyticsConn
	clients.AnalyticsService = analytics_pb.NewAnalyticsServiceClient(analyticsConn)

	return clients, nil
}

func (c *GRPCClients) Close() {
	if c.userConn != nil {
		_ = c.userConn.Close()
	}
	if c.fundsConn != nil {
		_ = c.fundsConn.Close()
	}
	if c.notifConn != nil {
		_ = c.notifConn.Close()
	}
	if c.analyticsConn != nil {
		_ = c.analyticsConn.Close()
	}
}

func (c *GRPCClients) ValidateToken(ctx context.Context, token string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.UserService.ValidateToken(ctx, &user_pb.ValidateTokenRequest{
		AccessToken: token,
	})
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	if !resp.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return resp.UserId, nil
}
