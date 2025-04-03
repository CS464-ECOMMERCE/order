package models

import "time"

type PaymentStatus string

var (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusError     PaymentStatus = "error"
)

type Payment struct {
	OrderId       uint64        `json:"order_id" gorm:"primaryKey;not null"` // PK and FK
	Order         Order         `gorm:"foreignKey:OrderId;references:ID;constraint:OnDelete:CASCADE"`
	TransactionId string        `json:"transaction_id"`
	Status        PaymentStatus `json:"status" gorm:"default:pending"`
	Gateway       string        `json:"gateway"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}
