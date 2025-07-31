package transactions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
)

// TransferRequest struct'ı transfer.go'dan kullanılıyor

func TransferPendingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" TransferPendingHandler tetiklendi")

	fromUserRaw := r.Context().Value(middleware.UserIDKey)
	fromUserID, ok := fromUserRaw.(int)
	if !ok || fromUserID == 0 {
		http.Error(w, "Kimlik doğrulama başarısız", http.StatusUnauthorized)
		return
	}

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

	// Sadece "pending" kayıt açıyoruz, worker sonra işleyip sonuca göre update edecek
	tx := models.Transaction{
		FromUserID:  uint(fromUserID),
		ToUserID:    uint(req.ToUserID),
		Amount:      req.Amount,
		Type:        "transfer",
		Status:      "pending",
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&tx).Error; err != nil {
		log.Println(" Pending kayıt oluşturulamadı:", err)
		http.Error(w, "İşlem başlatılamadı", http.StatusInternalServerError)
		return
	}

	log.Printf(" Transfer talebi kaydedildi. ID: %d\n", tx.ID)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "Transfer talebiniz alınmıştır. Kısa süre içinde işlenecektir."}`))
}
