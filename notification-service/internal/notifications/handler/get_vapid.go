package handler

import (
	"context"

	"github.com/cg-2025-crutch/backend/notification-service/internal/grpc/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *GRPCHandler) GetVapidKey(context.Context, *emptypb.Empty) (*gen.GetVapidKeyResponse, error) {
	vapidKey := h.service.GetVapidKey()

	return &gen.GetVapidKeyResponse{
		VapidKey: vapidKey,
	}, nil
}
