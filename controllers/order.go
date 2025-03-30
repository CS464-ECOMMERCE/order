package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"order/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderController handles HTTP requests for orders
type OrderController struct {
	orderService *services.OrderService
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewOrderController creates a new order controller
func NewOrderController() *OrderController {
	return &OrderController{
		orderService: services.NewOrderService(),
	}
}

// SetupRoutes registers the routes for the order controller
func (c *OrderController) SetupRoutes(router *gin.Engine) {
	orders := router.Group("/api/orders")
	{
		orders.POST("/checkout", c.PlaceOrder)
		orders.GET("/:id", c.GetOrder)
		orders.GET("/user/:userId", c.GetOrdersByUser)
		orders.GET("/merchant/:merchantId", c.GetOrdersByMerchant)
		orders.PUT("/:id/status", c.UpdateOrderStatus)
		orders.POST("/:id/cancel", c.CancelOrder)
		orders.DELETE("/:id", c.DeleteOrder)
	}
}

// PlaceOrder handles the checkout process
func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	// Extract session ID from cookie or header
	sessionId, err := ctx.Cookie("session_id")
	if err != nil {
		sessionId = ctx.GetHeader("X-Session-ID")
		if sessionId == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Session ID is required",
			})
			return
		}
	}

	// Extract user ID from request
	userIdStr := ctx.Query("userId")
	if userIdStr == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Place the order
	order, err := c.orderService.PlaceOrder(sessionId, userId)
	if err != nil {
		log.Printf("Error placing order: %v", err)
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to place order: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Order placed successfully",
		Data:    order,
	})
}

// GetOrder retrieves an order by ID
func (c *OrderController) GetOrder(ctx *gin.Context) {
	// Extract order ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	// Get the order
	order, err := c.orderService.GetOrder(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Order not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    order,
	})
}

// GetOrdersByUser retrieves all orders for a user
func (c *OrderController) GetOrdersByUser(ctx *gin.Context) {
	// Extract user ID from path
	userIdStr := ctx.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get the orders
	orders, err := c.orderService.GetOrdersByUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to get orders",
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    orders,
	})
}

// GetOrdersByMerchant retrieves all orders for a merchant
func (c *OrderController) GetOrdersByMerchant(ctx *gin.Context) {
	// Extract merchant ID from path
	merchantIdStr := ctx.Param("merchantId")
	merchantId, err := strconv.ParseUint(merchantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid merchant ID",
		})
		return
	}

	// Get the orders
	orders, err := c.orderService.GetOrdersByMerchant(merchantId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to get orders",
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    orders,
	})
}

// UpdateOrderStatus updates the status of an order
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	// Extract order ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	// Extract status from request body
	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Update the status
	if err := c.orderService.UpdateOrderStatus(id, request.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to update order status: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Order status updated successfully",
	})
}

// CancelOrder cancels an order
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	// Extract order ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	// Cancel the order
	if err := c.orderService.CancelOrder(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to cancel order: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Order cancelled successfully",
	})
}

// DeleteOrder deletes an order
func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	// Extract order ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	// Delete the order
	if err := c.orderService.DeleteOrder(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to delete order: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Order deleted successfully",
	})
}
