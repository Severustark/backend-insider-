package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
	"golang.org/x/crypto/bcrypt"
)

type UpdateMeRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UpdateMeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("✅ UpdateMeHandler tetiklendi")

	// JWT'den kullanıcı ID'sini al
	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var req UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Şifre hashlenemedi", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashed)
	}

	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Kullanıcı güncellenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bilgiler başarıyla güncellendi"})
}
