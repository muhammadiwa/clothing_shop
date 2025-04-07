package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"clothing-shop-api/internal/api/handlers"
	"clothing-shop-api/internal/api/middleware"
	"clothing-shop-api/internal/config"
	"clothing-shop-api/internal/domain/services"
	"clothing-shop-api/internal/repository"
	"clothing-shop-api/pkg/database"
)

// @title           Clothing Shop API
// @version         1.0
// @description     Indonesian online clothing store API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Connect to database
	database.Connect()
	defer database.Close()
	db := database.GetDB()

	// Set up repositories
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	variantRepo := repository.NewProductVariantRepository(db)
	imageRepo := repository.NewProductImageRepository(db)

	// Set up services using interfaces
	productService := services.NewProductService(productRepo, categoryRepo, variantRepo, imageRepo)

	// Set up handlers
	productHandler := handlers.NewProductHandler(productService)

	// Initialize router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Public routes
	v1 := router.Group("/api/v1")

	// Authentication
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/verify-email", authHandler.VerifyEmail)
	auth.POST("/forgot-password", authHandler.ForgotPassword)
	auth.POST("/reset-password", authHandler.ResetPassword)

	// Products
	products := v1.Group("/products")
	products.GET("", productHandler.GetProducts)
	products.GET("/:id", productHandler.GetProduct)
	products.GET("/categories", productHandler.GetCategories)

	// Protected routes - Customer
	customer := v1.Group("/")
	customer.Use(middleware.AuthMiddleware())

	// Cart
	// cart := customer.Group("/cart")
	// cart.GET("", cartHandler.GetCart)
	// cart.POST("", cartHandler.AddToCart)
	// cart.PUT("/:id", cartHandler.UpdateCartItem)
	// cart.DELETE("/:id", cartHandler.RemoveFromCart)

	// Wishlist
	// wishlist := customer.Group("/wishlist")
	// wishlist.GET("", wishlistHandler.GetWishlist)
	// wishlist.POST("", wishlistHandler.AddToWishlist)
	// wishlist.DELETE("/:id", wishlistHandler.RemoveFromWishlist)

	// Orders
	// orders := customer.Group("/orders")
	// orders.POST("", orderHandler.CreateOrder)
	// orders.GET("", orderHandler.GetUserOrders)
	// orders.GET("/:id", orderHandler.GetOrderDetails)

	// Payments
	// payments := customer.Group("/payments")
	// payments.POST("/:order_id", paymentHandler.CreatePayment)
	// payments.GET("/:id", paymentHandler.GetPaymentStatus)
	// payments.POST("/notification", paymentHandler.HandleNotification)

	// Protected routes - Admin
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	// Product management
	adminProducts := admin.Group("/products")
	adminProducts.POST("", productHandler.CreateProduct)
	adminProducts.PUT("/:id", productHandler.UpdateProduct)
	adminProducts.DELETE("/:id", productHandler.DeleteProduct)

	// Add Swagger documentation route before starting the server
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	log.Printf("Starting server on port %s...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
