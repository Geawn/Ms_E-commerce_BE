package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/database"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/graphql"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/service"
	"github.com/Geawn/Ms_E-commerce_BE/order-service/proto"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		ctx := context.WithValue(c.Request.Context(), "token", token)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

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

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Connect to user service
	userConn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	// Connect to product service
	productConn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	userService := service.NewUserService(proto.NewUserServiceClient(userConn))
	productService := service.NewProductService(productConn)
	orderService := service.NewOrderService(orderRepo, userService, productService, db, redisClient)

	// Create resolver
	resolver := &graphql.Resolver{
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	}

	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Add authentication middleware
	router.Use(authMiddleware())

	// Add GraphQL playground in development
	if os.Getenv("ENV") != "production" {
		router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	}

	// Add GraphQL endpoint
	router.POST("/query", gin.WrapH(srv))

	// Create HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8081" // Default port if not set
	}
	srvHTTP := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server started at http://localhost:%s", httpPort)
		if err := srvHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := srvHTTP.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
