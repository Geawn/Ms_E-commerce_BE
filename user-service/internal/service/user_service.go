package service

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/repository"
)

type UserService interface {
	GetCurrentUser(ctx context.Context, userID string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetCurrentUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}
