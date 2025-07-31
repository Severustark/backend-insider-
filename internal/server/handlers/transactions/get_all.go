package transactions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

type TransactionResponse struct {
	ID          uint    `json:"id"`
	FromUserID  *uint   `json:"from_user_id,omitempty"`
	ToUserID    *uint   `json:"to_user_id,omitempty"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

func GetAllTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" GetAllTransactionsHandler tetiklendi")

	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var transactions []TransactionResponse
	err := db.DB.Raw(`
		SELECT id, from_user_id, to_user_id, amount, type, status, description, created_at
		FROM transactions
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
	`, userID, userID).Scan(&transactions).Error

	if err != nil {
		http.Error(w, "İşlem geçmişi alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
