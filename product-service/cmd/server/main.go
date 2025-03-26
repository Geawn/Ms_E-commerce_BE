package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/database"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/database/migration"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/event"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/graphql"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/service"
	pb "github.com/Geawn/Ms_E-commerce_BE/product-service/proto"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get the underlying sql.DB and defer its closing
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := migration.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize RabbitMQ
	rabbitmq, err := event.NewRabbitMQService(
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
	)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo)
	grpcService := service.NewGRPCProductService(productService)

	// Initialize GraphQL resolver
	resolver := graphql.NewResolver(productService)

	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	// Initialize Gin router
	router := gin.Default()

	// Add GraphQL playground in development
	if os.Getenv("ENV") != "production" {
		router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	}

	// Add GraphQL endpoint
	router.POST("/query", gin.WrapH(srv))

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, grpcService)
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down servers...")
		grpcServer.GracefulStop()
		httpServer.Shutdown(context.Background())
	}()

	// Start servers
	go func() {
		log.Printf("Starting gRPC server on :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	log.Printf("Starting HTTP server on :%s", cfg.HTTPPort)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
