package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rkweber-max/checkout-backend/internal/checkout/domain"
	"github.com/rkweber-max/checkout-backend/internal/checkout/service"
)

type CheckoutHandler struct {
	service *service.CheckoutService
}

func NewCheckoutHandler(service *service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{service: service}
}

func (h *CheckoutHandler) Checkout(c *gin.Context) {
	var request domain.CheckoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}
	
	order, err := h.service.ProcessOrder(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
