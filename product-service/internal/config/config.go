package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis config
	RedisHost string
	RedisPort string

	// RabbitMQ config
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string

	// Server config
	HTTPPort  string
	GRPCPort  string
	Env       string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		// Database config
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "product_db"),

		// Redis config
		RedisHost: getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort: getEnvOrDefault("REDIS_PORT", "6379"),

		// RabbitMQ config
		RabbitMQHost:     getEnvOrDefault("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnvOrDefault("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnvOrDefault("RABBITMQ_USER", "guest"),
		RabbitMQPassword: getEnvOrDefault("RABBITMQ_PASSWORD", "guest"),

		// Server config
		HTTPPort: getEnvOrDefault("HTTP_PORT", "8080"),
		GRPCPort: getEnvOrDefault("GRPC_PORT", "50051"),
		Env:      getEnvOrDefault("ENV", "development"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
