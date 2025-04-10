package handlers

import (
	"net/http"

	"github.com/NLstn/clubs/auth"
)

func withAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler := http.HandlerFunc(h)
		auth.AuthMiddleware(handler).ServeHTTP(w, r)
	}
}

func Handler_v1() http.Handler {
	mux := http.NewServeMux()

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

	mux.HandleFunc("/api/v1/clubs/{clubid}/members", withAuth(handleClubMembers))

	mux.HandleFunc("/api/v1/clubs/{clubid}/members/{memberid}", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handleClubMemberDelete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

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
