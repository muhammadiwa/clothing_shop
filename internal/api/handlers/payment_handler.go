package handlers

import (
	"clothing-shop-api/internal/domain/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// @Summary Create payment for order
// @Description Creates a new payment for the specified order
// @Tags payments
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Param payment body CreatePaymentRequest true "Payment details"
// @Success 201 {object} models.Payment
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /payments/{order_id} [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.CreatePayment(uint(orderID), req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// @Summary Get payment status
// @Description Get the current status of a payment
// @Tags payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} models.Payment
// @Failure 404 {object} ErrorResponse
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.paymentService.GetPaymentStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if payment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Handle payment notification webhook
// @Description Handle payment status notification from payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Param notification body payment.PaymentNotification true "Payment notification"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /payments/notification [post]
func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	var notification payment.PaymentNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.paymentService.HandlePaymentNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification processed successfully"})
}
