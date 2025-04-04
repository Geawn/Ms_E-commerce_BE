package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Geawn/Ms_E-commerce_BE/account-service/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// Configure GORM
	config := &gorm.Config{
		PrepareStmt: true,                                // Enable prepared statement cache
		Logger:      logger.Default.LogMode(logger.Info), // Enable detailed logging
		// Disable automatic foreign key constraints
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Drop existing tables if they exist
	log.Println("Dropping existing tables...")
	// Drop user_tokens first because it has foreign key to account_users
	if err := db.Migrator().DropTable(&models.UserToken{}); err != nil {
		log.Println("Warning: Error dropping UserToken table:", err)
	}
	if err := db.Migrator().DropTable(&models.AccountUser{}); err != nil {
		log.Println("Warning: Error dropping AccountUser table:", err)
	}

	// Create account_users table first
	log.Println("Creating account_users table...")
	err = db.AutoMigrate(&models.AccountUser{})
	if err != nil {
		log.Fatal("Error migrating AccountUser table:", err)
	}

	// Verify account_users table was created
	if !db.Migrator().HasTable(&models.AccountUser{}) {
		log.Fatal("Failed to create account_users table")
	}
	log.Println("Successfully created account_users table")

	// Then create user_tokens table
	log.Println("Creating user_tokens table...")
	err = db.AutoMigrate(&models.UserToken{})
	if err != nil {
		log.Fatal("Error migrating UserToken table:", err)
	}

	// Verify user_tokens table was created
	if !db.Migrator().HasTable(&models.UserToken{}) {
		log.Fatal("Failed to create user_tokens table")
	}
	log.Println("Successfully created user_tokens table")

	// Drop existing foreign key if exists
	log.Println("Dropping existing foreign key constraint...")
	err = db.Exec(`ALTER TABLE user_tokens DROP CONSTRAINT IF EXISTS fk_user_tokens_account_users`).Error
	if err != nil {
		log.Fatal("Error dropping existing foreign key constraint:", err)
	}
	log.Println("Successfully dropped existing foreign key constraint")

	// Create foreign key constraint with proper ON DELETE and ON UPDATE actions
	log.Println("Creating foreign key constraint...")
	err = db.Exec(`ALTER TABLE user_tokens 
		ADD CONSTRAINT fk_user_tokens_account_users 
		FOREIGN KEY (user_id) 
		REFERENCES account_users(user_id) 
		ON DELETE CASCADE 
		ON UPDATE CASCADE`).Error
	if err != nil {
		log.Fatal("Error creating foreign key constraint:", err)
	}
	log.Println("Successfully created foreign key constraint")

	// Verify foreign key constraint
	log.Println("Verifying foreign key constraint...")
	var count int64
	if err := db.Model(&models.UserToken{}).Count(&count).Error; err != nil {
		log.Fatal("Error verifying foreign key constraint:", err)
	}
	log.Println("Foreign key constraint verified successfully")

	DB = db
	log.Println("Successfully connected to database")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
