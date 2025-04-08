package models

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UpdateProfileRequest represents an update profile request
type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

// AddressRequest represents an address request
type AddressRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	IsDefault  bool   `json:"is_default"`
}

// ProductRequest represents a product request
type ProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	SKU         string  `json:"sku" validate:"required"`
	Weight      float64 `json:"weight" validate:"required,gt=0"`
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	IsActive    bool    `json:"is_active"`
}

// ProductVariantRequest represents a product variant request
type ProductVariantRequest struct {
	Size  string  `json:"size" validate:"required"`
	Color string  `json:"color" validate:"required"`
	Stock int     `json:"stock" validate:"required,gte=0"`
	Price float64 `json:"price" validate:"omitempty,gte=0"`
	SKU   string  `json:"sku" validate:"required"`
}

// CategoryRequest represents a category request
type CategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id" validate:"omitempty,uuid"`
}

// CartItemRequest represents a cart item request
type CartItemRequest struct {
	ProductID string  `json:"product_id" validate:"required,uuid"`
	VariantID *string `json:"variant_id" validate:"omitempty,uuid"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
}

// WishlistItemRequest represents a wishlist item request
type WishlistItemRequest struct {
	ProductID string  `json:"product_id" validate:"required,uuid"`
	VariantID *string `json:"variant_id" validate:"omitempty,uuid"`
}

// OrderRequest represents an order request
type OrderRequest struct {
	ShippingAddressID string  `json:"shipping_address_id" validate:"required,uuid"`
	ShippingMethod    string  `json:"shipping_method" validate:"required"`
	ShippingCost      float64 `json:"shipping_cost" validate:"required,gte=0"`
	Notes             string  `json:"notes"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID        string `json:"order_id" validate:"required,uuid"`
	Method         string `json:"method" validate:"required"`
	PaymentGateway string `json:"payment_gateway" validate:"required"`
	PaymentDetails string `json:"payment_details"`
}

// ReviewRequest represents a review request
type ReviewRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}

// ShippingCostRequest represents a shipping cost request
type ShippingCostRequest struct {
	Origin      int    `json:"origin" validate:"required"`
	Destination int    `json:"destination" validate:"required"`
	Weight      int    `json:"weight" validate:"required,gt=0"`
	Courier     string `json:"courier" validate:"required"`
}
