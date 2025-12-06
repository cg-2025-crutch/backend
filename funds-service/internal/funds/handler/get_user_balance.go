package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetUserBalance(ctx context.Context, req *pb.GetUserBalanceRequest) (*pb.GetUserBalanceResponse, error) {
	balance, err := h.service.GetUserBalance(ctx, req.UserUid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	return &pb.GetUserBalanceResponse{
		Balance: h.balanceToProto(balance),
	}, nil
}
