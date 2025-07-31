package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransactionService struct {
	DB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{DB: db}
}

func (s *TransactionService) Transfer(ctx context.Context, fromUserID, toUserID int64, amount float64) error {
	if fromUserID == toUserID {
		return errors.New("cannot transfer to the same account")
	}

	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var senderBalance float64
		if err := tx.Raw(`SELECT amount FROM balances WHERE user_id = ? FOR UPDATE`, fromUserID).Scan(&senderBalance).Error; err != nil {
			return fmt.Errorf("failed to fetch sender balance: %w", err)
		}

		if senderBalance < amount {
			return errors.New("insufficient balance")
		}

		if err := tx.Exec(`UPDATE balances SET amount = amount - ?, last_updated_at = NOW() WHERE user_id = ?`, amount, fromUserID).Error; err != nil {
			return fmt.Errorf("failed to deduct balance: %w", err)
		}

		if err := tx.Exec(`UPDATE balances SET amount = amount + ?, last_updated_at = NOW() WHERE user_id = ?`, amount, toUserID).Error; err != nil {
			return fmt.Errorf("failed to add balance: %w", err)
		}

		if err := tx.Exec(`
			INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at)
			VALUES (?, ?, ?, 'transfer', 'completed', ?)`,
			fromUserID, toUserID, amount, time.Now(),
		).Error; err != nil {
			return fmt.Errorf("failed to insert transaction: %w", err)
		}

		return nil
	})
}
