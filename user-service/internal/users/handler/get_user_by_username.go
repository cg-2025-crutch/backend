package handler

import (
	"context"

	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "empty username")
	}

	user, err := h.service.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByUsernameResponse{
		User: h.modelToProto(user),
	}, nil
}
