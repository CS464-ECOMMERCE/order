package storage

import (
	"errors"
	"order/models"

	"gorm.io/gorm"
)

// OrderInterface defines the operations for order storage
type OrderInterface interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	GetOrder(id uint64, tx *gorm.DB) (*models.Order, error)
	GetOrdersByUserId(userId uint64) ([]*models.Order, error)
	GetOrdersByMerchantId(merchantId uint64) ([]*models.Order, error)
	UpdateOrderStatus(id uint64, status string) error
	DeleteOrder(id uint64) error
	UpdatePaymentStatus(id uint64, status models.PaymentStatus, tx *gorm.DB) error
}

// OrderDB implements OrderInterface
type OrderDB struct {
	read  *gorm.DB
	write *gorm.DB
}

// NewOrderTable creates a new OrderDB instance
func NewOrderTable(read, write *gorm.DB) OrderInterface {
	// Auto migrate the order and order_item tables
	if err := read.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
		panic(err)
	}
	return &OrderDB{
		read:  read,
		write: write,
	}
}

// CreateOrder creates a new order
func (o *OrderDB) CreateOrder(order *models.Order) (*models.Order, error) {
	tx := o.write.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (o *OrderDB) GetOrder(id uint64, tx *gorm.DB) (*models.Order, error) {
	var order models.Order

	db := tx
	if db == nil {
		db = o.read
	}

	if err := db.Preload("OrderItems").Where("id = ?", id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// GetOrdersByUserId retrieves all orders for a user
func (o *OrderDB) GetOrdersByUserId(userId uint64) ([]*models.Order, error) {
	var orders []*models.Order

	if err := o.read.Preload("OrderItems").Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrdersByMerchantId retrieves all orders for a merchant
func (o *OrderDB) GetOrdersByMerchantId(merchantId uint64) ([]*models.Order, error) {
	var orders []*models.Order

	// Join orders with order_items, then with products to filter by merchant_id
	if err := o.read.
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Where("products.merchant_id = ?", merchantId).
		Preload("OrderItems").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrderStatus updates the status of an order
func (o *OrderDB) UpdateOrderStatus(id uint64, status string) error {
	result := o.write.Model(&models.Order{}).Where("id = ?", id).Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

// DeleteOrder deletes an order
func (o *OrderDB) DeleteOrder(id uint64) error {
	tx := o.write.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Delete order items first
	if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete order
	if err := tx.Where("id = ?", id).Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// UpdatePaymentStatus updates the payment status of an order
func (o *OrderDB) UpdatePaymentStatus(id uint64, status models.PaymentStatus, tx *gorm.DB) error {
	db := tx
	if db == nil {
		db = o.write
	}

	result := db.Model(&models.Order{}).Where("id = ?", id).Update("payment_status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}
