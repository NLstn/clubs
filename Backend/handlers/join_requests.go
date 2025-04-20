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

	events, err := models.GetJoinRequestsForClub(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// endpoint: POST /api/v1/joinRequests/{requestid}/accept
func handleAcceptJoinRequest(w http.ResponseWriter, r *http.Request) {

	userID := extractUserID(r)

	requestID := extractPathParam(r, "joinRequests")
	if _, err := uuid.Parse(requestID); err != nil {
		http.Error(w, "Invalid request ID format", http.StatusBadRequest)
		return
	}

	canEdit, err := models.GetUserCanEditJoinRequest(userID, requestID)
	if err != nil || !canEdit {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = models.AcceptJoinRequest(requestID, userID)
	if err != nil {
		http.Error(w, "Failed to accept join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: POST /api/v1/joinRequests/{requestid}/reject
func handleRejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	userID := extractUserID(r)

	requestID := extractPathParam(r, "joinRequests")
	if _, err := uuid.Parse(requestID); err != nil {
		http.Error(w, "Invalid request ID format", http.StatusBadRequest)
		return
	}

	canEdit, err := models.GetUserCanEditJoinRequest(userID, requestID)
	if err != nil || !canEdit {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = models.RejectJoinRequest(requestID)
	if err != nil {
		http.Error(w, "Failed to reject join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/joinRequests
func handleGetUserJoinRequests(w http.ResponseWriter, r *http.Request) {

	type ApiJoinRequest struct {
		ID       string `json:"id"`
		ClubName string `json:"clubName"`
	}

	userID := extractUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := models.GetUserJoinRequests(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var apiRequests []ApiJoinRequest
	for _, request := range requests {
		club, err := models.GetClubByID(request.ClubID)
		if err != nil {
			http.Error(w, "Failed to get club information", http.StatusInternalServerError)
			return
		}
		apiRequest := ApiJoinRequest{
			ID:       request.ID,
			ClubName: club.Name,
		}
		apiRequests = append(apiRequests, apiRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiRequests)
}
