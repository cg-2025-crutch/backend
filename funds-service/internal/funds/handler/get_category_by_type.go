package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetCategoriesByType(ctx context.Context, req *pb.GetCategoriesByTypeRequest) (*pb.GetCategoriesByTypeResponse, error) {
	categories, err := h.service.GetCategoriesByType(ctx, req.Type)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get categories: %v", err)
	}

	pbCategories := make([]*pb.Category, 0, len(categories))
	for _, c := range categories {
		pbCategories = append(pbCategories, h.categoryToProto(c))
	}

	return &pb.GetCategoriesByTypeResponse{
		Categories: pbCategories,
	}, nil
}
