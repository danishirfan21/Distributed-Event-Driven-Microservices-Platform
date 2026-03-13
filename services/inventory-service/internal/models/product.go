package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	Description   string         `json:"description"`
	Price         float64        `gorm:"not null" json:"price"`
	StockQuantity int32          `gorm:"not null" json:"stock_quantity"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type StockHistory struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProductID string    `gorm:"not null;index" json:"product_id"`
	Change    int32     `gorm:"not null" json:"change"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
