package service

import (
	"context"

	"github.com/yourusername/product-service/internal/models"
	"github.com/yourusername/product-service/internal/repository"
)

type ProductService interface {
	GetProductBySlug(ctx context.Context, slug string) (*models.Product, error)
	ListProducts(ctx context.Context, limit, offset int) ([]models.Product, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}

type productService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *productService) ListProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *productService) ListProductsByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]models.Product, error) {
	return s.repo.ListByCategory(ctx, categorySlug, limit, offset)
}

func (s *productService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]models.Product, error) {
	return s.repo.Search(ctx, query, limit, offset)
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
