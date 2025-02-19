package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db *gorm.DB

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
		path := strings.Trim(r.URL.Path, "/")
		segments := strings.Split(path, "/")

		if len(segments) == 4 && segments[3] != "" {
			id := segments[3]
			var club models.Club
			result := db.First(&club, "id = ?", id)

			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Club not found", http.StatusNotFound)
				return
			}
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(club)
			return
		}

		var clubs []models.Club
		if result := db.Find(&clubs); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clubs)

	case http.MethodPost:
		var club models.Club
		if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		club.ID = uuid.New().String()

		if result := db.Create(&club); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(club)
	}
}

func handleClubMembers(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")

	if len(segments) != 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	clubID := segments[3]

	var club models.Club
	if result := db.First(&club, "id = ?", clubID); result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var members []models.Member
		if result := db.Where("club_id = ?", clubID).Find(&members); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)

	case http.MethodPost:
		var member models.Member
		if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if member.Email == "" || member.Name == "" {
			http.Error(w, "Email and name are required", http.StatusBadRequest)
			return
		}

		member.ID = uuid.New().String()
		member.ClubID = clubID

		if result := db.Create(&member); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(member)

	case http.MethodDelete:
		memberID := r.URL.Query().Get("id")
		if memberID == "" {
			http.Error(w, "Member ID parameter is required", http.StatusBadRequest)
			return
		}

		result := db.Where("id = ? AND club_id = ?", memberID, clubID).Delete(&models.Member{})
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	dbPort, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		log.Fatal("DATABASE_PORT environment variable is required")
	}

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal("DATABASE_PORT must be an integer")
	}

	dbUser := os.Getenv("DATABASE_USER")
	if dbUser == "" {
		log.Fatal("DATABASE_USER environment variable is required")
	}

	dbUserPassword := os.Getenv("DATABASE_USER_PASSWORD")
	if dbUserPassword == "" {
		log.Fatal("DATABASE_USER_PASSWORD environment variable is required")
	}

	config := &database.Config{
		Host:     dbUrl,
		Port:     dbPortInt,
		User:     dbUser,
		Password: dbUserPassword,
		DBName:   "clubs",
	}

	db, err = database.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

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
