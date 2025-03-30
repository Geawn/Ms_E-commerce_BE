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

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ProductService) ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*models.Product, error) {
	id, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByCategory(ctx, uint(id), limit, offset)
}

func (s *ProductService) ListByCollection(ctx context.Context, collectionID string, limit, offset int) ([]*models.Product, error) {
	id, err := strconv.ParseUint(collectionID, 10, 32)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByCollection(ctx, uint(id), limit, offset)
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*models.Product, error) {
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

func (s *ProductService) CreateCollection(ctx context.Context, collection *models.Collection) error {
	return s.repo.CreateCollection(ctx, collection)
}

func (s *ProductService) UpdateCollection(ctx context.Context, collection *models.Collection) error {
	return s.repo.UpdateCollection(ctx, collection)
}

func (s *ProductService) DeleteCollection(ctx context.Context, id uint) error {
	return s.repo.DeleteCollection(ctx, id)
}

func (s *ProductService) GetCollectionBySlug(ctx context.Context, slug string) (*models.Collection, error) {
	return s.repo.GetCollectionBySlug(ctx, slug)
}

func (s *ProductService) AddProductToCollection(ctx context.Context, productID, collectionID uint) error {
	return s.repo.AddProductToCollection(ctx, productID, collectionID)
}

func (s *ProductService) RemoveProductFromCollection(ctx context.Context, productID, collectionID uint) error {
	return s.repo.RemoveProductFromCollection(ctx, productID, collectionID)
}

func (s *ProductService) CreateReview(ctx context.Context, review *models.Review) error {
	return s.repo.CreateReview(ctx, review)
}

func (s *ProductService) UpdateReview(ctx context.Context, review *models.Review) error {
	return s.repo.UpdateReview(ctx, review)
}

func (s *ProductService) DeleteReview(ctx context.Context, id uint) error {
	return s.repo.DeleteReview(ctx, id)
}

func (s *ProductService) GetProductReviews(ctx context.Context, productID uint) ([]*models.Review, error) {
	return s.repo.GetProductReviews(ctx, productID)
}

func (s *ProductService) UpdateProductRating(ctx context.Context, productID uint) error {
	reviews, err := s.repo.GetProductReviews(ctx, productID)
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		return nil
	}

	var totalRating float64
	for _, review := range reviews {
		totalRating += review.Rating
	}

	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	product.Rating = totalRating / float64(len(reviews))
	return s.repo.Update(ctx, product)
}

func (s *ProductService) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	return s.repo.GetCategoryBySlug(ctx, slug)
}
