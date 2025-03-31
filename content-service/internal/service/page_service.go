package service

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/repository"
)

type PageService struct {
	pageRepo *repository.PageRepository
}

func NewPageService(pageRepo *repository.PageRepository) *PageService {
	return &PageService{
		pageRepo: pageRepo,
	}
}

func (s *PageService) GetBySlug(ctx context.Context, slug string) (*models.Page, error) {
	return s.pageRepo.GetBySlug(ctx, slug)
} 