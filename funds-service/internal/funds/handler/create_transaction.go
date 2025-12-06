package handler

import (
	"context"
	"time"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction date format: %v", err)
	}

	input := models.CreateTransactionInput{
		UserUID:         req.UserUid,
		CategoryID:      req.CategoryId,
		Type:            req.Type,
		Amount:          req.Amount,
		Title:           req.Title,
		Description:     req.Description,
		TransactionDate: transactionDate,
	}

	transaction, err := h.service.CreateTransaction(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	return &pb.CreateTransactionResponse{
		Transaction: h.transactionToProto(transaction),
	}, nil
}
