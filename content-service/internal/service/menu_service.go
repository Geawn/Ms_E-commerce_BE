package service

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
}

func NewMenuService(menuRepo *repository.MenuRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
	}
}

func (s *MenuService) GetBySlugAndChannel(ctx context.Context, slug, channel string) (*models.Menu, error) {
	return s.menuRepo.GetBySlugAndChannel(ctx, slug, channel)
} 