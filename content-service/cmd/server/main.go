package main

import (
	"fmt"
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/database/migration"
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

	// Load config
	cfg := config.LoadConfig()
	log.Printf("Loaded config: Database URL: %s, Redis Addr: %s, Port: %s", 
		cfg.DatabaseURL, cfg.RedisAddr, cfg.Port)

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// Initialize PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := migration.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
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
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 