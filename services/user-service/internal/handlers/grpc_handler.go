package handlers

import (
	"context"

	pb "github.com/username/dist-ecommerce-go/proto/user"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/service"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.svc.CreateUser(ctx, req.Email, req.Password, req.FullName)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.svc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

func (h *UserHandler) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	user, err := h.svc.ValidateUser(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.ValidateUserResponse{Valid: false}, nil
	}

	return &pb.ValidateUserResponse{
		Valid: true,
		User: &pb.User{
			Id:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}
