package models

import (
	"time"
)

type OrderStatus string
type PaymentStatus string

var (
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"

	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Order represents a user's order
type Order struct {
	Id                uint64        `json:"id" gorm:"primaryKey"`
	UserId            uint64        `json:"user_id"`
	Total             float32       `json:"total"`
	Status            OrderStatus   `json:"status" gorm:"gorm:type:order_status;default:processing"`
	TransactionId     string        `json:"transaction_id"`
	CheckoutSessionId string        `json:"checkout_session_id"`
	PaymentStatus     PaymentStatus `json:"payment_status" gorm:"type:payment_status;default:pending"`
	OrderItems        []OrderItem   `json:"order_items" gorm:"foreignKey:OrderId"`
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	OrderId   uint64    `json:"order_id" gorm:"primaryKey"`
	ProductId uint64    `json:"product_id" gorm:"primaryKey"`
	Quantity  uint64    `json:"quantity"`
	Price     float32   `json:"price"` // Storing price at time of purchase
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
