package migration

import (
	"log"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate các model first
	err := db.AutoMigrate(
		&model.Order{},
		&model.OrderLine{},
	)
	if err != nil {
		return err
	}

	// Add default values for amount and currency columns
	if err := db.Exec("ALTER TABLE orders ADD COLUMN IF NOT EXISTS total_amount decimal DEFAULT 0").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE orders ADD COLUMN IF NOT EXISTS currency text DEFAULT 'USD'").Error; err != nil {
		return err
	}

	// Add default values for price and currency in order_lines
	if err := db.Exec("ALTER TABLE order_lines ADD COLUMN IF NOT EXISTS price decimal DEFAULT 0").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE order_lines ADD COLUMN IF NOT EXISTS currency text DEFAULT 'USD'").Error; err != nil {
		return err
	}

	// Check and add sample orders if needed
	var orderCount int64
	db.Model(&model.Order{}).Count(&orderCount)
	if orderCount == 0 {
		// Tạo sample orders
		orders := []model.Order{
			{
				ID:            "1",
				UserID:        "1",
				Number:        "ORD-001",
				Created:       time.Now(),
				TotalAmount:   1029.98,
				Currency:      "USD",
				PaymentStatus: "PAID",
				Lines: []*model.OrderLine{
					{
						ID:          "1",
						OrderID:     "1",
						VariantID:   "1",
						VariantName: "128GB Black",
						ProductSlug: "iphone-13",
						Price:       999.99,
						Currency:    "USD",
						Quantity:    1,
					},
					{
						ID:          "2",
						OrderID:     "1",
						VariantID:   "2",
						VariantName: "Small White",
						ProductSlug: "t-shirt",
						Price:       29.99,
						Currency:    "USD",
						Quantity:    1,
					},
				},
			},
			{
				ID:            "2",
				UserID:        "2",
				Number:        "ORD-002",
				Created:       time.Now(),
				TotalAmount:   59.98,
				Currency:      "USD",
				PaymentStatus: "PENDING",
				Lines: []*model.OrderLine{
					{
						ID:          "3",
						OrderID:     "2",
						VariantID:   "2",
						VariantName: "Small White",
						ProductSlug: "t-shirt",
						Price:       29.99,
						Currency:    "USD",
						Quantity:    2,
					},
				},
			},
		}

		for _, order := range orders {
			if err := db.Create(&order).Error; err != nil {
				log.Printf("Error creating order %s: %v", order.Number, err)
			}
		}
	}

	return nil
}
