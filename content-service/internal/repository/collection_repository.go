package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) GetByID(ctx context.Context, id uint) (*models.Collection, error) {
	var collection models.Collection
	if err := r.db.First(&collection, id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
} 