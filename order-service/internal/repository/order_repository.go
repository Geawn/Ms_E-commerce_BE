package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	GetTotalByUserID(ctx context.Context, userID string) (int, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Lines").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).
		Preload("Lines").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	result := make([]*model.Order, len(orders))
	for i := range orders {
		result[i] = &orders[i]
	}
	return result, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *orderRepository) GetTotalByUserID(ctx context.Context, userID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
