package auth

import (
	"errors"
	"time"

	"fashion-shop/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in a JWT
type JWTClaims struct {
	UserID uint        `json:"user_id"`
	Role   entity.Role `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// JWTService defines the interface for JWT operations
type JWTService interface {
	GenerateAccessToken(userID uint, role entity.Role) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	GeneratePasswordResetToken(userID uint) (string, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshToken(tokenString string) (*JWTClaims, error)
	ValidatePasswordResetToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	accessSecret  string
	refreshSecret string
	resetSecret   string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	resetExpiry   time.Duration
}

// NewJWTService creates a new JWTService instance
func NewJWTService(accessSecret, refreshSecret, resetSecret string, accessExpiry, refreshExpiry, resetExp refreshSecret, resetSecret string, accessExpiry, refreshExpiry, resetExpiry time.Duration) JWTService {
	return &jwtService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		resetSecret:   resetSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		resetExpiry:   resetExpiry,
	}
}

// GenerateAccessToken generates a new access token
func (s *jwtService) GenerateAccessToken(userID uint, role entity.Role) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

// GenerateRefreshToken generates a new refresh token
func (s *jwtService) GenerateRefreshToken(userID uint) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

// GeneratePasswordResetToken generates a new password reset token
func (s *jwtService) GeneratePasswordResetToken(userID uint) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.resetExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.resetSecret))
}

// ValidateAccessToken validates an access token
func (s *jwtService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates a refresh token
func (s *jwtService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidatePasswordResetToken validates a password reset token
func (s *jwtService) ValidatePasswordResetToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.resetSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
