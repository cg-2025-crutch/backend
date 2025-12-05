package handler

import (
	"context"
	"errors"

	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	l := log.FromContext(ctx)

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	tokenPair, user, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		l.Error("Failed to login", zap.Error(err))
		return nil, status.Error(codes.Internal, "login failed")
	}

	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		User:         h.modelToProto(user),
	}, nil
}
