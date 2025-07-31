package balances

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

type DepositRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("💰 DepositHandler tetiklendi")

	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Kullanıcı kimliği alınamadı", http.StatusUnauthorized)
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "Geçersiz tutar", http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Veritabanı işlemi başlatılamadı", http.StatusInternalServerError)
		return
	}

	// Balance yoksa oluştur, varsa güncelle
	var balance models.Balance
	if err := tx.Where("user_id = ?", userID).FirstOrCreate(&balance, models.Balance{
		UserID:        uint(userID),
		Amount:        0,
		LastUpdatedAt: time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Bakiye oluşturulamadı", http.StatusInternalServerError)
		return
	}

	balance.Amount += req.Amount
	balance.LastUpdatedAt = time.Now()

	if err := tx.Save(&balance).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Bakiye güncellenemedi", http.StatusInternalServerError)
		return
	}

	// Transaction kaydı
	if err := tx.Exec(`
		INSERT INTO transactions (to_user_id, amount, type, status, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		userID, req.Amount, "credit", "completed", req.Description, time.Now(),
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
	w.Write([]byte(`{"message": "Para yatırma işlemi başarıyla tamamlandı"}`))
}
