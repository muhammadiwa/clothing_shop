package handler

import (
	"net/http"
	"strconv"

	"fashion-shop/internal/domain/usecase"
	"fashion-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase    usecase.UserUseCase
	addressUseCase usecase.AddressUseCase
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userUseCase usecase.UserUseCase, addressUseCase usecase.AddressUseCase) *UserHandler {
	return &UserHandler{
		userUseCase:    userUseCase,
		addressUseCase: addressUseCase,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	user, err := h.userUseCase.Register(c, request.Email, request.Password, request.Name, request.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	accessToken, refreshToken, err := h.userUseCase.Login(c, request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	accessToken, err := h.userUseCase.RefreshToken(c, request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
	})
}

// GetProfile handles getting the user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.userUseCase.GetProfile(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile handles updating the user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	var request struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	user, err := h.userUseCase.UpdateProfile(c, userID, request.Name, request.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}

// ChangePassword handles changing the user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")
	var request struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	err := h.userUseCase.ChangePassword(c, userID, request.OldPassword, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// RequestPasswordReset handles requesting a password reset
func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	err := h.userUseCase.RequestPasswordReset(c, request.Email)
	if err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you will receive a password reset link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you will receive a password reset link"})
}

// ResetPassword handles resetting a password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	err := h.userUseCase.ResetPassword(c, request.Token, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// Logout handles user logout
func (h *UserHandler) Logout(c *gin.Context) {
	userID := c.GetUint("userID")
	err := h.userUseCase.Logout(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetUsers handles getting all users (admin only)
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, count, err := h.userUseCase.GetUsers(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"meta": gin.H{
			"total": count,
			"page":  page,
			"limit": limit,
			"pages": (count + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetUserByID handles getting a user by ID (admin only)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userUseCase.GetUserByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ToggleUserActive handles toggling a user's active status (admin only)
func (h *UserHandler) ToggleUserActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		IsActive bool `json:"is_active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	err = h.userUseCase.ToggleUserActive(c, uint(id), request.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// ResetUserPassword handles resetting a user's password (admin only)
func (h *UserHandler) ResetUserPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetValidationErrorMessage(err)})
		return
	}

	err = h.userUseCase.ResetUserPassword(c, uint(id), request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User password reset successfully"})
}
