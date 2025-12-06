package handler

import (
	"context"
	"errors"

	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	l := log.FromContext(ctx)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	dto := models.UpdateUserDTO{
		UID:          uid,
		Username:     req.Username,
		FirstName:    req.FirstName,
		SecondName:   req.SecondName,
		Age:          req.Age,
		Salary:       req.Salary,
		WorkSphereID: req.WorkSphereId,
	}

	user, err := h.service.UpdateUser(ctx, dto)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		l.Error("Failed to update user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.UpdateUserResponse{
		User: h.modelToProto(user),
	}, nil
}
