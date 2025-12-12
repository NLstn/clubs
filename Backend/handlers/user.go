package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/models"
)

func registerUserRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/me/sessions", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetMySessions(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/me/sessions/", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleDeleteMySession(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/me/sessions
func handleGetMySessions(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessions, err := user.GetActiveSessions()
	if err != nil {
		http.Error(w, "Failed to get sessions", http.StatusInternalServerError)
		return
	}

	// Get current refresh token from X-Refresh-Token header to identify current session
	currentRefreshToken := r.Header.Get("X-Refresh-Token")
	var hashedCurrentRefreshToken string
	if currentRefreshToken != "" {
		hashedCurrentRefreshToken = models.HashToken(currentRefreshToken)
	}

	// Transform sessions for response
	type SessionResponse struct {
		ID        string    `json:"id"`
		UserAgent string    `json:"userAgent"`
		IPAddress string    `json:"ipAddress"`
		CreatedAt time.Time `json:"createdAt"`
		IsCurrent bool      `json:"isCurrent"`
	}

	var result []SessionResponse
	for _, session := range sessions {
		result = append(result, SessionResponse{
			ID:        session.ID,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			CreatedAt: session.CreatedAt,
			IsCurrent: hashedCurrentRefreshToken != "" && session.Token == hashedCurrentRefreshToken,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// endpoint: DELETE /api/v1/me/sessions/{sessionId}
func handleDeleteMySession(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract session ID from URL path
	path := r.URL.Path
	sessionID := ""
	if len(path) > len("/api/v1/me/sessions/") {
		sessionID = path[len("/api/v1/me/sessions/"):]
	}

	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	if err := user.DeleteSession(sessionID); err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


