package grpc

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/service"
	pb "github.com/Geawn/Ms_E-commerce_BE/user-service/proto"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}

func (s *Server) GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	user, err := s.userService.GetCurrentUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Convert uint ID to string
	id := fmt.Sprintf("%d", user.ID)

	// Create avatar if exists
	var avatar *pb.Avatar
	if user.Profile.Avatar != nil {
		avatar = &pb.Avatar{
			Url: user.Profile.Avatar.URL,
			Alt: user.Profile.Avatar.Alt,
		}
	}

	return &pb.GetCurrentUserResponse{
		User: &pb.User{
			Id:        id,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    avatar,
		},
	}, nil
}
