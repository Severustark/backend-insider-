package routes

import (
	"log"
	"net/http"

	"github.com/Severustark/movietracker-backend/internal/server/handlers/auth"
	"github.com/Severustark/movietracker-backend/internal/server/handlers/balances"
	"github.com/Severustark/movietracker-backend/internal/server/handlers/transactions"
	"github.com/Severustark/movietracker-backend/internal/server/handlers/users"
	"github.com/Severustark/movietracker-backend/internal/server/middleware"
	"github.com/go-chi/chi/v5"
)

func ProtectedRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/protected/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected pong"))
		})

		r.Put("/test", func(w http.ResponseWriter, r *http.Request) {
			log.Println(" Test endpoint tetiklendi")
			w.Write([]byte(" Test endpoint çalıştı"))
		})
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/transactions/transfer", transactions.TransferHandler)
		})

		r.Get("/users", users.GetAllUsersHandler)
		r.Get("/users/{id}", users.GetUserByIDHandler)
		r.Put("/users/{id}", users.UpdateUserHandler)
		r.Delete("/users/{id}", users.DeleteUserHandler)
		r.Get("/balances/{user_id}", balances.GetBalanceHandler)
		r.Post("/balances/deposit", balances.DepositHandler)
		r.Post("/transactions", transactions.CreateTransactionHandler)
		r.Get("/transactions/user/{id}", transactions.GetUserTransactionsHandler)
		r.Get("/transactions", transactions.GetAllTransactionsHandler)
		r.Post("/refresh", auth.RefreshTokenHandler)

		r.Post("/transactions/credit", transactions.CreditTransactionHandler)
		r.Post("/transactions/debit", transactions.DebitTransactionHandler)
		r.Get("/api/v1/transactions", transactions.GetAllTransactionsHandler)
		r.Get("/balances/{id}", balances.GetBalanceHandler)
		r.Post("/balances/deposit", balances.DepositHandler)
		r.Get("/users/me", users.GetMeHandler)
		r.Put("/users/me", users.UpdateMeHandler)
		r.Post("/transactions/transfer", transactions.TransferHandler)
		r.Get("/balances/current", balances.GetCurrentBalanceHandler)
		r.Get("/balances/historical", balances.GetHistoricalBalanceHandler)
		r.Get("/transactions/{id}", transactions.GetTransactionByIDHandler)
		r.Get("/transactions/history", transactions.TransactionHistoryHandler)

	})
}
