package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewMenuRepository(db *gorm.DB, redis *redis.Client) *MenuRepository {
	return &MenuRepository{
		db:    db,
		redis: redis,
	}
}

func (r *MenuRepository) GetBySlugAndChannel(ctx context.Context, slug, channel string) (*models.Menu, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("menu:%s:%s", slug, channel)
	cachedData, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var menu models.Menu
		if err := json.Unmarshal([]byte(cachedData), &menu); err == nil {
			return &menu, nil
		}
	}

	// If not in Redis, get from database
	var menu models.Menu
	if err := r.db.Preload("Items").
		Preload("Items.Category").
		Preload("Items.Collection").
		Preload("Items.Page").
		Preload("Items.Children").
		Preload("Items.Children.Category").
		Preload("Items.Children.Collection").
		Preload("Items.Children.Page").
		Where("slug = ? AND channel = ?", slug, channel).
		First(&menu).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	menuJSON, err := json.Marshal(menu)
	if err == nil {
		r.redis.Set(ctx, cacheKey, menuJSON, 24*time.Hour)
	}

	return &menu, nil
} 