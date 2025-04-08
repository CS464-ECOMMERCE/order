package controllers

import (
	"context"
	pb "order/proto"
	"order/services"
)

// OrderController handles HTTP requests for orders
type OrderController struct {
	pb.UnimplementedOrderServiceServer
	orderService *services.OrderService
}

// NewOrderController creates a new order controller
func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// GetOrder implements the GetOrder RPC method
func (s *OrderController) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.orderService.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetOrdersByUser implements the GetOrdersByUser RPC method
func (s *OrderController) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.GetOrdersResponse, error) {
	orders, err := s.orderService.GetOrdersByUser(req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetOrdersResponse{
		Orders: orders,
	}, nil
}

// GetOrdersByMerchant implements the GetOrdersByMerchant RPC method
func (s *OrderController) GetOrdersByMerchant(ctx context.Context, req *pb.GetOrdersByMerchantRequest) (*pb.GetOrdersResponse, error) {
	orders, err := s.orderService.GetOrdersByMerchant(req.MerchantId)
	if err != nil {
		return nil, err
	}

	return &pb.GetOrdersResponse{
		Orders: orders,
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

	return order, nil
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

	return order, nil
}

// UpdatePaymentStatus implements UpdatePaymentStatus RPC method
func (s *OrderController) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.Empty, error) {
	if err := s.orderService.UpdatePaymentStatus(req.Event, req.OrderId); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
