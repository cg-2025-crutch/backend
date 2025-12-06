package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetTransactionById(ctx context.Context, req *pb.GetTransactionByIdRequest) (*pb.GetTransactionByIdResponse, error) {
	transaction, err := h.service.GetTransactionById(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "transaction not found: %v", err)
	}

	return &pb.GetTransactionByIdResponse{
		Transaction: h.transactionToProto(transaction),
	}, nil
}
