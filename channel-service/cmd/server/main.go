package main

import (
	"context"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/database"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/graphql"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/repository"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

// Defining the Graphql handler
func graphqlHandler(resolver *graphql.Resolver) gin.HandlerFunc {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// Setting up Gin
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Create repository
	repo, err := repository.NewRedisRepository(redisAddr)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Run migration
	ctx := context.Background()
	if err := database.MigrateChannels(ctx, repo); err != nil {
		log.Printf("Warning: Failed to migrate channels: %v", err)
	}

	// Create resolver
	resolver, err := graphql.NewResolver(redisAddr)
	if err != nil {
		log.Fatalf("Failed to create resolver: %v", err)
	}

	// Setting up routes
	r.GET("/", playgroundHandler())
	r.POST("/query", graphqlHandler(resolver))

	log.Printf("Starting HTTP server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
