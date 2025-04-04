package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/model"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/repository"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo      repository.OrderRepository
	userService    *UserService
	productService *ProductService
	db             *gorm.DB
	redis          *redis.Client
}

func NewOrderService(orderRepo repository.OrderRepository, userService *UserService, productService *ProductService, db *gorm.DB, redis *redis.Client) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		userService:    userService,
		productService: productService,
		db:             db,
		redis:          redis,
	}
}

// CalculateTotal calculates the total amount for an order based on current product prices
func (s *OrderService) CalculateTotal(ctx context.Context, order *model.Order) error {
	var totalAmount float64
	var currency string

	for _, line := range order.Lines {
		// Get current price from product service
		product, err := s.productService.GetProductDetails(ctx, line.ProductSlug, "default")
		if err != nil {
			return fmt.Errorf("failed to get product details: %w", err)
		}

		// Find the correct variant
		for _, variant := range product.Variants {
			if variant.ID == line.VariantID {
				price, err := strconv.ParseFloat(variant.Pricing.Price.Gross.Amount, 64)
				if err != nil {
					return fmt.Errorf("failed to parse price: %w", err)
				}
				quantity := line.Quantity
				totalAmount += price * float64(quantity)
				currency = variant.Pricing.Price.Gross.Currency
				break
			}
		}
	}

	// Update total amount
	order.TotalAmount = totalAmount
	order.Currency = currency

	return nil
}

// UpdateOrder updates an order and recalculates the total
func (s *OrderService) UpdateOrder(ctx context.Context, orderID string, updates map[string]interface{}) error {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get current order
	var order model.Order
	if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update specified fields
	if err := tx.Model(&order).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Recalculate total
	if err := s.CalculateTotal(ctx, &order); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to calculate total: %w", err)
	}

	// Save order with new total
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save order: %w", err)
	}

	// Clear cache
	cacheKey := fmt.Sprintf("user_orders:%s", order.UserID)
	s.redis.Del(ctx, cacheKey)

	return tx.Commit().Error
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string, page, perPage int) ([]*model.Order, int, error) {
	offset := (page - 1) * perPage
	orders, err := s.orderRepo.GetByUserID(ctx, userID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.orderRepo.GetTotalByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, lines []*model.OrderLine) (*model.Order, error) {
	// Get user info
	_, err := s.userService.GetCurrentUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Create order
	order := &model.Order{
		UserID:        userID,
		Number:        generateOrderNumber(),
		Created:       time.Now(),
		PaymentStatus: model.PaymentStatusPending,
		Lines:         lines,
	}

	// Calculate total amount
	totalAmount := 0.0
	for _, line := range lines {
		// Get product details
		product, err := s.productService.GetProductDetails(ctx, line.ProductSlug, "default")
		if err != nil {
			return nil, fmt.Errorf("failed to get product details: %w", err)
		}

		// Find matching variant
		var variant *model.ProductVariant
		for _, v := range product.Variants {
			if v.ID == line.VariantID {
				variant = v
				break
			}
		}
		if variant == nil {
			return nil, fmt.Errorf("variant not found: %s", line.VariantID)
		}

		// Calculate line total
		price, err := strconv.ParseFloat(variant.Pricing.Price.Gross.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price: %w", err)
		}
		lineTotal := price * float64(line.Quantity)
		totalAmount += lineTotal

		// Set line price and currency
		line.Price = price
		line.Currency = variant.Pricing.Price.Gross.Currency
		line.VariantName = variant.Name
	}

	order.TotalAmount = totalAmount
	order.Currency = "USD" // Default currency

	// Save order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status model.PaymentStatusEnum) error {
	return s.orderRepo.UpdateStatus(ctx, id, string(status))
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}
