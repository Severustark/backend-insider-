package routes

import (
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/server/handlers/auth"
	"github.com/go-chi/chi/v5"
)

func PublicRoutes() http.Handler {
	r := chi.NewRouter()

	// ✅ Sağlık kontrolü (isteğe bağlı)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// ✅ Login route'u burada tanımlanmalı:
	r.Post("/api/v1/auth/login", auth.LoginHandler)

	// Refresh token endpoint'i varsa onu da burada tanımla:
	r.Post("/api/v1/auth/refresh", auth.RefreshTokenHandler)

	return r
}
