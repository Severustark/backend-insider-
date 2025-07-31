package models

import (
	"errors"
	"regexp"
	"time"
)

// User struct
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `gorm:"-" json:"-"` // sadece uygulama içinde kullanılır, DB ve JSON'da yok sayılır
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate validates the User fields
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	if u.Role == "" {
		return errors.New("role is required")
	}

	return nil
}
