package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
)

func Handler_v1() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/auth/requestMagicLink", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleRequestMagicLink(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/verifyMagicLink", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			verifyMagicLink(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/refreshToken", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleRefreshToken(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/logout", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleLogout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Authenticated endpoints with general API rate limiting
	mux.Handle("/api/v1/clubs", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAllClubs(w, r)
		case http.MethodPost:
			handleCreateClub(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetClubByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/members", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubMembers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/members/{memberid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handleClubMemberDelete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/events", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubEvents(w, r)
		case http.MethodPost:
			handleCreateClubEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/events/{eventid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handleDeleteClubEvent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/me", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetMe(w, r)
		case http.MethodPut:
			handleUpdateMe(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

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

	return LoggingMiddleware(CorsMiddleware(mux))
}

func withAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler := http.HandlerFunc(h)
		auth.AuthMiddleware(handler).ServeHTTP(w, r)
	}
}

func extractPathParam(r *http.Request, param string) string {
	parts := strings.Split(r.URL.Path, "/")
	for i, part := range parts {
		if part == param && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func extractUser(r *http.Request) models.User {
	userID := r.Context().Value(auth.UserIDKey).(string)
	if userID == "" {
		return models.User{}
	}
	user, err := models.GetUserByID(userID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		return models.User{}
	}
	return user
}
