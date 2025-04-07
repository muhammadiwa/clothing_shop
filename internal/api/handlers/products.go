package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "clothing-shop-api/internal/domain/models"
    "clothing-shop-api/internal/domain/services"
)

type ProductHandler struct {
    ProductService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
    return &ProductHandler{ProductService: productService}
}

// GetProducts retrieves all products
func (h *ProductHandler) GetProducts(c *gin.Context) {
    products, err := h.ProductService.GetAllProducts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, products)
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
    id := c.Param("id")
    product, err := h.ProductService.GetProductByID(id)
    if err != nil {
        if err == services.ErrProductNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    createdProduct, err := h.ProductService.CreateProduct(product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, createdProduct)
}

// UpdateProduct updates an existing product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
    id := c.Param("id")
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    updatedProduct, err := h.ProductService.UpdateProduct(id, product)
    if err != nil {
        if err == services.ErrProductNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct deletes a product by ID
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    err := h.ProductService.DeleteProduct(id)
    if err != nil {
        if err == services.ErrProductNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusNoContent, nil)
}