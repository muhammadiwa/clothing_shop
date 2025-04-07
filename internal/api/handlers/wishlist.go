package handlers

import (
	"net/http"
	"strconv"

	"clothing-shop-api/internal/domain/services"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	wishlistService *services.WishlistService
}

func NewWishlistHandler(wishlistService *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

// AddToWishlist handles adding an item to the wishlist
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.wishlistService.AddToWishlist(userID.(uint), req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to wishlist"})
}

// RemoveFromWishlist handles removing an item from the wishlist
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist item ID"})
		return
	}

	// Get wishlist to check if the item belongs to the user
	wishlistItems, err := h.wishlistService.GetUserWishlist(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var authorized bool
	for _, item := range wishlistItems {
		if item.ID == uint(id) {
			authorized = true
			break
		}
	}

	if !authorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this wishlist item"})
		return
	}

	err = h.wishlistService.RemoveFromWishlist(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from wishlist"})
}

// GetWishlist handles retrieving all items in the user's wishlist
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	wishlistItems, err := h.wishlistService.GetUserWishlist(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": wishlistItems})
}
