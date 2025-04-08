package middleware

import (
	"net/http"
	"time"

	"github.com/fashion-shop/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// SetupGlobalMiddleware sets up all global middleware for the application
func SetupGlobalMiddleware(router *gin.Engine, cfg *config.Config) {
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting middleware
	rate := limiter.Rate{
		Period: cfg.RateLimitDuration,
		Limit:  int64(cfg.RateLimitRequests),
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)
	router.Use(ginlimiter.NewMiddleware(rateLimiter))

	// Recovery middleware
	router.Use(gin.Recovery())

	// Request logging middleware
	router.Use(RequestLogger())

	// Security headers middleware
	router.Use(SecurityHeaders())
}

// RequestLogger logs request details
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		// Log using Gin's logger
		gin.DefaultWriter.Write([]byte(
			"[GIN] " + time.Now().Format("2006/01/02 - 15:04:05") +
				" | " + statusCode + " | " +
				latency.String() + " | " +
				clientIP + " | " +
				method + " | " +
				path + "\n"))
	}
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// AuthMiddleware verifies JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Validate token
		// This is a placeholder - actual implementation will be in auth service
		tokenString := authHeader[7:] // Remove "Bearer " prefix
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// AdminMiddleware ensures the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First apply auth middleware
		AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// Check if user has admin role
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateToken validates a JWT token and returns the claims
// This is a placeholder - actual implementation will be in auth service
func ValidateToken(tokenString string) (*Claims, error) {
	// Placeholder - will be implemented in auth service
	return nil, nil
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
