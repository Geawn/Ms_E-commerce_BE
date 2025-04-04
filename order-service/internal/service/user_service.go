package service

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/proto"
)

type UserService struct {
	client proto.UserServiceClient
}

func NewUserService(client proto.UserServiceClient) *UserService {
	return &UserService{
		client: client,
	}
}

func (s *UserService) GetCurrentUser(ctx context.Context, token string) (*proto.GetCurrentUserResponse, error) {
	resp, err := s.client.GetCurrentUser(ctx, &proto.GetCurrentUserRequest{
		Token: token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	return resp, nil
}
