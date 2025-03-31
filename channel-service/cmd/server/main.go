package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/database"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/graphql"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/repository"
)

const defaultPort = "8080"

func main() {
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

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
