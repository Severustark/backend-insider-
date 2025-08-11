package routes

import (
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/server/handlers/auth"
	"github.com/go-chi/chi/v5"
)

func PublicRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Post("/api/v1/auth/login", auth.LoginHandler)

	r.Post("/api/v1/auth/refresh", auth.RefreshTokenHandler)

	return r
}
