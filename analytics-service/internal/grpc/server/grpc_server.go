package server

import (
	"context"
	"fmt"
	"net"

	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/handler"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/config"
	pb "github.com/cg-2025-crutch/backend/analytics-service/internal/grpc/gen/analytics_service"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	cfg    config.ServerConfig
	server *grpc.Server
}

func NewGrpcServer(cfg config.ServerConfig) *GrpcServer {
	return &GrpcServer{
		cfg:    cfg,
		server: grpc.NewServer(),
	}
}

func (g *GrpcServer) RegisterGRPC(handler *handler.AnalyticsHandler) {
	pb.RegisterAnalyticsServiceServer(g.server, handler)
}

func (g *GrpcServer) Run(ctx context.Context) error {
	l := log.FromContext(ctx)

	addr := fmt.Sprintf("%s:%s", g.cfg.GRPCHost, g.cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	l.Infof("gRPC server starting on %s", addr)

	errChan := make(chan error, 1)

	go func() {
		if err := g.server.Serve(lis); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		l.Info("Shutting down gRPC server...")
		g.server.GracefulStop()
		return nil
	case err := <-errChan:
		return err
	}
}
