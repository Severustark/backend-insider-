package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	customMiddleware "github.com/Severustark/movietracker-backend/internal/server/middleware"
	"github.com/Severustark/movietracker-backend/internal/server/routes"
)

func StartHTTPServer(addr string) {
	r := chi.NewRouter()

	// Built-in Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Custom Middlewares
	r.Use(customMiddleware.CORSMiddleware)
	r.Use(customMiddleware.SecurityHeaders)
	r.Use(customMiddleware.RateLimiter)

	// Route'lar
	r.Route("/api/v1", func(r chi.Router) {
		routes.AuthRoutes(r)
		routes.ProtectedRoutes(r)
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Printf(" HTTP Server listening on %s\n", addr)
	http.ListenAndServe(addr, r)
}
