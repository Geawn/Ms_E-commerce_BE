package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]models.Product, error)
	ListByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]models.Product, error)
	Search(ctx context.Context, query string, limit, offset int) ([]models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Where("slug = ?", slug).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) ListByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
