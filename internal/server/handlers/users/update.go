package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" UpdateUserHandler çağrıldı")

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Geçersiz istek gövdesi", http.StatusBadRequest)
		log.Println(" JSON parse hatası:", err)
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		log.Println(" Kullanıcı bulunamadı:", err)
		return
	}

	user.Username = input.Username
	user.Email = input.Email
	user.Role = input.Role
	user.UpdatedAt = time.Now()

	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError)
		log.Println(" Güncelleme hatası:", err)
		return
	}

	response := UserResponse{
		ID:       uint(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(" JSON encode hatası:", err)
	}
}
