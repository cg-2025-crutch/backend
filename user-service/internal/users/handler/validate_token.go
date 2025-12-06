package handler

import (
	"context"
	"errors"

	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/jwt"
	"github.com/cg-2025-crutch/backend/user-service/internal/infrastructure/log"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	l := log.FromContext(ctx)

	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	claims, err := h.service.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return &pb.ValidateTokenResponse{
				Valid: false,
			}, nil
		}
		if errors.Is(err, jwt.ErrInvalidToken) {
			return &pb.ValidateTokenResponse{
				Valid: false,
			}, nil
		}
		l.Error("Failed to validate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "token validation failed")
	}

	return &pb.ValidateTokenResponse{
		Valid:     true,
		UserId:    claims.UserID.String(),
		Username:  claims.Username,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}
