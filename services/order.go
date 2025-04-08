package services

import (
	"errors"
	"order/models"
	pb "order/proto"
	"order/storage"
	"time"

	"github.com/stripe/stripe-go/v81"
	"gorm.io/gorm"
)

// OrderService provides operations for managing orders
type OrderService struct{}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{}
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(id uint64) (*pb.Order, error) {
	order, err := storage.GetInstance().Order.GetOrder(id, nil)
	if err != nil {
		return nil, err
	}

	return convertToProtoOrder(order), nil
}

// GetOrdersByUser retrieves all orders for a user
func (s *OrderService) GetOrdersByUser(userId uint64) ([]*pb.Order, error) {
	order, err := storage.GetInstance().Order.GetOrdersByUserId(userId)
	if err != nil {
		return nil, err
	}

	protoOrders := make([]*pb.Order, len(order))
	for i, ord := range order {
		protoOrders[i] = convertToProtoOrder(ord)
		if err != nil {
			return nil, err
		}
	}

	return protoOrders, nil
}

// GetOrdersByMerchant retrieves all orders containing products from a specific merchant
func (s *OrderService) GetOrdersByMerchant(merchantId uint64) ([]*pb.Order, error) {
	order, err := storage.GetInstance().Order.GetOrdersByMerchantId(merchantId)
	if err != nil {
		return nil, err
	}

	protoOrders := make([]*pb.Order, len(order))
	for i, ord := range order {
		protoOrders[i] = convertToProtoOrder(ord)
		if err != nil {
			return nil, err
		}
	}

	return protoOrders, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id uint64, status string) error {
	// Validate status
	if status != string(models.OrderStatusProcessing) &&
		status != string(models.OrderStatusCancelled) &&
		status != string(models.OrderStatusCompleted) {
		return errors.New("invalid status")
	}

	return storage.GetInstance().Order.UpdateOrderStatus(id, status)
}

// CancelOrder cancels an order and restores inventory
func (s *OrderService) CancelOrder(id uint64) error {
	tx := storage.GetInstance().BeginTransaction()
	// Get the order
	order, err := storage.GetInstance().Order.GetOrder(id, tx)
	if err != nil {
		return err
	}

	// Check if order is already cancelled or completed
	if order.Status == models.OrderStatusCancelled {
		return errors.New("order is already cancelled")
	}
	if order.Status == models.OrderStatusCompleted {
		return errors.New("cannot cancel a completed order")
	}

	err = s.handleRevertOrderItems(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// UpdatePaymentStatus updates an order payment status
// and, potentially rolling back an unpaid order
func (s *OrderService) UpdatePaymentStatus(stripeEvent string, orderId uint64) error {
	tx := storage.GetInstance().BeginTransaction()
	var err error

	switch stripeEvent {
	case string(stripe.EventTypeCheckoutSessionCompleted):
		err = storage.GetInstance().Order.UpdatePaymentStatus(orderId, models.PaymentStatusCompleted, tx)
	case string(stripe.EventTypeCheckoutSessionExpired):
		err = s.handleRevertOrderItems(orderId, tx)
		if err != nil {
			break
		}
		err = storage.GetInstance().Order.UpdatePaymentStatus(orderId, models.PaymentStatusCancelled, tx)
	default:
		err = nil
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *OrderService) handleRevertOrderItems(orderId uint64, tx *gorm.DB) error {
	order, err := storage.GetInstance().Order.GetOrder(orderId, tx)

	if err != nil {
		return err
	}

	// Revert all the quantities
	for _, item := range order.OrderItems {
		if err = storage.GetInstance().Product.RevertProductQuantity(item.ProductId, item.Quantity, tx); err != nil {
			return err
		}
	}

	return nil
}

// convertToProtoOrder converts a model order to a protobuf order
func convertToProtoOrder(order *models.Order) *pb.Order {
	orderItems := make([]*pb.OrderItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		product, err := storage.GetInstance().Product.GetProduct(item.ProductId, nil)

		var productImage, productName string

		// Handle error if product not found;
		// Product can be deleted
		if err == nil {
			productName = product.Name

			if len(product.Images) > 0 {
				productImage = product.Images[0] // get only 1 image
			}
		}

		orderItems[i] = &pb.OrderItem{
			OrderId:      item.OrderId,
			ProductId:    item.ProductId,
			Quantity:     item.Quantity,
			Price:        item.Price,
			ProductName:  productName,
			ProductImage: productImage,
			CreatedAt:    item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    item.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &pb.Order{
		Id:                order.Id,
		UserId:            order.UserId,
		Total:             order.Total,
		Status:            string(order.Status),
		TransactionId:     order.TransactionId,
		CheckoutSessionId: order.CheckoutSessionId,
		PaymentStatus:     string(order.PaymentStatus),
		OrderItems:        orderItems,
		CreatedAt:         order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         order.UpdatedAt.Format(time.RFC3339),
	}
}
