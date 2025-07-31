package services

import (
	"errors"

	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) error
	Authenticate(username, password string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Register: kullanıcı oluştururken şifreyi hash'ler
func (s *userService) Register(user *models.User) error {
	// Şifre hash'leme
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPass)

	// Kullanıcıyı veritabanına ekle
	return s.userRepo.Create(user)
}

// Authenticate: kullanıcı giriş kontrolü
func (s *userService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("username or password is incorrect")
	}

	// Şifre doğrulama
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("username or password is incorrect")
	}

	return user, nil
}
