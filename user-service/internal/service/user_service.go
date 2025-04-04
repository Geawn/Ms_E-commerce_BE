package service

import (
	"context"
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/repository"
)

type UserService interface {
	GetCurrentUser(ctx context.Context, userID string) (*models.User, error)
	CreateUser(ctx context.Context, email, firstName, lastName string) (uint64, error)
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
	log.Printf("Service: GetCurrentUser called with userID: %s", userID)
	
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Service: Error getting user from repository: %v", err)
		return nil, err
	}
	
	log.Printf("Service: User found in repository: %+v", user)
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, email, firstName, lastName string) (uint64, error) {
	log.Printf("Service: CreateUser called with email: %s, firstName: %s, lastName: %s", email, firstName, lastName)
	
	userID, err := s.userRepo.CreateUser(ctx, email, firstName, lastName)
	if err != nil {
		log.Printf("Service: Error creating user in repository: %v", err)
		return 0, err
	}
	
	log.Printf("Service: User created with ID: %d", userID)
	return userID, nil
}
