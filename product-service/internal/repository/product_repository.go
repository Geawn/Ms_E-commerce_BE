package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/yourusername/product-service/internal/models"
)

type ProductRepository interface {
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]models.Product, error)
	ListByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]models.Product, error)
	Search(ctx context.Context, query string, limit, offset int) ([]models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepository(db *gorm.DB, redisClient *redis.Client) ProductRepository {
	return &productRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product

	// Try to get from cache first
	cacheKey := "product:" + slug
	if _, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
		// TODO: Unmarshal cached data
	}

	// If not in cache, get from database
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&product).Error; err != nil {
		return nil, err
	}

	// Cache the result
	// TODO: Marshal and cache product data

	return &product, nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]models.Product, error) {
	var products []models.Product

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) ListByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]models.Product, error) {
	var products []models.Product

	if err := r.db.WithContext(ctx).Joins("JOIN categories ON products.category_id = categories.id").
		Where("categories.slug = ?", categorySlug).
		Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]models.Product, error) {
	var products []models.Product

	if err := r.db.WithContext(ctx).Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).Offset(offset).Find(&products).Error; err != nil {
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
