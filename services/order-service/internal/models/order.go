package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusCompleted OrderStatus = "COMPLETED"
)

type Order struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     string         `gorm:"not null;index" json:"user_id"`
	TotalPrice float64        `gorm:"not null" json:"total_price"`
	Status     OrderStatus    `gorm:"not null;default:'PENDING'" json:"status"`
	Items      []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID   string  `gorm:"not null;index" json:"order_id"`
	ProductID string  `gorm:"not null" json:"product_id"`
	Quantity  int32   `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"unit_price"`
}
