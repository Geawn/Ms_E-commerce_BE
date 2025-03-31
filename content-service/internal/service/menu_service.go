package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/repository"
	"github.com/go-redis/redis/v8"
)

type MenuService struct {
	menuRepo       *repository.MenuRepository
	pageRepo       *repository.PageRepository
	categoryRepo   *repository.CategoryRepository
	collectionRepo *repository.CollectionRepository
	redis          *redis.Client
}

func NewMenuService(
	menuRepo *repository.MenuRepository,
	pageRepo *repository.PageRepository,
	categoryRepo *repository.CategoryRepository,
	collectionRepo *repository.CollectionRepository,
	redis *redis.Client,
) *MenuService {
	return &MenuService{
		menuRepo:       menuRepo,
		pageRepo:       pageRepo,
		categoryRepo:   categoryRepo,
		collectionRepo: collectionRepo,
		redis:          redis,
	}
}

func (s *MenuService) GetBySlugAndChannel(ctx context.Context, slug, channel string) (*models.Menu, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("menu:%s:%s", slug, channel)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var menu models.Menu
		if err := json.Unmarshal([]byte(cached), &menu); err == nil {
			return &menu, nil
		}
	}

	// If not in Redis, get from DB
	menu, err := s.menuRepo.GetBySlugAndChannel(ctx, slug, channel)
	if err != nil {
		return nil, err
	}

	// Cache in Redis for 1 hour
	if menuJSON, err := json.Marshal(menu); err == nil {
		s.redis.Set(ctx, cacheKey, menuJSON, time.Hour)
	}

	return menu, nil
}

func (s *MenuService) GetItemsByMenuID(ctx context.Context, menuID uint) ([]*models.MenuItem, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("menu:%d:items", menuID)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var items []*models.MenuItem
		if err := json.Unmarshal([]byte(cached), &items); err == nil {
			return items, nil
		}
	}

	// If not in Redis, get from DB
	items, err := s.menuRepo.GetItemsByMenuID(ctx, menuID)
	if err != nil {
		return nil, err
	}

	// Cache in Redis for 1 hour
	if itemsJSON, err := json.Marshal(items); err == nil {
		s.redis.Set(ctx, cacheKey, itemsJSON, time.Hour)
	}

	return items, nil
}

func (s *MenuService) GetCategoryByID(ctx context.Context, id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *MenuService) GetCollectionByID(ctx context.Context, id uint) (*models.Collection, error) {
	return s.collectionRepo.GetByID(ctx, id)
}

func (s *MenuService) GetPageBySlug(ctx context.Context, slug string) (*models.Page, error) {
	return s.pageRepo.GetBySlug(ctx, slug)
}

func (s *MenuService) GetChildrenByParentID(ctx context.Context, parentID uint) ([]*models.MenuItem, error) {
	return s.menuRepo.GetChildrenByParentID(ctx, parentID)
} 