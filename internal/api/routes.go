package api

import (
    "github.com/gin-gonic/gin"
    "clothing-shop-api/internal/api/handlers"
    "clothing-shop-api/internal/api/middleware"
)

func SetupRoutes(router *gin.Engine) {
    // Public routes
    router.POST("/auth/login", handlers.Login)
    router.POST("/auth/register", handlers.Register)

    // Protected routes
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    
    protected.GET("/products", handlers.GetProducts)
    protected.GET("/products/:id", handlers.GetProductByID)
    protected.POST("/products", handlers.CreateProduct)
    protected.PUT("/products/:id", handlers.UpdateProduct)
    protected.DELETE("/products/:id", handlers.DeleteProduct)

    protected.GET("/orders", handlers.GetOrders)
    protected.POST("/orders", handlers.CreateOrder)

    protected.GET("/users/:id", handlers.GetUserProfile)
    protected.PUT("/users/:id", handlers.UpdateUserProfile)
}