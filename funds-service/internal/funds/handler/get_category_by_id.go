package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetCategoryById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.GetCategoryByIdResponse, error) {
	category, err := h.service.GetCategoryById(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "category not found: %v", err)
	}

	return &pb.GetCategoryByIdResponse{
		Category: h.categoryToProto(category),
	}, nil
}
