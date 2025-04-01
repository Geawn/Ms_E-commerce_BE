package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
