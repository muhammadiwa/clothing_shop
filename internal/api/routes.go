package api

import (
	"github.com/gin-gonic/gin"

	"clothing-shop-api/internal/api/handlers"
	"clothing-shop-api/internal/api/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/api/v1")

	// Authentication
	auth := public.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/verify", handlers.VerifyEmail)
	auth.POST("/forgot-password", handlers.ForgotPassword)
	auth.POST("/reset-password", handlers.ResetPassword)

	// Products and categories
	products := public.Group("/products")
	products.GET("", handlers.ListProducts)
	products.GET("/:id", handlers.GetProduct)
	products.GET("/categories", handlers.ListCategories)
	products.GET("/categories/:id", handlers.GetCategory)
	products.GET("/:id/reviews", handlers.GetProductReviews)

	// Location data
	locations := public.Group("/locations")
	locations.GET("/provinces", handlers.GetProvinces)
	locations.GET("/cities", handlers.GetCities)
	locations.GET("/districts", handlers.GetDistricts)

	// Protected routes - Customer
	customer := router.Group("/api/v1")
	customer.Use(middleware.AuthMiddleware())

	// User profile
	profile := customer.Group("/profile")
	profile.GET("", handlers.GetUserProfile)
	profile.PUT("", handlers.UpdateUserProfile)
	profile.GET("/addresses", handlers.ListUserAddresses)
	profile.POST("/addresses", handlers.AddUserAddress)
	profile.PUT("/addresses/:id", handlers.UpdateUserAddress)
	profile.DELETE("/addresses/:id", handlers.DeleteUserAddress)

	// Shopping
	cart := customer.Group("/cart")
	cart.GET("", handlers.GetCart)
	cart.POST("", handlers.AddToCart)
	cart.PUT("/:id", handlers.UpdateCartItem)
	cart.DELETE("/:id", handlers.RemoveFromCart)

	wishlist := customer.Group("/wishlist")
	wishlist.GET("", handlers.GetWishlist)
	wishlist.POST("", handlers.AddToWishlist)
	wishlist.DELETE("/:id", handlers.RemoveFromWishlist)

	// Orders
	orders := customer.Group("/orders")
	orders.POST("", handlers.CreateOrder)
	orders.GET("", handlers.GetUserOrders)
	orders.GET("/:id", handlers.GetOrderDetails)
	orders.POST("/shipping-cost", handlers.CalculateShippingCost)

	// Payment
	payments := customer.Group("/payments")
	payments.POST("/:order_id", handlers.CreatePayment)
	payments.GET("/:id", handlers.GetPaymentStatus)

	// Reviews
	reviews := customer.Group("/reviews")
	reviews.POST("", handlers.CreateReview)
	reviews.PUT("/:id", handlers.UpdateReview)

	// Support tickets
	tickets := customer.Group("/support")
	tickets.GET("", handlers.GetUserTickets)
	tickets.POST("", handlers.CreateTicket)
	tickets.GET("/:id", handlers.GetTicketDetails)
	tickets.POST("/:id/reply", handlers.AddTicketReply)
	tickets.PUT("/:id/close", handlers.CloseTicket)

	// Protected routes - Admin
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	// Product management
	adminProducts := admin.Group("/products")
	adminProducts.POST("", handlers.CreateProduct)
	adminProducts.PUT("/:id", handlers.UpdateProduct)
	adminProducts.DELETE("/:id", handlers.DeleteProduct)
	adminProducts.POST("/:id/images", handlers.AddProductImage)
	adminProducts.DELETE("/images/:id", handlers.DeleteProductImage)

	// Category management
	adminCategories := admin.Group("/categories")
	adminCategories.POST("", handlers.CreateCategory)
	adminCategories.PUT("/:id", handlers.UpdateCategory)
	adminCategories.DELETE("/:id", handlers.DeleteCategory)

	// User management
	adminUsers := admin.Group("/users")
	adminUsers.GET("", handlers.ListUsers)
	adminUsers.GET("/:id", handlers.GetUser)
	adminUsers.PUT("/:id", handlers.UpdateUser)
	adminUsers.DELETE("/:id", handlers.DeleteUser)

	// Order management
	adminOrders := admin.Group("/orders")
	adminOrders.GET("", handlers.ListAllOrders)
	adminOrders.GET("/:id", handlers.GetOrderAdmin)
	adminOrders.PUT("/:id/status", handlers.UpdateOrderStatus)

	// Support ticket management
	adminTickets := admin.Group("/support")
	adminTickets.GET("", handlers.ListAllTickets)
	adminTickets.GET("/:id", handlers.GetTicketAdmin)
	adminTickets.POST("/:id/reply", handlers.AddAdminTicketReply)
	adminTickets.PUT("/:id/status", handlers.UpdateTicketStatus)

	// Sales reports
	adminReports := admin.Group("/reports")
	adminReports.GET("/sales", handlers.GetSalesReport)
	adminReports.GET("/products", handlers.GetProductsReport)

	// Webhook routes
	webhooks := public.Group("/webhooks")
	webhooks.POST("/midtrans", handlers.MidtransWebhook)
}
