package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetBySlugAndChannel(ctx context.Context, slug, channel string) (*models.Menu, error) {
	var menu models.Menu
	if err := r.db.Where("slug = ? AND channel = ?", slug, channel).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) GetItemsByMenuID(ctx context.Context, menuID uint) ([]*models.MenuItem, error) {
	var items []*models.MenuItem
	if err := r.db.Where("menu_id = ?", menuID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MenuRepository) GetChildrenByParentID(ctx context.Context, parentID uint) ([]*models.MenuItem, error) {
	var children []*models.MenuItem
	if err := r.db.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}
