package routes

import (
	"github.com/Severustark/movietracker-backend/internal/server/handlers/auth"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginHandler)
		r.Post("/refresh", auth.RefreshTokenHandler)
	})
}
