package services

import (
	"context"
	"fmt"
	"order/configs"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient wraps the Redis client with additional methods for order operations
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	config configs.EnvConfig
}

var redisClient *RedisClient

// GetRedisClient returns a singleton instance of the Redis client
func GetRedisClient() *RedisClient {
	if redisClient == nil {
		config := configs.GetEnvConfig()
		ctx := context.Background()

		client := redis.NewClient(&redis.Options{
			Addr:         config.RedisAddr,
			Password:     config.RedisPassword,
			DB:           config.RedisDB,
			PoolSize:     10,
			MinIdleConns: 5,
		})

		// Test connection
		if _, err := client.Ping(ctx).Result(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}

		fmt.Println("Successfully connected to Redis")

		redisClient = &RedisClient{
			client: client,
			ctx:    ctx,
			config: config,
		}
	}

	return redisClient
}

// Close closes the Redis client connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// getCartKey returns the Redis key for a user's cart
func (r *RedisClient) getCartKey(session_id string) string {
	return fmt.Sprintf("cart:%s", session_id)
}

// GetCart retrieves a user's cart from Redis
func (r *RedisClient) GetCart(session_id string) ([]byte, error) {
	key := r.getCartKey(session_id)
	return r.client.Get(r.ctx, key).Bytes()
}

// DeleteCart deletes a user's cart
func (r *RedisClient) DeleteCart(session_id string) error {
	key := r.getCartKey(session_id)
	return r.client.Del(r.ctx, key).Err()
}

// ExecuteWithLock executes a function with a distributed lock
func (r *RedisClient) ExecuteWithLock(key string, ttl time.Duration, fn func() error) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	// Set lock with NX option (only set if key doesn't exist)
	ok, err := r.client.SetNX(r.ctx, lockKey, 1, ttl).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !ok {
		return fmt.Errorf("failed to acquire lock, already locked")
	}

	defer r.client.Del(r.ctx, lockKey)

	return fn()
}
