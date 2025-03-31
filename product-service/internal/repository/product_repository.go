package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]*models.Product, error)
	ListByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]*models.Product, error)
	ListByCollection(ctx context.Context, collectionID uint, limit, offset int) ([]*models.Product, error)
	Search(ctx context.Context, query string, limit, offset int, sortBy string, sortDirection string) ([]*models.Product, error)
	GetTotalSearchResults(ctx context.Context, query string) (int, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	CreateCollection(ctx context.Context, collection *models.Collection) error
	UpdateCollection(ctx context.Context, collection *models.Collection) error
	DeleteCollection(ctx context.Context, id uint) error
	GetCollectionBySlug(ctx context.Context, slug string) (*models.Collection, error)
	AddProductToCollection(ctx context.Context, productID, collectionID uint) error
	RemoveProductFromCollection(ctx context.Context, productID, collectionID uint) error
	CreateReview(ctx context.Context, review *models.Review) error
	UpdateReview(ctx context.Context, review *models.Review) error
	DeleteReview(ctx context.Context, id uint) error
	GetProductReviews(ctx context.Context, productID uint) ([]*models.Review, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error)
	GetPosition(ctx context.Context, productID uint) (int, error)
	GetTotal(ctx context.Context) (int, error)
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
		Preload("Attributes").
		Preload("Collections").
		Preload("Reviews").
		Preload("Thumbnail").
		Where("slug = ?", slug).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Attributes").
		Preload("Collections").
		Preload("Reviews").
		Preload("Thumbnail").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Attributes").
		Preload("Collections").
		Preload("Reviews").
		Preload("Thumbnail").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	result := make([]*models.Product, len(products))
	for i := range products {
		result[i] = &products[i]
	}
	return result, nil
}

func (r *productRepository) ListByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]*models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Attributes").
		Preload("Collections").
		Preload("Reviews").
		Preload("Thumbnail").
		Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	result := make([]*models.Product, len(products))
	for i := range products {
		result[i] = &products[i]
	}
	return result, nil
}

func (r *productRepository) ListByCollection(ctx context.Context, collectionID uint, limit, offset int) ([]*models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Joins("JOIN product_collections ON products.id = product_collections.product_id").
		Where("product_collections.collection_id = ?", collectionID).
		Preload("Category").
		Preload("Variants").
		Preload("Attributes").
		Preload("Collections").
		Preload("Reviews").
		Preload("Thumbnail").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	result := make([]*models.Product, len(products))
	for i := range products {
		result[i] = &products[i]
	}
	return result, nil
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int, sortBy string, sortDirection string) ([]*models.Product, error) {
	var products []models.Product
	db := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Attributes").
		Preload("Collections").
		Preload("Reviews").
		Preload("Thumbnail").
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")

	// Apply sorting
	switch sortBy {
	case "NAME":
		db = db.Order("name " + sortDirection)
	case "PRICE":
		// Assuming price is stored in variants
		db = db.Joins("LEFT JOIN product_variants ON products.id = product_variants.product_id").
			Order("product_variants.price " + sortDirection)
	case "CREATED_AT":
		db = db.Order("created_at " + sortDirection)
	case "UPDATED_AT":
		db = db.Order("updated_at " + sortDirection)
	default:
		db = db.Order("created_at DESC")
	}

	err := db.Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	result := make([]*models.Product, len(products))
	for i := range products {
		result[i] = &products[i]
	}
	return result, nil
}

func (r *productRepository) GetTotalSearchResults(ctx context.Context, query string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
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

func (r *productRepository) CreateCollection(ctx context.Context, collection *models.Collection) error {
	return r.db.WithContext(ctx).Create(collection).Error
}

func (r *productRepository) UpdateCollection(ctx context.Context, collection *models.Collection) error {
	return r.db.WithContext(ctx).Save(collection).Error
}

func (r *productRepository) DeleteCollection(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Collection{}, id).Error
}

func (r *productRepository) GetCollectionBySlug(ctx context.Context, slug string) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.WithContext(ctx).
		Preload("Products").
		Where("slug = ?", slug).
		First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *productRepository) AddProductToCollection(ctx context.Context, productID, collectionID uint) error {
	return r.db.WithContext(ctx).Create(&models.ProductCollection{
		ProductID:    productID,
		CollectionID: collectionID,
	}).Error
}

func (r *productRepository) RemoveProductFromCollection(ctx context.Context, productID, collectionID uint) error {
	return r.db.WithContext(ctx).
		Where("product_id = ? AND collection_id = ?", productID, collectionID).
		Delete(&models.ProductCollection{}).Error
}

func (r *productRepository) CreateReview(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *productRepository) UpdateReview(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *productRepository) DeleteReview(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}

func (r *productRepository) GetProductReviews(ctx context.Context, productID uint) ([]*models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("product_id = ?", productID).
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	result := make([]*models.Review, len(reviews))
	for i := range reviews {
		result[i] = &reviews[i]
	}
	return result, nil
}

func (r *productRepository) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Preload("Products").
		Where("slug = ?", slug).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *productRepository) GetPosition(ctx context.Context, productID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id <= ?", productID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count) - 1, nil
}

func (r *productRepository) GetTotal(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
