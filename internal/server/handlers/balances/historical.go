package balances

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

type HistoricalBalance struct {
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"` // credit or debit
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}

func GetHistoricalBalanceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" GetHistoricalBalanceHandler tetiklendi")

	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var history []HistoricalBalance
	if err := db.DB.
		Table("transactions").
		Select("amount, type, status, created_at as occurred_at").
		Where("to_user_id = ? OR from_user_id = ?", userID, userID).
		Order("created_at desc").
		Scan(&history).Error; err != nil {
		http.Error(w, "Geçmiş bakiye getirilemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
