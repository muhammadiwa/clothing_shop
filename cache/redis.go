package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fashion-shop/config"
	"github.com/go-redis/redis/v8"
)

// RedisClient is the Redis client
var RedisClient *redis.Client

// Initialize initializes the Redis client
func Initialize(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Ping Redis to check connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

// Set sets a key-value pair in Redis
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()

	// Convert value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set key-value pair
	return RedisClient.Set(ctx, key, jsonValue, expiration).Err()
}

// Get gets a value from Redis
func Get(key string, dest interface{}) error {
	ctx := context.Background()

	// Get value
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	// Unmarshal JSON
	return json.Unmarshal([]byte(val), dest)
}

// Delete deletes a key from Redis
func Delete(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func Exists(key string) (bool, error) {
	ctx := context.Background()
	val, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// SetHash sets a hash in Redis
func SetHash(key string, field string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()

	// Convert value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set hash field
	err = RedisClient.HSet(ctx, key, field, jsonValue).Err()
	if err != nil {
		return err
	}

	// Set expiration if provided
	if expiration > 0 {
		return RedisClient.Expire(ctx, key, expiration).Err()
	}

	return nil
}

// GetHash gets a hash field from Redis
func GetHash(key string, field string, dest interface{}) error {
	ctx := context.Background()

	// Get hash field
	val, err := RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		return err
	}

	// Unmarshal JSON
	return json.Unmarshal([]byte(val), dest)
}

// DeleteHash deletes a hash field from Redis
func DeleteHash(key string, field string) error {
	ctx := context.Background()
	return RedisClient.HDel(ctx, key, field).Err()
}

// GetAllHash gets all hash fields from Redis
func GetAllHash(key string) (map[string]string, error) {
	ctx := context.Background()
	return RedisClient.HGetAll(ctx, key).Result()
}
