package main

import (
	"log"
	"os"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/graphql"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Initialize PostgreSQL
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	pageRepo := repository.NewPageRepository(db, rdb)
	menuRepo := repository.NewMenuRepository(db, rdb)

	// Initialize services
	pageService := service.NewPageService(pageRepo)
	menuService := service.NewMenuService(menuRepo)

	// Initialize GraphQL resolver
	resolver := graphql.NewResolver(pageService, menuService)

	// Initialize Gin router
	r := gin.Default()

	// Add GraphQL handler
	r.POST("/query", graphql.Handler(resolver))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 