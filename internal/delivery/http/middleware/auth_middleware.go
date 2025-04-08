package middleware

import (
	"net/http"
	"strings"

	"fashion-shop/internal/domain/entity"
	"fashion-shop/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware for authentication
type AuthMiddleware struct {
	jwtService auth.JWTService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(jwtService auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth requires authentication for a route
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID and role in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// RequireAdmin requires admin role for a route
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First require authentication
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Check if user has admin role
		role, exists := c.Get("userRole")
		if !exists || role.(entity.Role) != entity.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
