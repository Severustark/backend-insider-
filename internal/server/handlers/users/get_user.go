package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"

	"github.com/go-chi/chi/v5"
)

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
