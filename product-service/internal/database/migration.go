package database

import (
	"log"

	"github.com/yourusername/product-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate các model
	err := db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.ProductVariant{},
		&models.VariantAttribute{},
	)
	if err != nil {
		return err
	}

	// Thêm dữ liệu mẫu
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		// Tạo categories
		categories := []models.Category{
			{
				Name:        "Electronics",
				Slug:        "electronics",
				Description: "Electronic devices and accessories",
			},
			{
				Name:        "Clothing",
				Slug:        "clothing",
				Description: "Fashion items and accessories",
			},
			{
				Name:        "Books",
				Slug:        "books",
				Description: "Books and publications",
			},
		}

		for _, category := range categories {
			if err := db.Create(&category).Error; err != nil {
				log.Printf("Error creating category %s: %v", category.Name, err)
			}
		}

		// Tạo products
		products := []models.Product{
			{
				Name:        "iPhone 13",
				Slug:        "iphone-13",
				Description: "Latest iPhone model with advanced features",
				Price:       999.99,
				Stock:       100,
				CategoryID:  1, // Electronics
				Variants: []models.ProductVariant{
					{
						Name:  "128GB",
						Price: 999.99,
						Stock: 50,
						Attributes: []models.VariantAttribute{
							{Name: "Storage", Value: "128GB"},
							{Name: "Color", Value: "Black"},
						},
					},
					{
						Name:  "256GB",
						Price: 1099.99,
						Stock: 50,
						Attributes: []models.VariantAttribute{
							{Name: "Storage", Value: "256GB"},
							{Name: "Color", Value: "Black"},
						},
					},
				},
			},
			{
				Name:        "T-Shirt",
				Slug:        "t-shirt",
				Description: "Comfortable cotton t-shirt",
				Price:       29.99,
				Stock:       200,
				CategoryID:  2, // Clothing
				Variants: []models.ProductVariant{
					{
						Name:  "Small",
						Price: 29.99,
						Stock: 100,
						Attributes: []models.VariantAttribute{
							{Name: "Size", Value: "S"},
							{Name: "Color", Value: "White"},
						},
					},
					{
						Name:  "Medium",
						Price: 29.99,
						Stock: 100,
						Attributes: []models.VariantAttribute{
							{Name: "Size", Value: "M"},
							{Name: "Color", Value: "White"},
						},
					},
				},
			},
			{
				Name:        "The Go Programming Language",
				Slug:        "the-go-programming-language",
				Description: "A comprehensive guide to Go programming",
				Price:       49.99,
				Stock:       50,
				CategoryID:  3, // Books
				Variants: []models.ProductVariant{
					{
						Name:  "Paperback",
						Price: 49.99,
						Stock: 25,
						Attributes: []models.VariantAttribute{
							{Name: "Format", Value: "Paperback"},
							{Name: "Language", Value: "English"},
						},
					},
					{
						Name:  "Hardcover",
						Price: 69.99,
						Stock: 25,
						Attributes: []models.VariantAttribute{
							{Name: "Format", Value: "Hardcover"},
							{Name: "Language", Value: "English"},
						},
					},
				},
			},
		}

		for _, product := range products {
			if err := db.Create(&product).Error; err != nil {
				log.Printf("Error creating product %s: %v", product.Name, err)
			}
		}
	}

	return nil
} 