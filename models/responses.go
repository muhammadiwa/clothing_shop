package models

// RegisterResponse represents a user registration response
type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// LoginResponse represents a user login response
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// ForgotPasswordResponse represents a forgot password response
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordResponse represents a reset password response
type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// LogoutResponse represents a logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

// ChangePasswordResponse represents a change password response
type ChangePasswordResponse struct {
	Message string `json:"message"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// UpdateProfileResponse represents an update profile response
type UpdateProfileResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

// AddressResponse represents an address response
type AddressResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	IsDefault  bool   `json:"is_default"`
}

// ProductResponse represents a product response
type ProductResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Price       float64                  `json:"price"`
	Stock       int                      `json:"stock"`
	SKU         string                   `json:"sku"`
	Weight      float64                  `json:"weight"`
	CategoryID  string                   `json:"category_id"`
	Category    CategoryResponse         `json:"category"`
	Images      []ProductImageResponse   `json:"images"`
	Variants    []ProductVariantResponse `json:"variants"`
	Tags        []TagResponse            `json:"tags"`
	IsActive    bool                     `json:"is_active"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

// ProductImageResponse represents a product image response
type ProductImageResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
}

// ProductVariantResponse represents a product variant response
type ProductVariantResponse struct {
	ID    string  `json:"id"`
	Size  string  `json:"size"`
	Color string  `json:"color"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
	SKU   string  `json:"sku"`
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ParentID    *string           `json:"parent_id"`
	Parent      *CategoryResponse `json:"parent,omitempty"`
}

// TagResponse represents a tag response
type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CartResponse represents a cart response
type CartResponse struct {
	ID    string             `json:"id"`
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}

// CartItemResponse represents a cart item response
type CartItemResponse struct {
	ID       string                  `json:"id"`
	Product  ProductResponse         `json:"product"`
	Variant  *ProductVariantResponse `json:"variant,omitempty"`
	Quantity int                     `json:"quantity"`
	Price    float64                 `json:"price"`
	Total    float64                 `json:"total"`
}

// WishlistResponse represents a wishlist response
type WishlistResponse struct {
	ID    string                 `json:"id"`
	Items []WishlistItemResponse `json:"items"`
}

// WishlistItemResponse represents a wishlist item response
type WishlistItemResponse struct {
	ID      string                  `json:"id"`
	Product ProductResponse         `json:"product"`
	Variant *ProductVariantResponse `json:"variant,omitempty"`
}

// OrderResponse represents an order response
type OrderResponse struct {
	ID              string              `json:"id"`
	OrderNumber     string              `json:"order_number"`
	Status          string              `json:"status"`
	Items           []OrderItemResponse `json:"items"`
	ShippingAddress AddressResponse     `json:"shipping_address"`
	ShippingMethod  string              `json:"shipping_method"`
	ShippingCost    float64             `json:"shipping_cost"`
	TrackingNumber  string              `json:"tracking_number"`
	SubTotal        float64             `json:"sub_total"`
	Tax             float64             `json:"tax"`
	Discount        float64             `json:"discount"`
	GrandTotal      float64             `json:"grand_total"`
	Notes           string              `json:"notes"`
	Payment         *PaymentResponse    `json:"payment,omitempty"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

// OrderItemResponse represents an order item response
type OrderItemResponse struct {
	ID       string                  `json:"id"`
	Product  ProductResponse         `json:"product"`
	Variant  *ProductVariantResponse `json:"variant,omitempty"`
	Quantity int                     `json:"quantity"`
	Price    float64                 `json:"price"`
	Total    float64                 `json:"total"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID             string  `json:"id"`
	Amount         float64 `json:"amount"`
	Method         string  `json:"method"`
	Status         string  `json:"status"`
	TransactionID  string  `json:"transaction_id"`
	PaymentGateway string  `json:"payment_gateway"`
	RefundAmount   float64 `json:"refund_amount"`
	RefundReason   string  `json:"refund_reason"`
	RefundDate     string  `json:"refund_date"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ReviewResponse represents a review response
type ReviewResponse struct {
	ID        string          `json:"id"`
	Product   ProductResponse `json:"product"`
	User      UserResponse    `json:"user"`
	Rating    int             `json:"rating"`
	Comment   string          `json:"comment"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Message     string  `json:"message"`
	Type        string  `json:"type"`
	IsRead      bool    `json:"is_read"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
	CreatedAt   string  `json:"created_at"`
}

// ProvinceResponse represents a province response
type ProvinceResponse struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

// CityResponse represents a city response
type CityResponse struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

// ShippingCostResponse represents a shipping cost response
type ShippingCostResponse struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        []struct {
		Value int    `json:"value"`
		Etd   string `json:"etd"`
		Note  string `json:"note"`
	} `json:"cost"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}

// PaginationResponse represents a pagination response
type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}
