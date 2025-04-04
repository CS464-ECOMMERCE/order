package services

import (
	"errors"
	"fmt"
	"order/models"
	"order/storage"
	"time"

	"github.com/stripe/stripe-go/v81"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// CartItem represents an item in the cart as stored in cartClient
type CartItem struct {
	Id       uint64 `json:"id"`
	Quantity uint64 `json:"quantity"`
}

// Cart represents a user's shopping cart as stored in cartClient
type Cart struct {
	SessionId string     `json:"session_id"`
	Items     []CartItem `json:"items"`
}

// OrderService provides operations for managing orders
type OrderService struct {
	productClient *ProductService
	cartClient    *CartService
	redis         *RedisClient
}

// NewOrderService creates a new order service
func NewOrderService(productConn, cartConn *grpc.ClientConn) *OrderService {
	return &OrderService{
		productClient: NewProductService(productConn),
		cartClient:    NewCartService(cartConn),
		redis:         GetRedisClient(),
	}
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
	// Get the order
	order, err := storage.GetInstance().Order.GetOrder(id, nil)
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

	// Define a transaction function to be executed with a cartClient lock
	cancelFn := func() error {
		// Restore inventory for each product
		for _, item := range order.OrderItems {
			product, err := s.productClient.GetProduct(item.ProductId)
			if err != nil {
				return fmt.Errorf("failed to get product for inventory restore: %w", err)
			}

			// Calculate new inventory
			newInventory := product.Inventory + item.Quantity
			if err := s.productClient.UpdateInventory(item.ProductId, newInventory); err != nil {
				return fmt.Errorf("failed to restore inventory: %w", err)
			}
		}

		// Update order status to cancelled
		if err := storage.GetInstance().Order.UpdateOrderStatus(id, "cancelled"); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	}

	// Execute the transaction with a distributed lock
	return s.redis.ExecuteWithLock(fmt.Sprintf("cancel_order:%d", id), 10*time.Second, cancelFn)
}

// DeleteOrder deletes an order
func (s *OrderService) DeleteOrder(id uint64) error {
	// Get the order
	order, err := storage.GetInstance().Order.GetOrder(id, nil)
	if err != nil {
		return err
	}

	// If order is not cancelled, cancel it first to restore inventory
	if order.Status != "cancelled" {
		if err := s.CancelOrder(id); err != nil {
			return fmt.Errorf("failed to cancel order before deletion: %w", err)
		}
	}

	return storage.GetInstance().Order.DeleteOrder(id)
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
