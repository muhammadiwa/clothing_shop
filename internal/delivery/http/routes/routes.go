package routes

import (
	"fashion-shop/internal/config"
	"fashion-shop/internal/delivery/http/handler"
	"fashion-shop/internal/delivery/http/middleware"
	"fashion-shop/internal/domain/usecase/impl"
	"fashion-shop/internal/infrastructure/auth"
	"fashion-shop/internal/infrastructure/persistence"
	"fashion-shop/internal/infrastructure/third_party"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes registers all routes
func RegisterRoutes(router *gin.Engine, repos *persistence.Repositories, redisClient *redis.Client, cfg *config.Config) {
	// Initialize services
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.Secret,
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
		1*time.Hour,
	)

	emailService := auth.NewSMTPEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	// Initialize third-party services
	rajaOngkirService := third_party.NewRajaOngkirService(cfg.RajaOngkir.APIKey, cfg.RajaOngkir.URL)
	midtransService := third_party.NewMidtransService(
		cfg.Midtrans.ServerKey,
		cfg.Midtrans.ClientKey,
		cfg.Midtrans.Environment,
	)

	// Initialize use cases
	userUseCase := impl.NewUserUseCase(repos.User, jwtService, emailService)
	addressUseCase := impl.NewAddressUseCase(repos.Address)
	productUseCase := impl.NewProductUseCase(repos.Product, repos.ProductImage, repos.ProductVariant, repos.Category)
	categoryUseCase := impl.NewCategoryUseCase(repos.Category)
	reviewUseCase := impl.NewReviewUseCase(repos.Review)
	cartUseCase := impl.NewCartUseCase(repos.Cart, repos.Product, repos.ProductVariant)
	wishlistUseCase := impl.NewWishlistUseCase(repos.Wishlist, repos.Product)
	orderUseCase := impl.NewOrderUseCase(
		repos.Order,
		repos.OrderItem,
		repos.Cart,
		repos.Product,
		repos.ProductVariant,
		repos.Address,
	)
	paymentUseCase := impl.NewPaymentUseCase(repos.Payment, repos.Order, midtransService)
	shippingUseCase := impl.NewShippingUseCase(rajaOngkirService)
	notificationUseCase := impl.NewNotificationUseCase(repos.Notification)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase, addressUseCase)
	productHandler := handler.NewProductHandler(productUseCase, categoryUseCase, reviewUseCase)
	cartHandler := handler.NewCartHandler(cartUseCase)
	wishlistHandler := handler.NewWishlistHandler(wishlistUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase, paymentUseCase, shippingUseCase)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// API versioning
	v1 := router.Group("/api/v1")

	// Public routes
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
			auth.POST("/forgot-password", userHandler.RequestPasswordReset)
			auth.POST("/reset-password", userHandler.ResetPassword)
		}

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.ListProducts)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/best-sellers", productHandler.GetBestSellers)
			products.GET("/new-arrivals", productHandler.GetNewArrivals)
			products.GET("/top-rated", productHandler.GetTopRated)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/slug/:slug", productHandler.GetProductBySlug)
			products.GET("/:id/reviews", productHandler.GetProductReviews)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", productHandler.ListCategories)
			categories.GET("/:id", productHandler.GetCategoryByID)
			categories.GET("/slug/:slug", productHandler.GetCategoryBySlug)
		}

		// Shipping routes
		shipping := v1.Group("/shipping")
		{
			shipping.GET("/provinces", orderHandler.GetProvinces)
			shipping.GET("/cities", orderHandler.GetCities)
		}

		// Payment webhook
		v1.POST(cfg.Midtrans.PaymentWebhook, paymentHandler.HandleWebhook)
	}

	// Protected routes (require authentication)
	protected := v1.Group("")
	protected.Use(authMiddleware.RequireAuth())
	{
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
			user.POST("/logout", userHandler.Logout)

			// Address routes
			addresses := user.Group("/addresses")
			{
				addresses.GET("", userHandler.GetAddresses)
				addresses.POST("", userHandler.CreateAddress)
				addresses.GET("/:id", userHandler.GetAddressByID)
				addresses.PUT("/:id", userHandler.UpdateAddress)
				addresses.DELETE("/:id", userHandler.DeleteAddress)
				addresses.PUT("/:id/default", userHandler.SetDefaultAddress)
				addresses.GET("/default", userHandler.GetDefaultAddress)
			}
		}

		// Cart routes
		cart := protected.Group("/cart")
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("/items", cartHandler.AddToCart)
			cart.PUT("/items/:id", cartHandler.UpdateCartItem)
			cart.DELETE("/items/:id", cartHandler.RemoveFromCart)
			cart.DELETE("", cartHandler.ClearCart)
		}

		// Wishlist routes
		wishlist := protected.Group("/wishlist")
		{
			wishlist.GET("", wishlistHandler.GetWishlist)
			wishlist.POST("", wishlistHandler.AddToWishlist)
			wishlist.DELETE("/:id", wishlistHandler.RemoveFromWishlist)
			wishlist.GET("/check/:product_id", wishlistHandler.IsInWishlist)
		}

		// Order routes
		orders := protected.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.GET("/number/:number", orderHandler.GetOrderByNumber)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		}

		// Payment routes
		payments := protected.Group("/payments")
		{
			payments.POST("", paymentHandler.ProcessPayment)
			payments.GET("/:id", paymentHandler.GetPaymentByID)
			payments.GET("/order/:order_id", paymentHandler.GetPaymentByOrderID)
		}

		// Shipping calculation
		shipping := protected.Group("/shipping")
		{
			shipping.POST("/calculate", orderHandler.CalculateShipping)
			shipping.GET("/track/:courier/:waybill", orderHandler.TrackShipment)
		}

		// Review routes
		reviews := protected.Group("/reviews")
		{
			reviews.POST("", productHandler.CreateReview)
			reviews.GET("/user", productHandler.GetUserReviews)
			reviews.PUT("/:id", productHandler.UpdateReview)
			reviews.DELETE("/:id", productHandler.DeleteReview)
		}

		// Notification routes
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/unread", notificationHandler.GetUnreadNotifications)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		}
	}

	// Admin routes (require admin role)
	admin := v1.Group("/admin")
	admin.Use(authMiddleware.RequireAdmin())
	{
		// User management
		users := admin.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id/active", userHandler.ToggleUserActive)
			users.PUT("/:id/reset-password", userHandler.ResetUserPassword)
		}

		// Product management
		products := admin.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
			products.POST("/bulk-upload", productHandler.BulkUploadProducts)

			// Product images
			products.POST("/:id/images", productHandler.UploadProductImage)
			products.DELETE("/images/:id", productHandler.DeleteProductImage)
			products.PUT("/images/:id/primary", productHandler.SetPrimaryImage)

			// Product variants
			products.POST("/:id/variants", productHandler.AddVariant)
			products.PUT("/variants/:id", productHandler.UpdateVariant)
			products.DELETE("/variants/:id", productHandler.DeleteVariant)
			products.PUT("/variants/:id/stock", productHandler.UpdateStock)
		}

		// Category management
		categories := admin.Group("/categories")
		{
			categories.POST("", productHandler.CreateCategory)
			categories.PUT("/:id", productHandler.UpdateCategory)
			categories.DELETE("/:id", productHandler.DeleteCategory)
			categories.POST("/:id/image", productHandler.UploadCategoryImage)
		}

		// Order management
		orders := admin.Group("/orders")
		{
			orders.GET("", orderHandler.GetAllOrders)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			orders.PUT("/:id/shipping", orderHandler.UpdateShippingInfo)
			orders.GET("/sales-report", orderHandler.GetSalesReport)
		}

		// Payment management
		payments := admin.Group("/payments")
		{
			payments.GET("", paymentHandler.ListPayments)
			payments.POST("/:id/refund", paymentHandler.RefundPayment)
		}
	}
}
