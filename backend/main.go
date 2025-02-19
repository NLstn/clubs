package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

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

		// Create custom response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK, // Default status
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"endpoint=%s method=%s status=%d duration=%v",
			r.URL.Path,
			r.Method,
			rw.status,
			duration,
		)
	})
}

func handleClubs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Clean the path and split it into segments
		path := strings.Trim(r.URL.Path, "/")
		segments := strings.Split(path, "/")

		// Check if we're requesting a specific club
		if len(segments) == 4 && segments[3] != "" { // ["api", "v1", "clubs", "{id}"]
			id := segments[3]
			mutex.RLock()
			club, ok := clubs[id]
			mutex.RUnlock()

			if !ok {
				http.Error(w, "Club not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(club)
			return
		}

		// List all clubs
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

	// FIXME: move clubs handler to its own file and find a smarter way to handle routes
	mux.HandleFunc("/api/v1/clubs", handleClubs)
	mux.HandleFunc("/api/v1/clubs/", handleClubs)

	handler := corsMiddleware(mux)
	handlerWithLogging := loggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
