package transactions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

type DebitRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func DebitTransactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" DebitTransactionHandler tetiklendi")

	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var req DebitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek formatı", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Tutar pozitif olmalı", http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Veritabanı işlemi başlatılamadı", http.StatusInternalServerError)
		return
	}

	var balance float64
	err := tx.Raw(`SELECT amount FROM balances WHERE user_id = ? FOR UPDATE`, userID).Scan(&balance).Error
	if err != nil {
		tx.Rollback()
		http.Error(w, "Bakiye okunamadı", http.StatusInternalServerError)
		return
	}

	if balance < req.Amount {
		tx.Rollback()
		http.Error(w, "Yetersiz bakiye", http.StatusBadRequest)
		return
	}

	//  Bakiyeden düş
	if err := tx.Exec(`UPDATE balances SET amount = amount - ?, last_updated_at = ? WHERE user_id = ?`,
		req.Amount, time.Now(), userID).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Bakiye güncellenemedi", http.StatusInternalServerError)
		return
	}

	if err := tx.Exec(`
		INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, description, created_at)
		VALUES (?, NULL, ?, ?, ?, ?, ?)`,
		userID, req.Amount, "debit", "completed", req.Description, time.Now(),
	).Error; err != nil {
		tx.Rollback()
		http.Error(w, "İşlem kaydedilemedi", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		http.Error(w, "Commit başarısız", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Bakiyeden düşme işlemi başarıyla tamamlandı"}`))
}
