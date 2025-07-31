package services

import (
	"context"
	"log"
	"sync"
	"time"
)

type Transaction struct {
	ID         int64
	FromUserID int64
	ToUserID   int64
	Amount     float64
}

type TransactionProcessor struct {
	queue      chan Transaction
	workerPool sync.WaitGroup
	workers    int
	service    *TransactionService
}

func NewTransactionProcessor(workers int, queueSize int, service *TransactionService) *TransactionProcessor {
	return &TransactionProcessor{
		queue:   make(chan Transaction, queueSize),
		workers: workers,
		service: service,
	}
}

func (tp *TransactionProcessor) Start() {
	log.Println("Starting workers...")
	for i := 0; i < tp.workers; i++ {
		tp.workerPool.Add(1)
		go tp.worker(i + 1)
	}
}

func (tp *TransactionProcessor) Stop() {
	close(tp.queue)
	tp.workerPool.Wait()
	log.Println("All workers stopped")
}

func (tp *TransactionProcessor) Submit(tx Transaction) {
	tp.queue <- tx
}

func (tp *TransactionProcessor) worker(id int) {
	defer tp.workerPool.Done()
	log.Printf("Worker %d started", id)

	for tx := range tp.queue {
		log.Printf("Worker %d processing transaction ID %d", id, tx.ID)
		err := tp.service.Transfer(context.Background(), tx.FromUserID, tx.ToUserID, tx.Amount)
		if err != nil {
			log.Printf("Transaction ID %d failed: %v", tx.ID, err)
		} else {
			log.Printf("Transaction ID %d processed successfully", tx.ID)
		}
		time.Sleep(500 * time.Millisecond) // işlem takibi için kısa bir bekleme
	}

	log.Printf("Worker %d stopped", id)
}
