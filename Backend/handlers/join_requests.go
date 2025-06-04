package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func registerJoinRequestRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/joinRequests", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleJoinRequestCreate(w, r)
		case http.MethodGet:
			handleGetJoinEvents(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/joinRequests", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUserJoinRequests(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/joinRequests/{requestid}/accept", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleAcceptJoinRequest(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/joinRequests/{requestid}/reject", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleRejectJoinRequest(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: POST /api/v1/clubs/{clubid}/joinRequests
func handleJoinRequestCreate(w http.ResponseWriter, r *http.Request) {

	type Payload struct {
		Email string `json:"email"`
	}

	clubID := extractPathParam(r, "clubs")
	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)
	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		http.Error(w, "Missing email", http.StatusBadRequest)
		return
	}

	err = club.CreateJoinRequest(payload.Email, user.ID)
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

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	user := extractUser(r)
	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := club.GetJoinRequests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// endpoint: POST /api/v1/joinRequests/{requestid}/accept
func handleAcceptJoinRequest(w http.ResponseWriter, r *http.Request) {

	user := extractUser(r)

	requestID := extractPathParam(r, "joinRequests")
	if _, err := uuid.Parse(requestID); err != nil {
		http.Error(w, "Invalid request ID format", http.StatusBadRequest)
		return
	}

	canEdit, err := user.GetUserCanEditJoinRequest(requestID)
	if err != nil || !canEdit {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = models.AcceptJoinRequest(requestID, user.ID)
	if err != nil {
		http.Error(w, "Failed to accept join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: POST /api/v1/joinRequests/{requestid}/reject
func handleRejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	requestID := extractPathParam(r, "joinRequests")
	if _, err := uuid.Parse(requestID); err != nil {
		http.Error(w, "Invalid request ID format", http.StatusBadRequest)
		return
	}

	canEdit, err := user.GetUserCanEditJoinRequest(requestID)
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

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := user.GetJoinRequests()
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
