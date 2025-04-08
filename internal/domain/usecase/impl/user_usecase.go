package impl

import (
	"context"
	"errors"
	"time"

	"fashion-shop/internal/domain/entity"
	"fashion-shop/internal/domain/repository"
	"fashion-shop/internal/domain/usecase"
	"fashion-shop/internal/infrastructure/auth"
	"fashion-shop/internal/utils"
)

type userUseCase struct {
	userRepo     repository.UserRepository
	jwtService   auth.JWTService
	emailService auth.EmailService
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(userRepo repository.UserRepository, jwtService auth.JWTService, emailService auth.EmailService) usecase.UserUseCase {
	return &userUseCase{
		userRepo:     userRepo,
		jwtService:   jwtService,
		emailService: emailService,
	}
}

// Register registers a new user
func (uc *userUseCase) Register(ctx context.Context, email, password, name, phone string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		Email:     email,
		Password:  hashedPassword,
		Name:      name,
		Phone:     phone,
		Role:      entity.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns JWT tokens
func (uc *userUseCase) Login(ctx context.Context, email, password string) (string, string, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return "", "", errors.New("account is inactive")
	}

	// Verify password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return "", "", err
	}

	// Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken refreshes an access token using a refresh token
func (uc *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", err
	}

	// Check if user is active
	if !user.IsActive {
		return "", errors.New("account is inactive")
	}

	// Generate new access token
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetProfile gets a user's profile
func (uc *userUseCase) GetProfile(ctx context.Context, userID uint) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

// UpdateProfile updates a user's profile
func (uc *userUseCase) UpdateProfile(ctx context.Context, userID uint, name, phone string) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Phone = phone
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

// ChangePassword changes a user's password
func (uc *userUseCase) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return uc.userRepo.ChangePassword(ctx, userID, hashedPassword)
}

// RequestPasswordReset sends a password reset email
func (uc *userUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	// Generate reset token
	resetToken, err := uc.jwtService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	// Send email
	subject := "Password Reset Request"
	body := "Click the link below to reset your password:\n\n"
	body += "https://your-website.com/reset-password?token=" + resetToken + "\n\n"
	body += "This link will expire in 1 hour."

	return uc.emailService.SendEmail(user.Email, subject, body)
}

// ResetPassword resets a user's password using a token
func (uc *userUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate token
	claims, err := uc.jwtService.ValidatePasswordResetToken(token)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return uc.userRepo.ChangePassword(ctx, claims.UserID, hashedPassword)
}

// Logout logs out a user
func (uc *userUseCase) Logout(ctx context.Context, userID uint) error {
	// In a stateless JWT system, there's no server-side logout
	// The client should discard the tokens
	return nil
}

// GetUsers gets a list of users (admin function)
func (uc *userUseCase) GetUsers(ctx context.Context, page, limit int) ([]*entity.User, int64, error) {
	offset := (page - 1) * limit
	users, count, err := uc.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Don't return passwords
	for _, user := range users {
		user.Password = ""
	}

	return users, count, nil
}

// GetUserByID gets a user by ID (admin function)
func (uc *userUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

// ToggleUserActive toggles a user's active status (admin function)
func (uc *userUseCase) ToggleUserActive(ctx context.Context, id uint, isActive bool) error {
	return uc.userRepo.ToggleActive(ctx, id, isActive)
}

// ResetUserPassword resets a user's password (admin function)
func (uc *userUseCase) ResetUserPassword(ctx context.Context, id uint, newPassword string) error {
	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return uc.userRepo.ChangePassword(ctx, id, hashedPassword)
}
