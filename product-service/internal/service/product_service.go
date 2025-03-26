package service

import (
	"context"
	"strconv"

	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ProductService) ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]models.Product, error) {
	id, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByCategory(ctx, uint(id), limit, offset)
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]models.Product, error) {
	return s.repo.Search(ctx, query, limit, offset)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
