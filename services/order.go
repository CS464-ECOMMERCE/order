package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"order/models"
	"order/storage"
)

// CartItem represents an item in the cart as stored in Redis
type CartItem struct {
	Id       uint64 `json:"id"`
	Quantity uint64 `json:"quantity"`
}

// Cart represents a user's shopping cart as stored in Redis
type Cart struct {
	SessionId string     `json:"session_id"`
	Items     []CartItem `json:"items"`
}

// OrderService provides operations for managing orders
type OrderService struct {
	redis         *RedisClient
	productClient *ProductClient
}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{
		redis:         GetRedisClient(),
		productClient: GetProductClient(),
	}
}

// PlaceOrder creates a new order from a user's cart
func (s *OrderService) PlaceOrder(sessionId string, userId uint64) (*models.Order, error) {
	// Get the user's cart
	cartData, err := s.redis.GetCart(sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	var cart Cart
	if err := json.Unmarshal(cartData, &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	// Validate cart is not empty
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Create the order
	order := &models.Order{
		UserId: userId,
		Status: "pending",
		Total:  0, // Will calculate as we process items
	}

	// Validate inventory and calculate total
	orderItems := make([]models.OrderItem, 0, len(cart.Items))

	// Define a transaction function to be executed with a Redis lock
	processFn := func() error {
		for _, item := range cart.Items {
			// Validate inventory
			product, err := s.productClient.ValidateInventory(item.Id, item.Quantity)
			if err != nil {
				return err
			}

			// Create order item
			orderItem := models.OrderItem{
				ProductId: item.Id,
				Quantity:  item.Quantity,
				Price:     product.Price,
			}
			orderItems = append(orderItems, orderItem)

			// Update total
			order.Total += product.Price * float32(item.Quantity)
		}

		// Save order to database
		order.OrderItems = orderItems
		savedOrder, err := storage.GetInstance().Order.CreateOrder(order)
		if err != nil {
			return fmt.Errorf("failed to save order: %w", err)
		}
		*order = *savedOrder

		// Update inventory for each product
		for _, item := range cart.Items {
			product, err := s.productClient.GetProduct(item.Id)
			if err != nil {
				return fmt.Errorf("failed to get product for inventory update: %w", err)
			}

			// Calculate new inventory
			newInventory := product.Inventory - item.Quantity
			if err := s.productClient.UpdateInventory(item.Id, newInventory); err != nil {
				return fmt.Errorf("failed to update inventory: %w", err)
			}
		}

		// Clear the cart
		if err := s.redis.DeleteCart(sessionId); err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		return nil
	}

	// Execute the transaction with a distributed lock
	err = s.redis.ExecuteWithLock(fmt.Sprintf("order:%s", sessionId), s.redis.config.RedisDefaultTTL, processFn)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(id uint64) (*models.Order, error) {
	return storage.GetInstance().Order.GetOrder(id)
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
	order, err := storage.GetInstance().Order.GetOrder(id)
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

	// Define a transaction function to be executed with a Redis lock
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
	return s.redis.ExecuteWithLock(fmt.Sprintf("cancel_order:%d", id), s.redis.config.RedisDefaultTTL, cancelFn)
}

// DeleteOrder deletes an order
func (s *OrderService) DeleteOrder(id uint64) error {
	// Get the order
	order, err := storage.GetInstance().Order.GetOrder(id)
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
