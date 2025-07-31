package balances

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBalanceHandler tetiklendi")

	idParam := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID <= 0 {
		http.Error(w, "Geçersiz kullanıcı ID", http.StatusBadRequest)
		return
	}

	var balance models.Balance
	if err := db.DB.Where("user_id = ?", userID).First(&balance).Error; err != nil {
		http.Error(w, "Bakiye bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
