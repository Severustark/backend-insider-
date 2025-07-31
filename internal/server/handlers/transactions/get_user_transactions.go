package transactions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

func GetUserTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" GetUserTransactionsHandler tetiklendi")

	idParam := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf(" Geçersiz kullanıcı ID: %v\n", idParam)
		http.Error(w, "Geçersiz kullanıcı ID", http.StatusBadRequest)
		return
	}

	var transactions []models.Transaction
	if err := db.DB.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at desc").
		Find(&transactions).Error; err != nil {
		log.Printf(" İşlemler alınamadı: %v\n", err)
		http.Error(w, "İşlemler alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
