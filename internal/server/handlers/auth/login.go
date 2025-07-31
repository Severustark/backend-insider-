package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Geçersiz şifre", http.StatusUnauthorized)
		return
	}

	accessToken, _ := utils.GenerateToken(int64(user.ID), user.Role)
	refreshToken, _ := utils.GenerateRefreshToken(int64(user.ID), user.Role)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Message:      "Giriş başarılı",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
