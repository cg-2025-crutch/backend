package clients

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	fundspb "github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/gen/funds_service"
	userpb "github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/gen/user_service"
)

// Clients содержит клиенты для всех внешних gRPC сервисов
type Clients struct {
	UserClient  userpb.UserServiceClient
	FundsClient fundspb.FundsServiceClient

	userConn  *grpc.ClientConn
	fundsConn *grpc.ClientConn
}

// NewClients создает новые gRPC клиенты
func NewClients(userServiceAddr, fundsServiceAddr string) (*Clients, error) {
	clients := &Clients{}

	// Connect to User Service
	userConn, err := grpc.NewClient(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-service: %w", err)
	}
	clients.userConn = userConn
	clients.UserClient = userpb.NewUserServiceClient(userConn)

	// Connect to Funds Service
	fundsConn, err := grpc.NewClient(
		fundsServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		clients.Close()
		return nil, fmt.Errorf("failed to connect to funds-service: %w", err)
	}
	clients.fundsConn = fundsConn
	clients.FundsClient = fundspb.NewFundsServiceClient(fundsConn)

	return clients, nil
}

// Close закрывает все gRPC соединения
func (c *Clients) Close() {
	if c.userConn != nil {
		_ = c.userConn.Close()
	}
	if c.fundsConn != nil {
		_ = c.fundsConn.Close()
	}
}
