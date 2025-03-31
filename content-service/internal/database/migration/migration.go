package migration

import (
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate các model
	err := db.AutoMigrate(
		&models.Menu{},
		&models.MenuItem{},
		&models.Page{},
	)
	if err != nil {
		return err
	}

	// Check and add default menu if needed
	var menuCount int64
	db.Model(&models.Menu{}).Count(&menuCount)
	if menuCount == 0 {
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
						Name:  "Products",
						Level: 1,
						URL:   "/products",
					},
					{
						Name:  "About Us",
						Level: 1,
						URL:   "/about",
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

	// Check and add default page if needed
	var pageCount int64
	db.Model(&models.Page{}).Count(&pageCount)
	if pageCount == 0 {
		pages := []models.Page{
			{
				Title:          "About Us",
				Slug:           "about",
				Content:        "Welcome to our store! We are dedicated to providing the best products and services.",
				SeoTitle:       "About Us - Our Story and Mission",
				SeoDescription: "Learn more about our company, our mission, and our commitment to customer satisfaction.",
				IsPublished:    true,
			},
			{
				Title:          "Contact Us",
				Slug:           "contact",
				Content:        "Get in touch with us! We're here to help.",
				SeoTitle:       "Contact Us - Get in Touch",
				SeoDescription: "Contact us for any questions, concerns, or support needs.",
				IsPublished:    true,
			},
		}

		for _, page := range pages {
			if err := db.Create(&page).Error; err != nil {
				log.Printf("Error creating page %s: %v", page.Title, err)
			}
		}
	}

	return nil
}
