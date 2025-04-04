package migration

import (
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate các model
	err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Address{},
		&models.Avatar{},
	)
	if err != nil {
		return err
	}

	// Check and add users if needed
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		users := []models.User{
			{
				Email:     "admin@example.com",
				FirstName: "Admin",
				LastName:  "User",
				Profile: models.Profile{
					PhoneNumber: "+1234567890",
					Addresses: []models.Address{
						{
							Street:     "123 Main St",
							City:       "New York",
							State:      "NY",
							PostalCode: "10001",
							Country:    "USA",
							IsDefault:  true,
						},
					},
					Avatar: &models.Avatar{
						URL: "https://example.com/default-avatar.jpg",
						Alt: "Default Avatar",
					},
				},
			},
			{
				Email:     "user@example.com",
				FirstName: "Regular",
				LastName:  "User",
				Profile: models.Profile{
					PhoneNumber: "+1987654321",
					Addresses: []models.Address{
						{
							Street:     "456 Oak Ave",
							City:       "Los Angeles",
							State:      "CA",
							PostalCode: "90001",
							Country:    "USA",
							IsDefault:  true,
						},
					},
					Avatar: &models.Avatar{
						URL: "https://example.com/default-avatar.jpg",
						Alt: "Default Avatar",
					},
				},
			},
		}

		for _, user := range users {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Error creating user %s: %v", user.Email, err)
			}
		}
	}

	return nil
}
