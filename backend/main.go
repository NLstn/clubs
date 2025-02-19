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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []Member `json:"-"`
}

type ClubResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

func handleClubMembers(w http.ResponseWriter, r *http.Request) {
	// Clean the path and split it into segments
	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")

	// Path should be ["api", "v1", "clubs", "{id}", "members"]
	if len(segments) != 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	clubID := segments[3]

	mutex.RLock()
	club, exists := clubs[clubID]
	mutex.RUnlock()

	if !exists {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// List all members
		w.Header().Set("Content-Type", "application/json")
		if club.Members == nil {
			json.NewEncoder(w).Encode([]Member{})
		} else {
			json.NewEncoder(w).Encode(club.Members)
		}

	case http.MethodPost:
		var member Member
		if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if member.Email == "" || member.Name == "" {
			http.Error(w, "Email and name are required", http.StatusBadRequest)
			return
		}

		// Generate new UUID for the member
		member.ID = uuid.New().String()

		mutex.Lock()
		if club.Members == nil {
			club.Members = make([]Member, 0)
		}
		club.Members = append(club.Members, member)
		clubs[clubID] = club
		mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(member)

	case http.MethodDelete:
		memberID := r.URL.Query().Get("id")
		if memberID == "" {
			http.Error(w, "Member ID parameter is required", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		if club.Members != nil {
			// Find and remove the member with the matching ID
			for i, member := range club.Members {
				if member.ID == memberID {
					// Remove the member by slicing
					club.Members = append(club.Members[:i], club.Members[i+1:]...)
					break
				}
			}
			clubs[clubID] = club
		}
		mutex.Unlock()

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	// FIXME: this looks super weird
	mux.HandleFunc("/api/v1/clubs", handleClubs)
	mux.HandleFunc("/api/v1/clubs/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		segments := strings.Split(path, "/")

		// If the path includes "members", use the members handler
		if len(segments) == 5 && segments[4] == "members" {
			handleClubMembers(w, r)
			return
		}

		// Otherwise use the clubs handler
		handleClubs(w, r)
	})

	handler := corsMiddleware(mux)
	handlerWithLogging := loggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
