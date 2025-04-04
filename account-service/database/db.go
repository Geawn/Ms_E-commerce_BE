package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Geawn/Ms_E-commerce_BE/account-service/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default values")
	}

	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "77780409")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "Ecommerce")
	sslMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&models.UserToken{})
	if err != nil {
		log.Fatal("Error migrating UserToken table:", err)
	}

	err = db.AutoMigrate(&models.AccountUser{})
	if err != nil {
		log.Fatal("Error migrating AccountUser table:", err)
	}

	DB = db
	log.Println("Successfully connected to database")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
