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

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(user *models.User) error {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
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
	return s.userRepo.Create(user)
}

func (s *UserService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidToken
	}

	// Update user
	user.IsVerified = true
	user.VerificationToken = ""
	return s.userRepo.Update(user)
}

func (s *UserService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
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

func (s *UserService) GeneratePasswordResetToken(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
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
	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidToken
	}

	// Update password and clear reset token
	if err := s.userRepo.UpdatePassword(user.ID, newPassword); err != nil {
		return err
	}

	// Clear reset token
	user.ResetToken = ""
	user.ResetTokenExpiry = nil
	return s.userRepo.Update(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	existingUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	existingUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(id)
}

// Helper function to generate random tokens
func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
