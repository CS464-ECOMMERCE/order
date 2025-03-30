package services

import (
	"context"
	"order/models"
	pb "order/proto"
	"time"
)

// GRPCServer implements the OrderService gRPC service
type GRPCServer struct {
	pb.UnimplementedOrderServiceServer
	orderService *OrderService
}

// NewGRPCServer creates a new gRPC server with the order service
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		orderService: NewOrderService(),
	}
}

// PlaceOrder implements the PlaceOrder RPC method
func (s *GRPCServer) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.Order, error) {
	order, err := s.orderService.PlaceOrder(req.SessionId, req.UserId)
	if err != nil {
		return nil, err
	}
	return convertToProtoOrder(order), nil
}

// GetOrder implements the GetOrder RPC method
func (s *GRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.orderService.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}
	return convertToProtoOrder(order), nil
}

// GetOrdersByUser implements the GetOrdersByUser RPC method
func (s *GRPCServer) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.GetOrdersResponse, error) {
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
func (s *GRPCServer) GetOrdersByMerchant(ctx context.Context, req *pb.GetOrdersByMerchantRequest) (*pb.GetOrdersResponse, error) {
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
func (s *GRPCServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
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
func (s *GRPCServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
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
func (s *GRPCServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.Empty, error) {
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
