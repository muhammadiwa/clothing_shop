package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/repository"
	"clothing-shop-api/pkg/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) repository.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User) error {
	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrEmailAlreadyExists
	}

	// Generate verification token
	verificationToken, err := generateRandomToken(32)
	if err != nil {
		return err
	}
	user.VerificationToken = verificationToken
	user.IsVerified = false

	// Save user
	return s.repo.Create(user)
}

func (s *userService) VerifyEmail(token string) error {
	user, err := s.repo.FindByVerificationToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidToken
	}

	// Update user
	user.IsVerified = true
	user.VerificationToken = ""
	return s.repo.Update(user)
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := utils.VerifyPassword(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) GeneratePasswordResetToken(email string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// Generate reset token
	resetToken, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}

	// Set expiry to 24 hours from now
	expiry := time.Now().Add(24 * time.Hour)
	user.ResetToken = resetToken
	user.ResetTokenExpiry = &expiry

	// Update user
	if err := s.repo.Update(user); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *userService) ResetPassword(token, newPassword string) error {
	user, err := s.repo.FindByResetToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidToken
	}

	// Update password and clear reset token
	if err := s.repo.UpdatePassword(user.ID, newPassword); err != nil {
		return err
	}

	// Clear reset token
	user.ResetToken = ""
	user.ResetTokenExpiry = nil
	return s.repo.Update(user)
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUser(user *models.User) error {
	existingUser, err := s.repo.FindByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	existingUser, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(id)
}

// Helper function to generate random tokens
func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
