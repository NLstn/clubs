package handlers

import (
	"net/http"
)

func Handler_v1() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/clubs", handleClubs)
	mux.HandleFunc("/api/v1/clubs/", handleClubs)
	mux.HandleFunc("/api/v1/clubs/{clubid}/members", handleClubMembers)
	mux.HandleFunc("/api/v1/clubs/{clubid}/events", handleClubEvents)
	mux.HandleFunc("/api/v1/clubs/{clubid}/events/", handleClubEvents)

	mux.HandleFunc("/api/v1/auth/requestMagicLink", requestMagicLink)
	mux.HandleFunc("/api/v1/auth/verifyMagicLink", verifyMagicLink)

	return mux
}
