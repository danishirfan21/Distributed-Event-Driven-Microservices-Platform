package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID                    string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID               string         `gorm:"not null;index" json:"order_id"`
	Amount                float64        `gorm:"not null" json:"amount"`
	Status                string         `gorm:"not null" json:"status"` // e.g., CONFIRMED, FAILED
	ProviderTransactionID string         `json:"provider_transaction_id"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}
