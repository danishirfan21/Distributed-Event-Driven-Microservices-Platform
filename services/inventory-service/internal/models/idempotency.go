package models

import (
	"time"
)

type ProcessedEvent struct {
	ID        string    `gorm:"primaryKey"`
	Service   string    `gorm:"primaryKey"`
	CreatedAt time.Time
}
