package auth

import (
	"encoding/json"
	"net/http"

	"os"

	"github.com/Severustark/movietracker-backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Geçersiz refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Token claim çözülemedi", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok1 := claims["user_id"].(float64)
	role, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		http.Error(w, "Eksik token verisi", http.StatusUnauthorized)
		return
	}
	userID := int64(userIDFloat)

	// Yeni access token üret
	newAccessToken, err := utils.GenerateToken(userID, role)
	if err != nil {
		http.Error(w, "Token üretilemedi", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(RefreshResponse{
		AccessToken: newAccessToken,
	})
}
