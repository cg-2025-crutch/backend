package handler

import (
	"context"
	"errors"

	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	l := log.FromContext(ctx)

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	dto := models.CreateUserDTO{
		Username:     req.Username,
		Password:     req.Password,
		FirstName:    req.FirstName,
		SecondName:   req.SecondName,
		Age:          req.Age,
		Salary:       req.Salary,
		WorkSphereID: req.WorkSphereId,
	}

	user, err := h.service.CreateUser(ctx, dto)
	if err != nil {
		if errors.Is(err, repository.ErrUsernameExists) {
			return nil, status.Error(codes.AlreadyExists, "username already exists")
		}
		l.Error("Failed to create user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.CreateUserResponse{
		User: h.modelToProto(user),
	}, nil
}
