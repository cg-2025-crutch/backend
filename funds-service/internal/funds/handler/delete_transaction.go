package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.DeleteTransactionResponse, error) {
	err := h.service.DeleteTransaction(ctx, req.Id, req.UserUid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete transaction: %v", err)
	}

	return &pb.DeleteTransactionResponse{
		Success: true,
	}, nil
}
