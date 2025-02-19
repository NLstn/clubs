package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Club struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage
var (
	clubs = make(map[string]Club)
	mutex = &sync.RWMutex{}
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf(
			"endpoint=%s method=%s duration=%v",
			r.URL.Path,
			r.Method,
			duration,
		)
	})
}

func handleClubs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		clubsList := make([]Club, 0, len(clubs))
		for _, club := range clubs {
			clubsList = append(clubsList, club)
		}
		mutex.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clubsList)

	case http.MethodPost:
		var club Club
		if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Generate new UUID for the club
		club.ID = uuid.New().String()

		mutex.Lock()
		clubs[club.ID] = club
		mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(club)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	// Single API route for all CRUD operations
	mux.HandleFunc("/api/v1/clubs", handleClubs)

	// Wrap with CORS and logging middleware
	handler := corsMiddleware(mux)
	handlerWithLogging := loggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
