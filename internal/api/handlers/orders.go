package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "clothing-shop-api/internal/domain/models"
    "clothing-shop-api/internal/domain/services"
)

type OrderHandler struct {
    OrderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
    return &OrderHandler{OrderService: orderService}
}

// CreateOrder handles the creation of a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    createdOrder, err := h.OrderService.CreateOrder(order)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, createdOrder)
}

// GetOrder handles retrieving an order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
    orderID := c.Param("id")
    order, err := h.OrderService.GetOrderByID(orderID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }

    c.JSON(http.StatusOK, order)
}

// GetAllOrders handles retrieving all orders
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
    orders, err := h.OrderService.GetAllOrders()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, orders)
}