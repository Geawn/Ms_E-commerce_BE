package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	CreateUser(ctx context.Context, email, firstName, lastName string) (uint64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	log.Printf("Repository: GetUserByID called with userID: %s", userID)
	
	var user models.User
	result := r.db.WithContext(ctx).First(&user, userID)
	if result.Error != nil {
		log.Printf("Repository: Error querying database: %v", result.Error)
		return nil, fmt.Errorf("failed to get user: %v", result.Error)
	}
	
	log.Printf("Repository: User found in database: %+v", user)
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, email, firstName, lastName string) (uint64, error) {
	log.Printf("Repository: CreateUser called with email: %s, firstName: %s, lastName: %s", email, firstName, lastName)
	
	user := &models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		log.Printf("Repository: Error creating user in database: %v", result.Error)
		return 0, fmt.Errorf("failed to create user: %v", result.Error)
	}
	
	log.Printf("Repository: User created in database with ID: %d", user.ID)
	return uint64(user.ID), nil
}
