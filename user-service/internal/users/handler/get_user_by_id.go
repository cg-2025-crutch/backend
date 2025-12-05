package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user uid is required")
	}

	user, err := h.service.GetUserByID(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserByIdResponse{
		User: h.modelToProto(user),
	}, nil
}
