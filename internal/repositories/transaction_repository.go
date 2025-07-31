package repositories

import (
	"database/sql"
	"errors"

	"github.com/Severustark/movietracker-backend/internal/models"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	GetByID(id int) (*models.Transaction, error)
	GetHistoryByUser(userID int) ([]*models.Transaction, error)
}

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(query,
		tx.FromUserID,
		tx.ToUserID,
		tx.Amount,
		tx.Type,
		tx.Status,
	).Scan(&tx.ID, &tx.CreatedAt)
}

func (r *transactionRepo) GetByID(id int) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `
		SELECT id, from_user_id, to_user_id, amount, type, status, created_at
		FROM transactions
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.FromUserID,
		&tx.ToUserID,
		&tx.Amount,
		&tx.Type,
		&tx.Status,
		&tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}
	return tx, err
}

func (r *transactionRepo) GetHistoryByUser(userID int) ([]*models.Transaction, error) {
	query := `
		SELECT id, from_user_id, to_user_id, amount, type, status, created_at
		FROM transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.FromUserID,
			&tx.ToUserID,
			&tx.Amount,
			&tx.Type,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
