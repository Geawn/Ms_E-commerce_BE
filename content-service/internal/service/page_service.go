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

type PageService struct {
	pageRepo *repository.PageRepository
	redis    *redis.Client
}

func NewPageService(pageRepo *repository.PageRepository, redis *redis.Client) *PageService {
	return &PageService{
		pageRepo: pageRepo,
		redis:    redis,
	}
}

func (s *PageService) GetBySlug(ctx context.Context, slug string) (*models.Page, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("page:%s", slug)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var page models.Page
		if err := json.Unmarshal([]byte(cached), &page); err == nil {
			return &page, nil
		}
	}

	// If not in Redis, get from DB
	page, err := s.pageRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Cache in Redis for 1 hour
	if pageJSON, err := json.Marshal(page); err == nil {
		s.redis.Set(ctx, cacheKey, pageJSON, time.Hour)
	}

	return page, nil
}

func (s *PageService) GetContent(ctx context.Context, slug string) (string, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("page:%s:content", slug)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		return cached, nil
	}

	// If not in Redis, get from DB
	page, err := s.pageRepo.GetBySlug(ctx, slug)
	if err != nil {
		return "", err
	}

	// Cache in Redis for 1 hour
	s.redis.Set(ctx, cacheKey, page.Content, time.Hour)

	return page.Content, nil
} 