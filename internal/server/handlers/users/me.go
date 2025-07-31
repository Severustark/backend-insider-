package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" GetMeHandler tetiklendi")

	// Context’ten kullanıcı ID’sini al
	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	// Hassas bilgileri döndürmüyoruz (örn. şifre)
	type response struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	res := response{
		ID:       uint(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
