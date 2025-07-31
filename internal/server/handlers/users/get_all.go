package users

import (
	"encoding/json"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		http.Error(w, "Kullanıcılar alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
