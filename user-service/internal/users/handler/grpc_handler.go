package handler

import (
	pb "github.com/cg-2025-crutch/backend/user-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/models"
	"github.com/cg-2025-crutch/backend/user-service/internal/users/service"
)

type GRPCHandler struct {
	pb.UnimplementedUserServiceServer
	service service.UserService
}

func NewGRPCHandler(service service.UserService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) modelToProto(user *models.User) *pb.User {
	return &pb.User{
		Id:           user.UID.String(),
		Username:     user.Username,
		FirstName:    user.FirstName,
		SecondName:   user.SecondName,
		Age:          user.Age,
		Salary:       user.Salary,
		WorkSphereId: user.WorkSphereID,
	}
}
