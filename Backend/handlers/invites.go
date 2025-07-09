package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func registerInviteRoutes(mux *http.ServeMux) {
	// Admin creates invites
	mux.Handle("/api/v1/clubs/{clubid}/invites", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateInvite(w, r)
		case http.MethodGet:
			handleGetClubInvites(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// User accepts/rejects invites
	mux.Handle("/api/v1/invites", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetUserInvites(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/invites/{inviteid}/accept", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleAcceptInvite(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/invites/{inviteid}/reject", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleRejectInvite(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: POST /api/v1/clubs/{clubid}/invites
func handleCreateInvite(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Email string `json:"email"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)
	if !club.IsOwner(user) && !club.IsAdmin(user) {
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

	err = club.CreateInvite(payload.Email, user.ID)
	if err != nil {
		http.Error(w, "Failed to create invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// endpoint: GET /api/v1/clubs/{clubid}/invites
func handleGetClubInvites(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	invites, err := club.GetInvites()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invites)
}

// endpoint: GET /api/v1/invites
func handleGetUserInvites(w http.ResponseWriter, r *http.Request) {
	type ApiInvite struct {
		ID       string `json:"id"`
		ClubName string `json:"clubName"`
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invites, err := user.GetUserInvites()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var apiInvites []ApiInvite
	for _, invite := range invites {
		club, err := models.GetClubByID(invite.ClubID)
		if err != nil {
			http.Error(w, "Failed to get club information", http.StatusInternalServerError)
			return
		}
		apiInvite := ApiInvite{
			ID:       invite.ID,
			ClubName: club.Name,
		}
		apiInvites = append(apiInvites, apiInvite)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiInvites)
}

// endpoint: POST /api/v1/invites/{inviteid}/accept
func handleAcceptInvite(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	inviteID := extractPathParam(r, "invites")
	if _, err := uuid.Parse(inviteID); err != nil {
		http.Error(w, "Invalid invite ID format", http.StatusBadRequest)
		return
	}

	canEdit, err := user.CanUserEditInvite(inviteID)
	if err != nil || !canEdit {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = models.AcceptInvite(inviteID, user.ID)
	if err != nil {
		http.Error(w, "Failed to accept invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: POST /api/v1/invites/{inviteid}/reject
func handleRejectInvite(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	inviteID := extractPathParam(r, "invites")
	if _, err := uuid.Parse(inviteID); err != nil {
		http.Error(w, "Invalid invite ID format", http.StatusBadRequest)
		return
	}

	canEdit, err := user.CanUserEditInvite(inviteID)
	if err != nil || !canEdit {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = models.RejectInvite(inviteID)
	if err != nil {
		http.Error(w, "Failed to reject invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
