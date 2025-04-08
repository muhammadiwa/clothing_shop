package controllers

import (
	"net/http"

	"github.com/fashion-shop/models"
	"github.com/fashion-shop/services"
	"github.com/gin-gonic/gin"
)

// UserController handles user requests
type UserController struct {
	userService    *services.UserService
	addressService *services.AddressService
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{
		userService:    services.NewUserService(),
		addressService: services.NewAddressService(),
	}
}

// GetProfile handles get user profile
// @Summary Get user profile
// @Description Get current user profile
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user
	user, err := c.userService.GetUserByID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"phone":      user.Phone,
			"role":       user.Role,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
		},
	})
}

// UpdateProfile handles update user profile
// @Summary Update user profile
// @Description Update current user profile
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param profile body models.UpdateProfileRequest true "Profile data"
// @Success 200 {object} models.UpdateProfileResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req models.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user profile
	err := c.userService.UpdateProfile(userID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

// GetAddresses handles get user addresses
// @Summary Get user addresses
// @Description Get current user addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Address
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /addresses [get]
func (c *UserController) GetAddresses(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get addresses
	addresses, err := c.addressService.GetAddressesByUserID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
	})
}

// AddAddress handles add user address
// @Summary Add user address
// @Description Add a new address for current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param address body models.AddressRequest true "Address data"
// @Success 201 {object} models.AddressResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /addresses [post]
func (c *UserController) AddAddress(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req models.AddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add address
	address, err := c.addressService.AddAddress(userID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"address": address,
	})
}

// UpdateAddress handles update user address
// @Summary Update user address
// @Description Update an existing address for current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Address ID"
// @Param address body models.AddressRequest true "Address data"
// @Success 200 {object} models.AddressResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /addresses/{id} [put]
func (c *UserController) UpdateAddress(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get address ID from path
	addressID := ctx.Param("id")

	// Parse request body
	var req models.AddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update address
	address, err := c.addressService.UpdateAddress(userID.(string), addressID, req)
	if err != nil {
		if err.Error() == "address not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"address": address,
	})
}

// DeleteAddress handles delete user address
// @Summary Delete user address
// @Description Delete an existing address for current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Address ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /addresses/{id} [delete]
func (c *UserController) DeleteAddress(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get address ID from path
	addressID := ctx.Param("id")

	// Delete address
	err := c.addressService.DeleteAddress(userID.(string), addressID)
	if err != nil {
		if err.Error() == "address not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Address deleted successfully",
	})
}

// SetDefaultAddress handles setting default user address
// @Summary Set default user address
// @Description Set an address as default for current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Address ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /addresses/{id}/default [put]
func (c *UserController) SetDefaultAddress(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get address ID from path
	addressID := ctx.Param("id")

	// Set default address
	err := c.addressService.SetDefaultAddress(userID.(string), addressID)
	if err != nil {
		if err.Error() == "address not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Default address set successfully",
	})
}
