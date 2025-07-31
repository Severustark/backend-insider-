package transactions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
)

type CreateTransactionRequest struct {
	FromUserID uint    `json:"from_user_id"`
	ToUserID   uint    `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" CreateTransactionHandler tetiklendi")

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf(" Geçersiz istek: %v", err)
		http.Error(w, "Geçersiz istek verisi", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Tutar pozitif olmalı", http.StatusBadRequest)
		return
	}

	var fromBalance models.Balance
	if err := db.DB.Where("user_id = ?", req.FromUserID).First(&fromBalance).Error; err != nil {
		log.Printf(" Gönderici bakiyesi bulunamadı: %v", err)
		http.Error(w, "Gönderici bakiyesi bulunamadı", http.StatusNotFound)
		return
	}

	if fromBalance.Amount < req.Amount {
		http.Error(w, "Yetersiz bakiye", http.StatusBadRequest)
		return
	}

	var toBalance models.Balance
	if err := db.DB.Where("user_id = ?", req.ToUserID).First(&toBalance).Error; err != nil {
		log.Printf(" Alıcı için yeni bakiye oluşturuluyor")
		toBalance = models.Balance{
			UserID: req.ToUserID,
			Amount: 0,
		}
		db.DB.Create(&toBalance)
	}

	// İşlemi gerçekleştirme
	fromBalance.Amount -= req.Amount
	toBalance.Amount += req.Amount

	tx := db.DB.Begin()
	if err := tx.Save(&fromBalance).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Gönderici bakiyesi güncellenemedi", http.StatusInternalServerError)
		return
	}
	if err := tx.Save(&toBalance).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Alıcı bakiyesi güncellenemedi", http.StatusInternalServerError)
		return
	}

	transaction := models.Transaction{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Amount:     req.Amount,
		Type:       "transfer",
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		http.Error(w, "İşlem kaydedilemedi", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	log.Println(" Transfer işlemi başarılı")
	json.NewEncoder(w).Encode(transaction)
}
