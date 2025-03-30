package configs

import (
	"os"
	"strconv"
	"time"
)

// EnvConfig holds environment configuration
type EnvConfig struct {
	// Database
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	// Redis
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	RedisDefaultTTL time.Duration

	// Service ports
	OrderPort string
	GrpcPort  string

	// Product Service
	ProductServiceHost string
	ProductServicePort string
}

// GetEnvConfig loads environment configuration
func GetEnvConfig() EnvConfig {
	redisDB, _ := strconv.Atoi(getEnvWithDefault("REDIS_DB", "0"))
	redisTTL, _ := strconv.Atoi(getEnvWithDefault("REDIS_DEFAULT_TTL", "3600"))

	return EnvConfig{
		// Database
		DbHost:     getEnvWithDefault("DB_HOST", "localhost"),
		DbPort:     getEnvWithDefault("DB_PORT", "5432"),
		DbUser:     getEnvWithDefault("DB_USER", "postgres"),
		DbPassword: getEnvWithDefault("DB_PASSWORD", "postgres"),
		DbName:     getEnvWithDefault("DB_NAME", "order_service"),

		// Redis
		RedisAddr:       getEnvWithDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnvWithDefault("REDIS_PASSWORD", ""),
		RedisDB:         redisDB,
		RedisDefaultTTL: time.Duration(redisTTL) * time.Second,

		// Service ports
		OrderPort: getEnvWithDefault("ORDER_PORT", "8080"),
		GrpcPort:  getEnvWithDefault("GRPC_PORT", "50053"),

		// Product Service
		ProductServiceHost: getEnvWithDefault("PRODUCT_SERVICE_HOST", "localhost"),
		ProductServicePort: getEnvWithDefault("PRODUCT_SERVICE_PORT", "50051"),
	}
}

// getEnvWithDefault gets an environment variable or returns the default
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
