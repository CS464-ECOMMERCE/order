package services

import (
	"errors"
	"order/models"
	"order/storage"

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
func (s *OrderService) GetOrder(id uint64) (*models.Order, error) {
	return storage.GetInstance().Order.GetOrder(id, nil)
}

// GetOrdersByUser retrieves all orders for a user
func (s *OrderService) GetOrdersByUser(userId uint64) ([]*models.Order, error) {
	return storage.GetInstance().Order.GetOrdersByUserId(userId)
}

// GetOrdersByMerchant retrieves all orders containing products from a specific merchant
func (s *OrderService) GetOrdersByMerchant(merchantId uint64) ([]*models.Order, error) {
	return storage.GetInstance().Order.GetOrdersByMerchantId(merchantId)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id uint64, status string) error {
	// Validate status
	if status != "pending" && status != "completed" && status != "cancelled" {
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
	if order.Status == "cancelled" {
		return errors.New("order is already cancelled")
	}
	if order.Status == "completed" {
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
