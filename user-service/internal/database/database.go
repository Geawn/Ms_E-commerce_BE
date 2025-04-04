package database

import (
	"fmt"
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto Migrate the schema in correct order
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error migrating users table:", err)
	}

	err = db.AutoMigrate(&models.Profile{})
	if err != nil {
		log.Fatal("Error migrating profiles table:", err)
	}

	err = db.AutoMigrate(&models.Address{})
	if err != nil {
		log.Fatal("Error migrating addresses table:", err)
	}

	err = db.AutoMigrate(&models.Avatar{})
	if err != nil {
		log.Fatal("Error migrating avatars table:", err)
	}

	log.Println("Successfully connected to database and migrated tables")
	return db, nil
}
