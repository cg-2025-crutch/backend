package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetUserTransactions(ctx context.Context, req *pb.GetUserTransactionsRequest) (*pb.GetUserTransactionsResponse, error) {
	limit := req.Limit
	if limit == 0 {
		limit = 50
	}

	transactions, total, err := h.service.GetUserTransactions(ctx, req.UserUid, limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	pbTransactions := make([]*pb.Transaction, 0, len(transactions))
	for _, t := range transactions {
		pbTransactions = append(pbTransactions, h.transactionToProto(t))
	}

	return &pb.GetUserTransactionsResponse{
		Transactions: pbTransactions,
		Total:        total,
	}, nil
}
