package database

import (
	"fmt"
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/database/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Kết nối database và chạy migration
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Chạy migration từ thư mục database/migration
	log.Println("Running database migrations...")
	if err := migration.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Println("Database migration completed successfully.")

	return db, nil
}
