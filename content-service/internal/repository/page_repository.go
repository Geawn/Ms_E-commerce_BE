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

type PageRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPageRepository(db *gorm.DB, redis *redis.Client) *PageRepository {
	return &PageRepository{
		db:    db,
		redis: redis,
	}
}

func (r *PageRepository) GetBySlug(ctx context.Context, slug string) (*models.Page, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("page:%s", slug)
	cachedData, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var page models.Page
		if err := json.Unmarshal([]byte(cachedData), &page); err == nil {
			return &page, nil
		}
	}

	// If not in Redis, get from database
	var page models.Page
	if err := r.db.Where("slug = ?", slug).First(&page).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	pageJSON, err := json.Marshal(page)
	if err == nil {
		r.redis.Set(ctx, cacheKey, pageJSON, 24*time.Hour)
	}

	return &page, nil
} 