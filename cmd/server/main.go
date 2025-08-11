package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Severustark/movietracker-backend/internal/db"
	"github.com/Severustark/movietracker-backend/internal/models"
	"github.com/Severustark/movietracker-backend/internal/server/routes"
	"github.com/Severustark/movietracker-backend/internal/services"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env dosyası yüklenemedi, varsayılan değerler kullanılacak")
	}

	// Veritabanını başlat
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	// Örnek transaction sayısını logla
	var count int64
	db.DB.Model(&models.Transaction{}).Where("from_user_id = ? OR to_user_id = ?", 1, 1).Count(&count)
	log.Printf("user_id=1 için transaction kayıt sayısı: %d\n", count)

	// Transaction servisini başlat
	transactionService := services.NewTransactionService(db.DB)
	tp := services.NewTransactionProcessor(3, 10, transactionService)
	tp.Start()

	// İşlem gönder
	go func() {
		transactions := []services.Transaction{
			{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100.00},
			{ID: 2, FromUserID: 1, ToUserID: 2, Amount: 200.00},
			{ID: 3, FromUserID: 2, ToUserID: 1, Amount: 300.00},
			{ID: 4, FromUserID: 2, ToUserID: 1, Amount: 400.00},
			{ID: 5, FromUserID: 1, ToUserID: 2, Amount: 500.00},
		}

		for _, tx := range transactions {
			tp.Submit(tx)
		}

		time.Sleep(3 * time.Second)
		tp.Stop()
		log.Println("All transactions processed")
	}()

	// HTTP router'ı al
	r := routes.PublicRoutes()

	log.Println("HTTP sunucusu http://localhost:8080 üzerinde çalışıyor...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Sunucu başlatılamadı: ", err)
	}
}
