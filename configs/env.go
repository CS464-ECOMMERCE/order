package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
)

// EnvConfig holds environment configuration
type EnvConfig struct {
	// Database
	PostgresConnStringMaster string
	PostgresConnStringSlave  string
	PostgresMaxIdleConns     int
	PostgresMaxOpenConns     int

	// Service ports
	GrpcPort string

	// Product Service
	ProductServiceAddr string
	CartServiceAddr    string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

var envConfig EnvConfig

// InitEnv initializes environment variables
func InitEnv() {
	// Load .env file if it exists
	err := godotenv.Load("/app/secrets/.env")
	if err != nil {
		fmt.Println("Warning: Error loading env file:", err)
	}

	// Populate envConfig
	envConfig = EnvConfig{
		// Database
		PostgresConnStringMaster: getEnv("POSTGRESQL_CONN_STRING_MASTER", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"),
		PostgresConnStringSlave:  getEnv("POSTGRESQL_CONN_STRING_SLAVE", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"),
		PostgresMaxIdleConns:     getEnvAsInt("POSTGRESQL_MAX_IDLE_CONNS", 5),
		PostgresMaxOpenConns:     getEnvAsInt("POSTGRESQL_MAX_OPEN_CONNS", 10),

		// Service ports
		GrpcPort: getEnv("ORDER_SERVICE_GRPC_PORT", "50052"),

		// Product Service
		ProductServiceAddr: getEnv("PRODUCT_SERVICE_ADDR", "product.default.svc.cluster.local:50050"),
		CartServiceAddr:    getEnv("CART_SERVICE_ADDR", "cart.default.svc.cluster.local:50050"),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "redis_password"),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
	}

	stripe.Key = getEnv("STRIPE_SECRET_KEY", "some-secret-key")

	fmt.Println("Order service environment variables initialized")
}

// GetEnvConfig returns the current environment configuration
func GetEnvConfig() EnvConfig {
	return envConfig
}

// Helper function to get environment variable with fallback
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get environment variable as int
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// Helper function to get environment variable as duration
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
