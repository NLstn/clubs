package handlers

import (
	"net/http"

	"github.com/NLstn/clubs/auth"
)

func withAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create an http.Handler from the handlerFunc
		handler := http.HandlerFunc(h)
		// Apply the middleware and serve the request
		auth.AuthMiddleware(handler).ServeHTTP(w, r)
	}
}

func Handler_v1() http.Handler {
	mux := http.NewServeMux()

	// Route to specific handlers based on method and path
	mux.HandleFunc("/api/v1/clubs", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAllClubs(w, r)
		case http.MethodPost:
			handleCreateClub(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/v1/clubs/", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetClubByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/v1/clubs/{clubid}/members", handleClubMembers)

	// Register event-related endpoints
	mux.HandleFunc("/api/v1/clubs/{clubid}/events", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubEvents(w, r)
		case http.MethodPost:
			handleCreateClubEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/v1/clubs/{clubid}/events/{eventid}", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handleDeleteClubEvent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/v1/auth/requestMagicLink", requestMagicLink)
	mux.HandleFunc("/api/v1/auth/verifyMagicLink", verifyMagicLink)

	return mux
}
