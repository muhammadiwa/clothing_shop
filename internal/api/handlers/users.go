package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "clothing-shop-api/internal/domain/services"
)

type UserHandler struct {
    userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

// GetUserProfile retrieves the user profile by ID
func (h *UserHandler) GetUserProfile(c *gin.Context) {
    userID := c.Param("id")
    user, err := h.userService.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// UpdateUserProfile updates the user profile
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
    userID := c.Param("id")
    var userUpdate services.UserUpdateRequest
    if err := c.ShouldBindJSON(&userUpdate); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    err := h.userService.UpdateUser(userID, userUpdate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}