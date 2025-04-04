package grpc

import (
	"context"
	"fmt"
	"log"

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

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Received CreateUser request - Email: %s, FirstName: %s, LastName: %s", req.Email, req.FirstName, req.LastName)

	userID, err := s.userService.CreateUser(ctx, req.Email, req.FirstName, req.LastName)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return &pb.CreateUserResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	log.Printf("Successfully created user with ID: %d", userID)
	return &pb.CreateUserResponse{
		UserId:  userID,
		Success: true,
	}, nil
}
