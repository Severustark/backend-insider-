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

func TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" TransactionHistoryHandler tetiklendi")

	userIDRaw := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Geçersiz kullanıcı kimliği", http.StatusUnauthorized)
		return
	}

	// Tarih aralığı parse et
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var startTime, endTime time.Time
	var err error

	if startStr != "" {
		startTime, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, "Geçersiz start tarihi (format: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}

	if endStr != "" {
		endTime, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			http.Error(w, "Geçersiz end tarihi (format: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}

	var transactions []models.Transaction
	query := db.DB.Where("from_user_id = ? OR to_user_id = ?", userID, userID)

	if !startTime.IsZero() && !endTime.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	} else if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	} else if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	if err := query.Order("created_at desc").Find(&transactions).Error; err != nil {
		http.Error(w, "İşlem geçmişi alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
