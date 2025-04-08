package controllers

import (
	"net/http"

	"github.com/fashion-shop/config"
	"github.com/fashion-shop/models"
	"github.com/fashion-shop/services"
	"github.com/fashion-shop/utils"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication requests
type AuthController struct {
	authService *services.AuthService
	userService *services.UserService
}

// NewAuthController creates a new auth controller
func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(cfg),
		userService: services.NewUserService(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password, name, address, and phone
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
	}

	if err := c.authService.RegisterUser(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log activity
	activity := &models.UserActivity{
		UserID:    user.ID,
		Activity:  "User registered",
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	c.userService.LogUserActivity(activity)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
	})
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User login data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Login user
	accessToken, refreshToken, err := c.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := c.userService.GetUserByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Log activity
	activity := &models.UserActivity{
		UserID:    user.ID,
		Activity:  "User logged in",
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	c.userService.LogUserActivity(activity)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.RefreshTokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Refresh token
	accessToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// ForgotPassword handles password reset request
// @Summary Request password reset
// @Description Request password reset via email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body models.ForgotPasswordRequest true "User email"
// @Success 200 {object} models.ForgotPasswordResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Generate reset token
	resetToken, err := c.authService.ForgotPassword(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real application, send email with reset link
	// For this example, we'll just return the token
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset link sent to email",
		"token":   resetToken, // In production, don't return this
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body models.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} models.ResetPasswordResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Reset password
	if err := c.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// Logout handles user logout
// @Summary Logout a user
// @Description Logout a user by invalidating refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.LogoutResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Logout user
	if err := c.authService.LogoutUser(userID.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity
	activity := &models.UserActivity{
		UserID:    userID.(string),
		Activity:  "User logged out",
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	c.userService.LogUserActivity(activity)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param password body models.ChangePasswordRequest true "Change password data"
// @Success 200 {object} models.ChangePasswordResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/password [put]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Change password
	if err := c.authService.ChangePassword(userID.(string), req.CurrentPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log activity
	activity := &models.UserActivity{
		UserID:    userID.(string),
		Activity:  "User changed password",
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	c.userService.LogUserActivity(activity)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
