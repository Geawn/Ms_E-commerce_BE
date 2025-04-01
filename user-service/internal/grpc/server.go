package grpc

import (
	"context"

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

	return &pb.GetCurrentUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar: &pb.Avatar{
				Url: user.Avatar.URL,
				Alt: user.Avatar.Alt,
			},
		},
	}, nil
}
