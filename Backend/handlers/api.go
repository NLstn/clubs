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

	mux.HandleFunc("/api/v1/clubs", withAuth(handleClubs))
	mux.HandleFunc("/api/v1/clubs/", withAuth(handleClubs))
	mux.HandleFunc("/api/v1/clubs/{clubid}/members", handleClubMembers)
	mux.HandleFunc("/api/v1/clubs/{clubid}/events", handleClubEvents)
	mux.HandleFunc("/api/v1/clubs/{clubid}/events/", handleClubEvents)

	mux.HandleFunc("/api/v1/auth/requestMagicLink", requestMagicLink)
	mux.HandleFunc("/api/v1/auth/verifyMagicLink", verifyMagicLink)

	return mux
}
