package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
	RedisConfig  RedisConfig
	RabbitConfig RabbitConfig
}

type ServerConfig struct {
	HTTPPort int
	GRPCPort int
	Env      string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host string
	Port int
}

type RabbitConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func LoadConfig() (*Config, error) {
	httpPort, err := strconv.Atoi(getEnvOrDefault("HTTP_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	grpcPort, err := strconv.Atoi(getEnvOrDefault("GRPC_PORT", "50051"))
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	redisPort, err := strconv.Atoi(getEnvOrDefault("REDIS_PORT", "6379"))
	if err != nil {
		return nil, err
	}

	rabbitPort, err := strconv.Atoi(getEnvOrDefault("RABBITMQ_PORT", "5672"))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerConfig: ServerConfig{
			HTTPPort: httpPort,
			GRPCPort: grpcPort,
			Env:      getEnvOrDefault("ENV", "development"),
		},
		DBConfig: DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", "Ecommerce"),
		},
		RedisConfig: RedisConfig{
			Host: getEnvOrDefault("REDIS_HOST", "127.0.0.1"),
			Port: redisPort,
		},
		RabbitConfig: RabbitConfig{
			Host:     getEnvOrDefault("RABBITMQ_HOST", "localhost"),
			Port:     rabbitPort,
			User:     getEnvOrDefault("RABBITMQ_USER", "guest"),
			Password: getEnvOrDefault("RABBITMQ_PASSWORD", "guest"),
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
