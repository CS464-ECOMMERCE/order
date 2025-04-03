package storage

import (
	"errors"
	"order/models"

	"github.com/stripe/stripe-go/v81"
	"gorm.io/gorm"
)

// PaymentInterface defines the operations for payment storage
type PaymentInterface interface {
	CreatePayment(payment *models.Payment) (*models.Payment, error)
	GetPayment(orderId uint64) (*models.Payment, error)
	UpdatePaymentStatusByStripe(orderId uint64, stripeStatus string) error
	DeletePayment(orderId uint64) error
}

// PaymentDB implements PaymentInterface
type PaymentDB struct {
	read  *gorm.DB
	write *gorm.DB
}

// NewPaymentTable creates a new PaymentDB instance
func NewPaymentTable(read, write *gorm.DB) PaymentInterface {
	// Auto migrate the payment table
	if err := read.AutoMigrate(&models.Payment{}); err != nil {
		panic(err)
	}
	return &PaymentDB{
		read:  read,
		write: write,
	}
}

// CreatePayment implements PaymentInterface.
func (db *PaymentDB) CreatePayment(payment *models.Payment) (*models.Payment, error) {
	tx := db.write.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create payment
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPayment implements PaymentInterface.
func (db *PaymentDB) GetPayment(orderId uint64) (*models.Payment, error) {
	var payment models.Payment
	if err := db.read.Where("order_id = ?", orderId).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &payment, nil
}

// UpdatePaymentStatus implements PaymentInterface.
func (db *PaymentDB) UpdatePaymentStatusByStripe(orderId uint64, stripeStatus string) error {
	tx := db.write.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var status models.PaymentStatus

	switch stripeStatus {
	case string(stripe.CheckoutSessionPaymentStatusPaid):
		status = models.PaymentStatusCompleted
	case string(stripe.CheckoutSessionPaymentStatusUnpaid):
		status = models.PaymentStatusPending
	case string(stripe.CheckoutSessionStatusExpired):
		status = models.PaymentStatusCancelled
	default:
		status = models.PaymentStatusError
	}

	// Update payment status
	if err := tx.Model(&models.Payment{}).Where("order_id = ?", orderId).Update("status", status).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// DeletePayment implements PaymentInterface.
func (db *PaymentDB) DeletePayment(orderId uint64) error {
	tx := db.write.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Delete payment
	if err := tx.Where("order_id = ?", orderId).Delete(&models.Payment{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
