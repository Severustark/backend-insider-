package balances

import (
	"encoding/json"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

func GetCurrentBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var balance models.Balance
	if err := db.DB.Where("user_id = ?", userID).First(&balance).Error; err != nil {
		http.Error(w, "Bakiye bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
