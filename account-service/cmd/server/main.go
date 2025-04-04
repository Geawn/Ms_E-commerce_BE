package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Geawn/Ms_E-commerce_BE/account-service/database"
	"github.com/Geawn/Ms_E-commerce_BE/account-service/handlers"
	"github.com/Geawn/Ms_E-commerce_BE/account-service/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Set working directory to project root
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get current file path")
	}
	projectRoot := filepath.Join(filepath.Dir(filename), "../..")
	if err := os.Chdir(projectRoot); err != nil {
		log.Fatal("Failed to change working directory:", err)
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	database.InitDB()

	// Set Gin mode
	ginMode := os.Getenv("ENV")
	if ginMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Serve static files
	r.Static("/static", "static")
	r.StaticFile("/", "static/index.html")

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/change-password", handlers.ChangePassword)
		protected.POST("/logout", handlers.Logout)
	}

	// Get port from environment variable or use default
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	log.Printf("Server starting on http://%s:%s\n", host, port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
