package service

import (
	"log"
	//"net/http"
	"sync"
	"time"

	stock "github.com/alireza-aliabadi/golang-stocking-update/internal/models"
)

type StoreUpdater struct {}

func NewStoreUpdater() *StoreUpdater {
	return &StoreUpdater{}
}

func  (s *StoreUpdater) UpdateAllStores(payload stock.StockUpdatePayload) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(payload.StoreIDs))

	log.Printf("Starting update for stock unit: %s across %d stores", payload.StockUnit, len(payload.StoreIDs))

	for _, storeID := range payload.StoreIDs {
		wg.Add(1)
		go func (id string) {
			defer wg.Done()
			// api latency simulation
			time.Sleep(500*time.Millisecond)
			// actual api call
			// if resp, err := http.Post(url, contentType, body); err != nil {
			// 	errChan <- err
			// }
			log.Printf("-> Updated Store %s (Stock: %d)", id, payload.NewStock)
		}(storeID)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		// returning first error found
		if err != nil {
			return err
		}
	}
	return nil
}