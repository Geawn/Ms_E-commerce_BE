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
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/config"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/database"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/database/migration"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/graphql"
	grpcserver "github.com/Geawn/Ms_E-commerce_BE/user-service/internal/grpc"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/service"
	pb "github.com/Geawn/Ms_E-commerce_BE/user-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		Addr: fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	grpcServer := grpcserver.NewServer(userService)

	// Initialize gRPC server
	s := grpc.NewServer(
		grpc.MaxConcurrentStreams(100),
		grpc.MaxRecvMsgSize(1024*1024), // 1MB
		grpc.MaxSendMsgSize(1024*1024), // 1MB
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
		)),
	)

	// Đăng ký service
	pb.RegisterUserServiceServer(s, grpcServer)

	// Enable reflection cho development tools
	reflection.Register(s)

	// Enable Prometheus metrics
	grpc_prometheus.Register(s)

	// Initialize Gin router
	router := gin.Default()

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Setup GraphQL server
	resolver := graphql.NewResolver(userService)
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: resolver,
	}))

	// Add GraphQL endpoints
	router.POST("/query", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	// Add GraphQL playground in development
	if cfg.ServerConfig.Env == "development" {
		playgroundHandler := playground.Handler("GraphQL", "/query")
		router.GET("/", func(c *gin.Context) {
			playgroundHandler.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Lắng nghe trên port
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerConfig.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerConfig.HTTPPort),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down servers...")
		s.GracefulStop()
		httpServer.Shutdown(context.Background())
	}()

	// Start servers
	go func() {
		log.Printf("Starting gRPC server on :%d", cfg.ServerConfig.GRPCPort)
		if err := s.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	log.Printf("Starting HTTP server on :%d in %s environment", cfg.ServerConfig.HTTPPort, cfg.ServerConfig.Env)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
