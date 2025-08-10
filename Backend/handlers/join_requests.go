package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func registerJoinRequestRoutes(mux *http.ServeMux) {
	// Admin views join requests
	mux.Handle("/api/v1/clubs/{clubid}/joinRequests", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetJoinRequests(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Admin gets join request count
	mux.Handle("/api/v1/clubs/{clubid}/joinRequests/count", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetJoinRequestCount(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Invite link generation
	mux.Handle("/api/v1/clubs/{clubid}/inviteLink", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetInviteLink(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// User joins via link
	mux.Handle("/api/v1/clubs/{clubid}/join", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleJoinClubViaLink(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Club info for invitation
	mux.Handle("/api/v1/clubs/{clubid}/info", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetClubInfo(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Admin accepts/rejects join requests
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

// endpoint: GET /api/v1/clubs/{clubid}/joinRequests
func handleGetJoinRequests(w http.ResponseWriter, r *http.Request) {
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

	requests, err := club.GetJoinRequests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// endpoint: GET /api/v1/clubs/{clubid}/joinRequests/count
func handleGetJoinRequestCount(w http.ResponseWriter, r *http.Request) {
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

	count, err := club.GetJoinRequestCount()
	if err != nil {
		http.Error(w, "Failed to get join request count", http.StatusInternalServerError)
		return
	}

	response := map[string]int64{"count": count}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// endpoint: POST /api/v1/joinRequests/{requestid}/accept
func handleAcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	requestID := extractPathParam(r, "joinRequests")
	if _, err := uuid.Parse(requestID); err != nil {
		http.Error(w, "Invalid request ID format", http.StatusBadRequest)
		return
	}

	err := models.AcceptJoinRequest(requestID, user.ID)
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

	err := models.RejectJoinRequest(requestID, user.ID)
	if err != nil {
		http.Error(w, "Failed to reject join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/clubs/{clubid}/inviteLink
func handleGetInviteLink(w http.ResponseWriter, r *http.Request) {
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

	// For now, we'll use a simple format: club ID as the invitation parameter
	// In production, you might want to add a secure token or expiration
	inviteLink := map[string]string{
		"inviteLink": "/join/" + clubID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inviteLink)
}

// endpoint: POST /api/v1/clubs/{clubid}/join
func handleJoinClubViaLink(w http.ResponseWriter, r *http.Request) {
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
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is already a member
	if club.IsMember(user) {
		http.Error(w, "User is already a member of this club", http.StatusConflict)
		return
	}

	// Check if user already has a pending join request
	hasPendingRequest, err := club.HasPendingJoinRequest(user.ID)
	if err != nil {
		http.Error(w, "Failed to check join request status", http.StatusInternalServerError)
		return
	}
	if hasPendingRequest {
		http.Error(w, "You already have a pending join request for this club", http.StatusConflict)
		return
	}

	// Check if user already has a pending invite
	hasPendingInvite, err := club.HasPendingInvite(user.Email)
	if err != nil {
		http.Error(w, "Failed to check invite status", http.StatusInternalServerError)
		return
	}
	if hasPendingInvite {
		http.Error(w, "You already have a pending invitation for this club. Please check your profile invitations page", http.StatusConflict)
		return
	}

	// Create a join request
	err = club.CreateJoinRequest(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to create join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// endpoint: GET /api/v1/clubs/{clubid}/info
func handleGetClubInfo(w http.ResponseWriter, r *http.Request) {
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
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is already a member
	isMember := club.IsMember(user)

	// Check if user has pending join request
	hasPendingRequest, err := club.HasPendingJoinRequest(user.ID)
	if err != nil {
		hasPendingRequest = false // Don't fail, just assume no pending request
	}

	// Check if user has pending invite
	hasPendingInvite, err := club.HasPendingInvite(user.Email)
	if err != nil {
		hasPendingInvite = false // Don't fail, just assume no pending invite
	}

	clubInfo := map[string]interface{}{
		"id":                club.ID,
		"name":              club.Name,
		"description":       club.Description,
		"isMember":          isMember,
		"hasPendingRequest": hasPendingRequest,
		"hasPendingInvite":  hasPendingInvite,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubInfo)
}
