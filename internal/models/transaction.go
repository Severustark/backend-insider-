package models

import (
	"errors"
	"time"
)

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FromUserID  uint      `json:"from_user_id"`
	ToUserID    uint      `json:"to_user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`   // e.g. "credit", "debit", "transfer"
	Status      string    `json:"status"` // e.g. "pending", "completed", "failed"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if t.Type == "" {
		return errors.New("transaction type is required")
	}
	if t.Status == "" {
		return errors.New("transaction status is required")
	}
	return nil
}
