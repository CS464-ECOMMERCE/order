package controllers

import (
	"context"
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

// UpdatePaymentStatus implements UpdatePaymentStatus RPC method
func (s *OrderController) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.Empty, error) {
	if err := s.orderService.UpdatePaymentStatus(req.Event, req.OrderId); err != nil {
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

	var status pb.OrderStatus
	switch order.Status {
	case "completed":
		status = pb.OrderStatus_ORDER_STATUS_PROCESSING
	case "cancelled":
		status = pb.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		status = pb.OrderStatus_ORDER_STATUS_PROCESSING
	}

	return &pb.Order{
		Id:         order.Id,
		UserId:     order.UserId,
		Total:      order.Total,
		Status:     status,
		OrderItems: orderItems,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}
}
