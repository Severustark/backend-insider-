package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteUserHandler tetiklendi")

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf(" Geçersiz ID: %v\n", idParam)
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}

	log.Printf(" Kullanıcı siliniyor: ID = %d\n", id)

	if err := db.DB.Where("user_id = ?", id).Delete(&models.Balance{}).Error; err != nil {
		log.Printf(" Balance silme hatası: %v\n", err)
		http.Error(w, "İlişkili bakiyeler silinemedi", http.StatusInternalServerError)
		return
	}

	result := db.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		log.Printf(" Kullanıcı silme hatası: %v\n", result.Error)
		http.Error(w, "Kullanıcı silinemedi", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf(" Kullanıcı bulunamadı: ID = %d\n", id)
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	log.Println("Kullanıcı başarıyla silindi")
	w.WriteHeader(http.StatusNoContent)
}
