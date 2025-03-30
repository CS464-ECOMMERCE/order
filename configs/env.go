package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	PORT            string
	API_LISTEN_HOST string

	POSTGRESQL_CONN_STRING_MASTER string
	POSTGRESQL_CONN_STRING_SLAVE  string
	POSTGRESQL_MAX_IDLE_CONNS     int
	POSTGRESQL_MAX_OPEN_CONNS     int

	// Redis
	REDIS_ADDR        string
	REDIS_PASSWORD    string
	REDIS_DB          int
	REDIS_DEFAULT_TTL time.Duration

	// Service ports
	ORDER_PORT string
	GRPC_PORT  string

	// Product Service
	PRODUCT_SERVICE_ADDR string
)

// EnvConfig holds environment configuration
type EnvConfig struct {
	// Database
	PostgresConnString string

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

var envConfig EnvConfig

// InitEnv initializes environment variables
func InitEnv() {
	// Load .env file if it exists
	envPath := "/app/secrets/.env"
	if os.Getenv("ENV") == "dev" {
		envPath = "./secrets/testing.env"
	}
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Warning: Error loading env file:", err)
	}

	// rest api
	PORT = getEnv("API_PORT", "8082")
	API_LISTEN_HOST = getEnv("API_LISTEN_HOST", "0.0.0.0")

	// postgres
	POSTGRESQL_CONN_STRING_MASTER = getEnv("POSTGRESQL_CONN_STRING_MASTER", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai")
	POSTGRESQL_CONN_STRING_SLAVE = getEnv("POSTGRESQL_CONN_STRING_SLAVE", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai")
	maxOpenConns, err := strconv.Atoi(getEnv("POSTGRESQL_MAX_OPEN_CONNS", "10"))
	if err != nil {
		panic("Invalid value for POSTGRESQL_MAX_OPEN_CONNS")
	}
	POSTGRESQL_MAX_OPEN_CONNS = maxOpenConns
	maxIdleConns, err := strconv.Atoi(getEnv("POSTGRESQL_MAX_IDLE_CONNS", "5"))
	if err != nil {
		panic("Invalid value for POSTGRESQL_MAX_IDLE_CONNS")
	}
	POSTGRESQL_MAX_IDLE_CONNS = maxIdleConns

	// Redis
	REDIS_ADDR = getEnv("REDIS_ADDR", "localhost:6379")
	REDIS_PASSWORD = getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	REDIS_DB = redisDB
	redisTTL, _ := strconv.Atoi(getEnv("REDIS_DEFAULT_TTL", "3600"))
	REDIS_DEFAULT_TTL = time.Duration(redisTTL) * time.Second

	// Service ports
	ORDER_PORT = getEnv("ORDER_PORT", "8082")
	GRPC_PORT = getEnv("GRPC_PORT", "50053")

	// Product Service
	PRODUCT_SERVICE_ADDR = getEnv("PRODUCT_SERVICE_ADDR", "localhost:50051")

	// Populate envConfig
	envConfig = EnvConfig{
		// Database
		PostgresConnString: POSTGRESQL_CONN_STRING_MASTER,

		// Redis
		RedisAddr:       REDIS_ADDR,
		RedisPassword:   REDIS_PASSWORD,
		RedisDB:         REDIS_DB,
		RedisDefaultTTL: REDIS_DEFAULT_TTL,

		// Service ports
		OrderPort: ORDER_PORT,
		GrpcPort:  GRPC_PORT,

		// Product Service
		ProductServiceHost: getEnv("PRODUCT_SERVICE_HOST", "localhost"),
		ProductServicePort: getEnv("PRODUCT_SERVICE_PORT", "50051"),
	}

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
