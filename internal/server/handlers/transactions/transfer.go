package transactions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

type TransferRequest struct {
	ToUserID    int     `json:"to_user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" TransferHandler tetiklendi")

	// 1. Gönderen ID
	fromUserRaw := r.Context().Value(middleware.UserIDKey)
	fromUserID, ok := fromUserRaw.(int)
	if !ok || fromUserID == 0 {
		http.Error(w, "Kimlik doğrulama başarısız", http.StatusUnauthorized)
		return
	}

	// 2. İstek çözümlemesi
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Tutar pozitif olmalı", http.StatusBadRequest)
		return
	}
	if req.ToUserID == fromUserID {
		http.Error(w, "Kendinize para gönderemezsiniz", http.StatusBadRequest)
		return
	}

	// 3. Transaction başlat
	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Veritabanı işlemi başlatılamadı", http.StatusInternalServerError)
		return
	}

	// 4. Gönderen bakiye kontrol
	var senderBalance float64
	if err := tx.Raw(`SELECT amount FROM balances WHERE user_id = ? FOR UPDATE`, fromUserID).Scan(&senderBalance).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Gönderen bakiyesi okunamadı", http.StatusInternalServerError)
		return
	}
	if senderBalance < req.Amount {
		tx.Rollback()
		http.Error(w, "Yetersiz bakiye", http.StatusBadRequest)
		return
	}

	// 5. Alıcı bakiye kontrol
	var exists int64
	if err := tx.Raw(`SELECT COUNT(*) FROM users WHERE id = ?`, req.ToUserID).Scan(&exists).Error; err != nil || exists == 0 {
		tx.Rollback()
		http.Error(w, "Alıcı bulunamadı", http.StatusNotFound)
		return
	}

	// 6. Gönderen bakiyesini azalt
	if err := tx.Exec(`UPDATE balances SET amount = amount - ?, last_updated_at = ? WHERE user_id = ?`,
		req.Amount, time.Now(), fromUserID).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Gönderen bakiyesi güncellenemedi", http.StatusInternalServerError)
		return
	}

	// 7. Alıcı bakiyesini artır
	if err := tx.Exec(`UPDATE balances SET amount = amount + ?, last_updated_at = ? WHERE user_id = ?`,
		req.Amount, time.Now(), req.ToUserID).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Alıcı bakiyesi güncellenemedi", http.StatusInternalServerError)
		return
	}

	// 8. Her iki kullanıcı için işlem kaydı
	if err := tx.Exec(`
		INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, description, created_at)
		VALUES (?, ?, ?, 'transfer', 'completed', ?, ?)`,
		fromUserID, req.ToUserID, req.Amount, req.Description, time.Now()).Error; err != nil {
		tx.Rollback()
		http.Error(w, "İşlem kaydı başarısız", http.StatusInternalServerError)
		return
	}

	// 9. Commit
	if err := tx.Commit().Error; err != nil {
		http.Error(w, "İşlem commit edilemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Transfer başarıyla tamamlandı"}`))
}
