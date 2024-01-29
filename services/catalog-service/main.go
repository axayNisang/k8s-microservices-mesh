package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type CatalogStore struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewCatalogStore() *CatalogStore {
	return &CatalogStore{
		items: map[string]Item{
			"item_k8s_book":   {ID: "item_k8s_book", Name: "Cloud Native Architecture Guide", Price: 49.99, Stock: 150},
			"item_serverless": {ID: "item_serverless", Name: "High Throughput Systems Design", Price: 59.99, Stock: 80},
			"item_devops_pro": {ID: "item_devops_pro", Name: "Enterprise Terraform & Helm", Price: 39.99, Stock: 200},
		},
	}
}

func main() {
	store := NewCatalogStore()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP","service":"catalog-service"}`))
	})

	http.HandleFunc("/catalog/items", func(w http.ResponseWriter, r *http.Request) {
		store.mu.RLock()
		defer store.mu.RUnlock()

		itemList := make([]Item, 0, len(store.items))
		for _, item := range store.items {
			itemList = append(itemList, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(itemList)
	})

	http.HandleFunc("/catalog/items/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/catalog/items/")
		store.mu.RLock()
		item, exists := store.items[id]
		store.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Item not found"}`))
			return
		}
		json.NewEncoder(w).Encode(item)
	})

	log.Println("[CatalogService] Listening on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
