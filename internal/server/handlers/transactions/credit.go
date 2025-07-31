package transactions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
)

type CreditRequest struct {
	ToUserID    uint    `json:"to_user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func CreditTransactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" CreditTransactionHandler tetiklendi")

	var req CreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek verisi", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 || req.ToUserID == 0 {
		http.Error(w, "Geçersiz kullanıcı veya tutar", http.StatusBadRequest)
		return
	}

	//  Transaction başlat
	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Veritabanı transaction başlatılamadı", http.StatusInternalServerError)
		return
	}

	//  Kullanıcıya ait bakiye var mı?
	var count int64
	err := tx.Table("balances").Where("user_id = ?", req.ToUserID).Count(&count).Error
	if err != nil {
		tx.Rollback()
		http.Error(w, "Bakiye kontrolü başarısız", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		// Yeni bakiye oluştur
		if err := tx.Exec(`
			INSERT INTO balances (user_id, amount, last_updated_at) VALUES (?, ?, ?)
		`, req.ToUserID, req.Amount, time.Now()).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Yeni bakiye oluşturulamadı", http.StatusInternalServerError)
			return
		}
	} else {
		// Var olan bakiyeye ekle
		if err := tx.Exec(`
			UPDATE balances SET amount = amount + ?, last_updated_at = ? WHERE user_id = ?
		`, req.Amount, time.Now(), req.ToUserID).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Bakiye güncellenemedi", http.StatusInternalServerError)
			return
		}
	}

	//  İşlemi transaction tablosuna yaz
	if err := tx.Exec(`
	INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, description, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
`, nil, req.ToUserID, req.Amount, "credit", "completed", req.Description, time.Now()).Error; err != nil {
		log.Printf(" Transaction kaydı eklenemedi: %v", err) // << burası önemli
		tx.Rollback()
		http.Error(w, "İşlem kaydedilemedi", http.StatusInternalServerError)
		return
	}

	//  Commit
	if err := tx.Commit().Error; err != nil {
		http.Error(w, "Transaction commit edilemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Kredi işlemi başarıyla tamamlandı"}`))
}
