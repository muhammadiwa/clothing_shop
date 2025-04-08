package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter is a middleware for rate limiting
type RateLimiter struct {
	redisClient *redis.Client
	maxRequests int
	duration    time.Duration
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(redisClient *redis.Client, maxRequests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		maxRequests: maxRequests,
		duration:    duration,
	}
}

// Middleware returns a middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		ctx := context.Background()

		// Get current count
		count, err := rl.redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
			c.Abort()
			return
		}

		// If key doesn't exist, create it
		if err == redis.Nil {
			_, err = rl.redisClient.Set(ctx, key, 1, rl.duration).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
				c.Abort()
				return
			}
			count = 1
		} else {
			// Increment count
			count, err = rl.redisClient.Incr(ctx, key).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
				c.Abort()
				return
			}
		}

		// Get TTL
		ttl, err := rl.redisClient.TTL(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
			c.Abort()
			return
		}

		// Set headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(rl.maxRequests-count))
		c.Header("X-RateLimit-Reset", strconv.Itoa(int(ttl.Seconds())))

		// Check if limit exceeded
		if count > rl.maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
