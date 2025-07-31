package models

import (
	"time"
)

type Balance struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	LastUpdatedAt time.Time `gorm:"autoUpdateTime" json:"last_updated_at"`
}
