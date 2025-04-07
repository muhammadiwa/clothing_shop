package handlers

import (
	"net/http"
	"strconv"

	"clothing-shop-api/internal/domain/services"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// AddToCart handles adding an item to the cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ProductVariantID uint `json:"product_variant_id" binding:"required"`
		Quantity         int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	err := h.cartService.AddToCart(userID.(uint), req.ProductVariantID, req.Quantity)
	if err != nil {
		if err == services.ErrInsufficientStock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
}

// UpdateCartItem handles updating a cart item quantity
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get cart item to check if it belongs to the user
	cartItems, err := h.cartService.GetUserCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var authorized bool
	for _, item := range cartItems {
		if item.ID == uint(id) {
			authorized = true
			break
		}
	}

	if !authorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this cart item"})
		return
	}

	err = h.cartService.UpdateCartItem(uint(id), req.Quantity)
	if err != nil {
		if err == services.ErrCartItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		if err == services.ErrInsufficientStock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

// RemoveFromCart handles removing an item from the cart
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	// Get cart item to check if it belongs to the user
	cartItems, err := h.cartService.GetUserCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var authorized bool
	for _, item := range cartItems {
		if item.ID == uint(id) {
			authorized = true
			break
		}
	}

	if !authorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this cart item"})
		return
	}

	err = h.cartService.RemoveFromCart(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// GetCart handles retrieving all items in the user's cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cartItems, err := h.cartService.GetUserCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate totals
	var (
		totalItems     int
		subtotal       float64
		discountAmount float64
	)

	for _, item := range cartItems {
		totalItems += item.Quantity

		// Calculate item price considering discounts
		itemPrice := item.ProductVariant.Price
		if item.ProductVariant.Product != nil && item.ProductVariant.Product.Discount > 0 {
			discountMultiplier := 1.0 - (item.ProductVariant.Product.Discount / 100.0)
			itemPrice = item.ProductVariant.Price * discountMultiplier

			// Track discount amount
			discountAmount += (item.ProductVariant.Price - itemPrice) * float64(item.Quantity)
		}

		subtotal += itemPrice * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"items": cartItems,
		"summary": gin.H{
			"total_items":     totalItems,
			"subtotal":        subtotal,
			"discount_amount": discountAmount,
		},
	})
}
