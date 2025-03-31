package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL    string
	RedisAddr      string
	RedisPassword  string
	Port           string
}

func LoadConfig() *Config {
	// Construct database URL from individual components
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Construct Redis address
	redisAddr := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	return &Config{
		DatabaseURL:   dbURL,
		RedisAddr:     redisAddr,
		RedisPassword: "", // No password set in .env
		Port:          os.Getenv("HTTP_PORT"),
	}
} 