package repositories

import (
	"database/sql"
	"errors"

	"github.com/Severustark/movietracker-backend/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *userRepo) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id=$1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE username=$1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *userRepo) Update(user *models.User) error {
	query := `
		UPDATE users SET username=$1, email=$2, password_hash=$3, role=$4, updated_at=NOW()
		WHERE id=$5
	`
	res, err := r.db.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.ID,
	)
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

func (r *userRepo) Delete(id int) error {
	query := `DELETE FROM users WHERE id=$1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
