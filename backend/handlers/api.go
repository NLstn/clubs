package handlers

import (
	"net/http"
)

func Handler_v1() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/clubs", handleClubs)
	mux.HandleFunc("/api/v1/clubs/", handleClubs)
	mux.HandleFunc("/api/v1/clubs/{clubid}/members", handleClubMembers)

	return mux
}
