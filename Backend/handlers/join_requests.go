package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

// endpoint: POST /api/v1/clubs/{clubid}/joinRequests
func handleJoinRequestCreate(w http.ResponseWriter, r *http.Request) {

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	userID := extractUserID(r)
	if !auth.IsOwnerOfClub(userID, clubID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var joinRequest models.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	joinRequest.ClubID = clubID

	if joinRequest.ClubID == "" || joinRequest.Email == "" {
		http.Error(w, "Missing club_id or email", http.StatusBadRequest)
		return
	}

	err := models.CreateJoinRequest(joinRequest.ClubID, joinRequest.Email)
	if err != nil {
		http.Error(w, "Failed to create join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// endpoint: GET /api/v1/clubs/{clubid}/joinRequests
func handleGetJoinEvents(w http.ResponseWriter, r *http.Request) {

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	userID := extractUserID(r)
	if !auth.IsOwnerOfClub(userID, clubID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := models.GetJoinRequests(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
