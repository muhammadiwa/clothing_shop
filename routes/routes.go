package routes

import (
	"github.com/fashion-shop/config"
	"github.com/fashion-shop/controllers"
	"github.com/fashion-shop/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// API version group
	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", controllers.HealthCheck)

	// Public routes
	setupPublicRoutes(api, cfg)

	// Protected routes (require authentication)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	setupProtectedRoutes(protected, cfg)

	// Admin routes (require admin role)
	admin := api.Group("/admin")
	admin.Use(middleware.AdminMiddleware())
	setupAdminRoutes(admin, cfg)

	// Swagger documentation
	router.GET("/swagger/*any", controllers.SwaggerHandler)
}

// setupPublicRoutes configures all public routes
func setupPublicRoutes(router *gin.RouterGroup, cfg *config.Config) {
	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)
	}

	// Product routes
	products := router.Group("/products")
	{
		products.GET("/", controllers.ListProducts)
		products.GET("/:id", controllers.GetProduct)
		products.GET("/search", controllers.SearchProducts)
	}

	// Category routes
	categories := router.Group("/categories")
	{
		categories.GET("/", controllers.ListCategories)
		categories.GET("/:id", controllers.GetCategory)
		categories.GET("/:id/products", controllers.GetCategoryProducts)
	}

	// Shipping routes
	shipping := router.Group("/shipping")
	{
		shipping.GET("/provinces", controllers.GetProvinces)
		shipping.GET("/cities", controllers.GetCities)
		shipping.POST("/cost", controllers.CalculateShippingCost)
	}
}

// setupProtectedRoutes configures all protected routes (require authentication)
func setupProtectedRoutes(router *gin.RouterGroup, cfg *config.Config) {
	// User profile routes
	profile := router.Group("/profile")
	{
		profile.GET("/", controllers.GetProfile)
		profile.PUT("/", controllers.UpdateProfile)
		profile.GET("/orders", controllers.GetUserOrders)
		profile.GET("/addresses", controllers.GetUserAddresses)
		profile.POST("/addresses", controllers.AddUserAddress)
		profile.PUT("/addresses/:id", controllers.UpdateUserAddress)
		profile.DELETE("/addresses/:id", controllers.DeleteUserAddress)
	}

	// Cart routes
	cart := router.Group("/cart")
	{
		cart.GET("/", controllers.GetCart)
		cart.POST("/items", controllers.AddCartItem)
		cart.PUT("/items/:id", controllers.UpdateCartItem)
		cart.DELETE("/items/:id", controllers.DeleteCartItem)
	}

	// Wishlist routes
	wishlist := router.Group("/wishlist")
	{
		wishlist.GET("/", controllers.GetWishlist)
		wishlist.POST("/items", controllers.AddWishlistItem)
		wishlist.DELETE("/items/:id", controllers.DeleteWishlistItem)
	}

	// Order routes
	orders := router.Group("/orders")
	{
		orders.POST("/", controllers.CreateOrder)
		orders.GET("/:id", controllers.GetOrder)
		orders.GET("/:id/track", controllers.TrackOrder)
	}

	// Payment routes
	payments := router.Group("/payments")
	{
		payments.POST("/", controllers.CreatePayment)
		payments.GET("/:id", controllers.GetPayment)
	}

	// Review routes
	reviews := router.Group("/reviews")
	{
		reviews.POST("/", controllers.CreateReview)
		reviews.PUT("/:id", controllers.UpdateReview)
		reviews.DELETE("/:id", controllers.DeleteReview)
	}

	// Notification routes
	notifications := router.Group("/notifications")
	{
		notifications.GET("/", controllers.GetNotifications)
		notifications.PUT("/:id/read", controllers.MarkNotificationAsRead)
	}

	// Auth routes that require authentication
	auth := router.Group("/auth")
	{
		auth.POST("/logout", controllers.Logout)
		auth.PUT("/password", controllers.ChangePassword)
	}
}

// setupAdminRoutes configures all admin routes
func setupAdminRoutes(router *gin.RouterGroup, cfg *config.Config) {
	// Product management
	products := router.Group("/products")
	{
		products.POST("/", controllers.CreateProduct)
		products.PUT("/:id", controllers.UpdateProduct)
		products.DELETE("/:id", controllers.DeleteProduct)
		products.POST("/bulk", controllers.BulkUploadProducts)
		products.PUT("/:id/stock", controllers.UpdateProductStock)
		products.POST("/:id/images", controllers.UploadProductImage)
		products.DELETE("/:id/images/:imageId", controllers.DeleteProductImage)
	}

	// Category management
	categories := router.Group("/categories")
	{
		categories.POST("/", controllers.CreateCategory)
		categories.PUT("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}

	// Order management
	orders := router.Group("/orders")
	{
		orders.GET("/", controllers.ListAllOrders)
		orders.PUT("/:id/status", controllers.UpdateOrderStatus)
		orders.POST("/:id/refund", controllers.ProcessRefund)
	}

	// User management
	users := router.Group("/users")
	{
		users.GET("/", controllers.ListUsers)
		users.GET("/:id", controllers.GetUser)
		users.PUT("/:id/block", controllers.BlockUser)
		users.PUT("/:id/unblock", controllers.UnblockUser)
		users.POST("/:id/reset-password", controllers.AdminResetUserPassword)
		users.GET("/:id/activity", controllers.GetUserActivityLogs)
	}

	// Reports
	reports := router.Group("/reports")
	{
		reports.GET("/sales", controllers.GenerateSalesReport)
		reports.GET("/inventory", controllers.GenerateInventoryReport)
		reports.GET("/users", controllers.GenerateUserReport)
	}
}
