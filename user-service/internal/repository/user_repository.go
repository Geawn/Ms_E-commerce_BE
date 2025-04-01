package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	// Add your database client here
	// db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	// Implement your database query here
	// This is just a mock implementation
	return &models.User{
		ID:        id,
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Avatar: &models.Avatar{
			URL: "https://example.com/avatar.jpg",
			Alt: "User avatar",
		},
	}, nil
}
