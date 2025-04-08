package services

import (
	"errors"

	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/repository"
	"clothing-shop-api/pkg/utils"
)

type AuthService struct {
	authRepo  repository.AuthRepository
	jwtSecret string
}

func NewAuthService(authRepo repository.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(email, password, name, phone string) error {
	// Check if user already exists
	existing, _ := s.authRepo.FindByEmail(email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Phone:    phone,
	}

	return s.authRepo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	user, err := s.authRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	// Update last login
	s.authRepo.UpdateLastLogin(user.ID)

	return accessToken, refreshToken, nil
}
