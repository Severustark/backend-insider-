package repositories

import (
	"database/sql"
	"errors"

	"github.com/Severustark/movietracker-backend/internal/models"
)

type BalanceRepository interface {
	GetByUserID(userID int) (*models.Balance, error)
	Update(balance *models.Balance) error
	Create(balance *models.Balance) error
}

type balanceRepo struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) BalanceRepository {
	return &balanceRepo{db: db}
}

func (r *balanceRepo) GetByUserID(userID int) (*models.Balance, error) {
	balance := &models.Balance{}
	query := `SELECT user_id, amount, last_updated_at FROM balances WHERE user_id=$1`
	err := r.db.QueryRow(query, userID).Scan(&balance.UserID, &balance.Amount, &balance.LastUpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("balance not found")
	}
	return balance, err
}

func (r *balanceRepo) Update(balance *models.Balance) error {
	query := `
		UPDATE balances SET amount=$1, last_updated_at=NOW() WHERE user_id=$2
	`
	res, err := r.db.Exec(query, balance.Amount, balance.UserID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *balanceRepo) Create(balance *models.Balance) error {
	query := `
		INSERT INTO balances (user_id, amount, last_updated_at) VALUES ($1, $2, NOW())
	`
	_, err := r.db.Exec(query, balance.UserID, balance.Amount)
	return err
}
