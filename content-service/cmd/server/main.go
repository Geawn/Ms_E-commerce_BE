package main

import (
	"context"
	"log"
	"os"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/database/migration"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/graph"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/service"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=content_service port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := migration.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize repositories
	menuRepo := repository.NewMenuRepository(db)
	pageRepo := repository.NewPageRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)

	// Initialize services
	menuService := service.NewMenuService(menuRepo, pageRepo, categoryRepo, collectionRepo, redisClient)
	pageService := service.NewPageService(pageRepo, redisClient)

	// Initialize resolver
	resolver := graph.NewResolver(menuService, pageService)

	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Setup Gin router
	r := gin.Default()

	// GraphQL playground
	r.GET("/", gin.WrapF(playground.Handler("GraphQL playground", "/query")))

	// GraphQL endpoint
	r.POST("/query", gin.WrapF(srv.ServeHTTP))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(r.Run(":" + port))
} 