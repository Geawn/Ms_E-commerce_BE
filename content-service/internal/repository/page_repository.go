package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"gorm.io/gorm"
)

type PageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) GetBySlug(ctx context.Context, slug string) (*models.Page, error) {
	var page models.Page
	if err := r.db.Where("slug = ?", slug).First(&page).Error; err != nil {
		return nil, err
	}
	return &page, nil
} 