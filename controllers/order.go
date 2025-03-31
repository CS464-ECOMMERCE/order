package controllers

import (
	"context"
	"fmt"
	"order/models"
	pb "order/proto"
	"order/services"
	"time"
)

// OrderController handles HTTP requests for orders
type OrderController struct {
	pb.UnimplementedOrderServiceServer
	orderService *services.OrderService
}

// // Response represents a standard API response
// type Response struct {
// 	Success bool        `json:"success"`
// 	Message string      `json:"message,omitempty"`
// 	Data    interface{} `json:"data,omitempty"`
// }

// NewOrderController creates a new order controller
func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// PlaceOrder implements the PlaceOrder RPC method
func (s *OrderController) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.Order, error) {
	fmt.Println("PlaceOrder==================")
	order, err := s.orderService.PlaceOrder(req.SessionId, req.UserId)
	if err != nil {
		return nil, err
	}
	return convertToProtoOrder(order), nil
}

// GetOrder implements the GetOrder RPC method
func (s *OrderController) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.orderService.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}
	return convertToProtoOrder(order), nil
}

// GetOrdersByUser implements the GetOrdersByUser RPC method
func (s *OrderController) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.GetOrdersResponse, error) {
	orders, err := s.orderService.GetOrdersByUser(req.UserId)
	if err != nil {
		return nil, err
	}

	protoOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = convertToProtoOrder(order)
	}

	return &pb.GetOrdersResponse{
		Orders: protoOrders,
	}, nil
}

// GetOrdersByMerchant implements the GetOrdersByMerchant RPC method
func (s *OrderController) GetOrdersByMerchant(ctx context.Context, req *pb.GetOrdersByMerchantRequest) (*pb.GetOrdersResponse, error) {
	orders, err := s.orderService.GetOrdersByMerchant(req.MerchantId)
	if err != nil {
		return nil, err
	}

	protoOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = convertToProtoOrder(order)
	}

	return &pb.GetOrdersResponse{
		Orders: protoOrders,
	}, nil
}

// UpdateOrderStatus implements the UpdateOrderStatus RPC method
func (s *OrderController) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	if err := s.orderService.UpdateOrderStatus(req.Id, req.Status); err != nil {
		return nil, err
	}

	order, err := s.orderService.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return convertToProtoOrder(order), nil
}

// CancelOrder implements the CancelOrder RPC method
func (s *OrderController) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
	if err := s.orderService.CancelOrder(req.Id); err != nil {
		return nil, err
	}

	order, err := s.orderService.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return convertToProtoOrder(order), nil
}

