package migration

import (
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate các model first
	err := db.AutoMigrate(
		&models.Page{},
		&models.Menu{},
		&models.MenuItem{},
		&models.Category{},
		&models.Collection{},
	)
	if err != nil {
		return err
	}

	// Check and add categories if needed
	var categoryCount int64
	db.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{
				Name: "Electronics",
				Slug: "electronics",
			},
			{
				Name: "Clothing",
				Slug: "clothing",
			},
			{
				Name: "Books",
				Slug: "books",
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
				Name: "Best Sellers",
				Slug: "best-sellers",
			},
			{
				Name: "New Arrivals",
				Slug: "new-arrivals",
			},
			{
				Name: "Special Offers",
				Slug: "special-offers",
			},
		}

		for _, collection := range collections {
			if err := db.Create(&collection).Error; err != nil {
				log.Printf("Error creating collection %s: %v", collection.Name, err)
			}
		}
	}

	// Check and add pages if needed
	var pageCount int64
	db.Model(&models.Page{}).Count(&pageCount)
	if pageCount == 0 {
		pages := []models.Page{
			{
				Title:          "About Us",
				Slug:           "about-us",
				Content:        "Welcome to our store! We are dedicated to providing the best products and services to our customers.",
				SeoTitle:       "About Us - Our Story and Mission",
				SeoDescription: "Learn more about our company, our mission, and our commitment to customer satisfaction.",
			},
			{
				Title:          "Contact",
				Slug:           "contact",
				Content:        "Get in touch with us! Email: support@example.com, Phone: (555) 123-4567",
				SeoTitle:       "Contact Us - Get in Touch",
				SeoDescription: "Contact us for any questions or support. We're here to help!",
			},
			{
				Title:          "Privacy Policy",
				Slug:           "privacy-policy",
				Content:        "This privacy policy describes how we collect, use, and protect your personal information.",
				SeoTitle:       "Privacy Policy - How We Protect Your Data",
				SeoDescription: "Learn about our privacy practices and how we protect your personal information.",
			},
		}

		for _, page := range pages {
			if err := db.Create(&page).Error; err != nil {
				log.Printf("Error creating page %s: %v", page.Title, err)
			}
		}
	}

	// Check and add menus if needed
	var menuCount int64
	db.Model(&models.Menu{}).Count(&menuCount)
	if menuCount == 0 {
		// Get existing categories, collections and pages
		var electronicsCategory models.Category
		db.Where("slug = ?", "electronics").First(&electronicsCategory)

		var clothingCategory models.Category
		db.Where("slug = ?", "clothing").First(&clothingCategory)

		var bestSellersCollection models.Collection
		db.Where("slug = ?", "best-sellers").First(&bestSellersCollection)

		var aboutUsPage models.Page
		db.Where("slug = ?", "about-us").First(&aboutUsPage)

		menus := []models.Menu{
			{
				Name:    "Main Menu",
				Slug:    "main-menu",
				Channel: "default",
				Items: []models.MenuItem{
					{
						Name:  "Home",
						Level: 1,
						URL:   "/",
					},
					{
						Name:       "Electronics",
						Level:      1,
						URL:        "/electronics",
						CategoryID: &electronicsCategory.ID,
					},
					{
						Name:       "Clothing",
						Level:      1,
						URL:        "/clothing",
						CategoryID: &clothingCategory.ID,
					},
					{
						Name:         "Best Sellers",
						Level:        1,
						URL:          "/best-sellers",
						CollectionID: &bestSellersCollection.ID,
					},
					{
						Name:   "About",
						Level:  1,
						URL:    "/about-us",
						PageID: &aboutUsPage.ID,
					},
					{
						Name:  "Contact",
						Level: 1,
						URL:   "/contact",
					},
				},
			},
			{
				Name:    "Footer Menu",
				Slug:    "footer-menu",
				Channel: "default",
				Items: []models.MenuItem{
					{
						Name:  "Privacy Policy",
						Level: 1,
						URL:   "/privacy-policy",
					},
					{
						Name:  "Terms of Service",
						Level: 1,
						URL:   "/terms",
					},
					{
						Name:  "Shipping Policy",
						Level: 1,
						URL:   "/shipping",
					},
				},
			},
		}

		for _, menu := range menus {
			if err := db.Create(&menu).Error; err != nil {
				log.Printf("Error creating menu %s: %v", menu.Name, err)
			}
		}
	}

	return nil
}
