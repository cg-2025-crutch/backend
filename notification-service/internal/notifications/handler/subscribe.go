package handler

import (
	"context"

	"github.com/cg-2025-crutch/backend/notification-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) Subscribe(ctx context.Context, req *gen.SubscribeReq) (*gen.SubscribeResp, error) {
	l := log.FromContext(ctx)

	err := h.service.SubscribeUser(ctx, req.UserId, req.Endpoint, req.P256Dh, req.Auth)
	if err != nil {
		l.Errorf("failed to subscribe user: %s", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &gen.SubscribeResp{
		Message: "success",
	}, nil
}