// DeleteOrder implements the DeleteOrder RPC method
func (s *OrderController) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.Empty, error) {
	if err := s.orderService.DeleteOrder(req.Id); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

// convertToProtoOrder converts a model order to a protobuf order
func convertToProtoOrder(order *models.Order) *pb.Order {
	orderItems := make([]*pb.OrderItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = &pb.OrderItem{
			OrderId:   item.OrderId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &pb.Order{
		Id:         order.Id,
		UserId:     order.UserId,
		Total:      order.Total,
		Status:     order.Status,
		OrderItems: orderItems,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}
}

// // PlaceOrder handles the checkout process
// func (c *OrderController) PlaceOrder(ctx *gin.Context) {
// 	// Extract session ID from cookie or header
// 	sessionId, err := ctx.Cookie("session_id")
// 	if err != nil {
// 		sessionId = ctx.GetHeader("X-Session-ID")
// 		if sessionId == "" {
// 			ctx.JSON(http.StatusBadRequest, Response{
// 				Success: false,
// 				Message: "Session ID is required",
// 			})
// 			return
// 		}
// 	}
//
// 	// Extract user ID from request
// 	userIdStr := ctx.Query("userId")
// 	if userIdStr == "" {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "User ID is required",
// 		})
// 		return
// 	}
//
// 	userId, err := strconv.ParseUint(userIdStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid user ID",
// 		})
// 		return
// 	}
//
// 	// Place the order
// 	order, err := c.orderService.PlaceOrder(sessionId, userId)
// 	if err != nil {
// 		log.Printf("Error placing order: %v", err)
// 		ctx.JSON(http.StatusInternalServerError, Response{
// 			Success: false,
// 			Message: "Failed to place order: " + err.Error(),
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusCreated, Response{
// 		Success: true,
// 		Message: "Order placed successfully",
// 		Data:    order,
// 	})
// }
//
// // GetOrder retrieves an order by ID
// func (c *OrderController) GetOrder(ctx *gin.Context) {
// 	// Extract order ID from path
// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid order ID",
// 		})
// 		return
// 	}
//
// 	// Get the order
// 	order, err := c.orderService.GetOrder(id)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, Response{
// 			Success: false,
// 			Message: "Order not found",
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, Response{
// 		Success: true,
// 		Data:    order,
// 	})
// }
//
// // GetOrdersByUser retrieves all orders for a user
// func (c *OrderController) GetOrdersByUser(ctx *gin.Context) {
// 	// Extract user ID from path
// 	userIdStr := ctx.Param("userId")
// 	userId, err := strconv.ParseUint(userIdStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid user ID",
// 		})
// 		return
// 	}
//
// 	// Get the orders
// 	orders, err := c.orderService.GetOrdersByUser(userId)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, Response{
// 			Success: false,
// 			Message: "Failed to get orders",
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, Response{
// 		Success: true,
// 		Data:    orders,
// 	})
// }
//
// // GetOrdersByMerchant retrieves all orders for a merchant
// func (c *OrderController) GetOrdersByMerchant(ctx *gin.Context) {
// 	// Extract merchant ID from path
// 	merchantIdStr := ctx.Param("merchantId")
// 	merchantId, err := strconv.ParseUint(merchantIdStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid merchant ID",
// 		})
// 		return
// 	}
//
// 	// Get the orders
// 	orders, err := c.orderService.GetOrdersByMerchant(merchantId)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, Response{
// 			Success: false,
// 			Message: "Failed to get orders",
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, Response{
// 		Success: true,
// 		Data:    orders,
// 	})
// }
//
// // UpdateOrderStatus updates the status of an order
// func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
// 	// Extract order ID from path
// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid order ID",
// 		})
// 		return
// 	}
//
// 	// Extract status from request body
// 	var request struct {
// 		Status string `json:"status"`
// 	}
//
// 	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid request body",
// 		})
// 		return
// 	}
//
// 	// Update the status
// 	if err := c.orderService.UpdateOrderStatus(id, request.Status); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, Response{
// 			Success: false,
// 			Message: "Failed to update order status: " + err.Error(),
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, Response{
// 		Success: true,
// 		Message: "Order status updated successfully",
// 	})
// }
//
// // CancelOrder cancels an order
// func (c *OrderController) CancelOrder(ctx *gin.Context) {
// 	// Extract order ID from path
// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid order ID",
// 		})
// 		return
// 	}
//
// 	// Cancel the order
// 	if err := c.orderService.CancelOrder(id); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, Response{
// 			Success: false,
// 			Message: "Failed to cancel order: " + err.Error(),
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, Response{
// 		Success: true,
// 		Message: "Order cancelled successfully",
// 	})
// }
//
// // DeleteOrder deletes an order
// func (c *OrderController) DeleteOrder(ctx *gin.Context) {
// 	// Extract order ID from path
// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, Response{
// 			Success: false,
// 			Message: "Invalid order ID",
// 		})
// 		return
// 	}
//
// 	// Delete the order
// 	if err := c.orderService.DeleteOrder(id); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, Response{
// 			Success: false,
// 			Message: "Failed to delete order: " + err.Error(),
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, Response{
// 		Success: true,
// 		Message: "Order deleted successfully",
// 	})
// }
