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

func GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTransactionByIDHandler tetiklendi")

	idParam := chi.URLParam(r, "id")
	transactionID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf(" Geçersiz işlem ID: %v\n", idParam)
		http.Error(w, "Geçersiz işlem ID", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	if err := db.DB.First(&transaction, transactionID).Error; err != nil {
		log.Printf(" İşlem bulunamadı: %v\n", err)
		http.Error(w, "İşlem bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}
