package migration

import (
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate các model first
	err := db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.ProductVariant{},
		&models.VariantAttribute{},
		&models.Collection{},
		&models.ProductCollection{},
		&models.ProductAttribute{},
		&models.Review{},
		&models.User{},
		&models.Image{},
	)
	if err != nil {
		return err
	}

	// Add size and format columns to images table if they don't exist
	if err := db.Exec("ALTER TABLE images ADD COLUMN IF NOT EXISTS size integer").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE images ADD COLUMN IF NOT EXISTS format text").Error; err != nil {
		return err
	}

	// Then add default values for amount and currency columns
	if err := db.Exec("ALTER TABLE products ADD COLUMN IF NOT EXISTS amount decimal DEFAULT 0").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE products ADD COLUMN IF NOT EXISTS currency text DEFAULT 'USD'").Error; err != nil {
		return err
	}

	// Check and add categories if needed
	var categoryCount int64
	db.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
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
	}

	// Check and add collections if needed
	var collectionCount int64
	db.Model(&models.Collection{}).Count(&collectionCount)
	if collectionCount == 0 {
		collections := []models.Collection{
			{
				Name:        "Best Sellers",
				Slug:        "best-sellers",
				Description: "Our most popular products",
			},
			{
				Name:        "New Arrivals",
				Slug:        "new-arrivals",
				Description: "Recently added products",
			},
			{
				Name:        "Special Offers",
				Slug:        "special-offers",
				Description: "Products with special discounts",
			},
		}

		for _, collection := range collections {
			if err := db.Create(&collection).Error; err != nil {
				log.Printf("Error creating collection %s: %v", collection.Name, err)
			}
		}
	}

	// Check and add users if needed
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		users := []models.User{
			{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			{
				Name:  "Jane Smith",
				Email: "jane@example.com",
			},
		}

		for _, user := range users {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Error creating user %s: %v", user.Name, err)
			}
		}
	}

	// Check and add products if needed
	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		// Tạo products
		products := []models.Product{
			{
				Name:        "iPhone 13",
				Slug:        "iphone-13",
				Description: "Latest iPhone model with advanced features",
				CategoryID:  1, // Electronics
				Rating:      4.5,
				Pricing: models.ProductPricing{
					PriceRange: models.PriceRange{
						Start: models.Price{
							Amount:   999.99,
							Currency: "USD",
						},
						Stop: models.Price{
							Amount:   999.99,
							Currency: "USD",
						},
					},
				},
				Thumbnail: &models.Image{
					URL:    "https://example.com/iphone-13.jpg",
					Alt:    "iPhone 13",
					Size:   1024,
					Format: "WEBP",
				},
				Attributes: []*models.ProductAttribute{
					{
						Name:   "Color",
						Values: []string{"Black", "White", "Blue"},
					},
					{
						Name:   "Storage",
						Values: []string{"128GB", "256GB", "512GB"},
					},
				},
				Variants: []*models.ProductVariant{
					{
						Name:  "128GB Black",
						Stock: 50,
						Pricing: models.ProductPricing{
							PriceRange: models.PriceRange{
								Start: models.Price{
									Amount:   999.99,
									Currency: "USD",
								},
								Stop: models.Price{
									Amount:   999.99,
									Currency: "USD",
								},
							},
						},
						Attributes: []*models.VariantAttribute{
							{Name: "Storage", Value: "128GB"},
							{Name: "Color", Value: "Black"},
						},
					},
				},
			},
			{
				Name:        "T-Shirt",
				Slug:        "t-shirt",
				Description: "Comfortable cotton t-shirt",
				CategoryID:  2, // Clothing
				Rating:      4.0,
				Pricing: models.ProductPricing{
					PriceRange: models.PriceRange{
						Start: models.Price{
							Amount:   29.99,
							Currency: "USD",
						},
						Stop: models.Price{
							Amount:   29.99,
							Currency: "USD",
						},
					},
				},
				Thumbnail: &models.Image{
					URL:    "https://example.com/t-shirt.jpg",
					Alt:    "T-Shirt",
					Size:   1024,
					Format: "WEBP",
				},
				Attributes: []*models.ProductAttribute{
					{
						Name:   "Size",
						Values: []string{"S", "M", "L", "XL"},
					},
					{
						Name:   "Color",
						Values: []string{"White", "Black", "Blue"},
					},
				},
				Variants: []*models.ProductVariant{
					{
						Name:  "Small White",
						Stock: 100,
						Pricing: models.ProductPricing{
							PriceRange: models.PriceRange{
								Start: models.Price{
									Amount:   29.99,
									Currency: "USD",
								},
								Stop: models.Price{
									Amount:   29.99,
									Currency: "USD",
								},
							},
						},
						Attributes: []*models.VariantAttribute{
							{Name: "Size", Value: "S"},
							{Name: "Color", Value: "White"},
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

		// Add products to collections if needed
		var productCollectionCount int64
		db.Model(&models.ProductCollection{}).Count(&productCollectionCount)
		if productCollectionCount == 0 {
			productCollections := []models.ProductCollection{
				{ProductID: 1, CollectionID: 1}, // iPhone 13 -> Best Sellers
				{ProductID: 2, CollectionID: 1}, // T-Shirt -> Best Sellers
				{ProductID: 1, CollectionID: 2}, // iPhone 13 -> New Arrivals
				{ProductID: 2, CollectionID: 3}, // T-Shirt -> Special Offers
			}

			for _, pc := range productCollections {
				if err := db.Create(&pc).Error; err != nil {
					log.Printf("Error creating product collection: %v", err)
				}
			}
		}
	}

	return nil
}
