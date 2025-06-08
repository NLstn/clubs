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

	// Register health check routes (no auth required)
	registerHealthRoutes(mux)

	registerAuthRoutes(mux)
	registerClubRoutes(mux)
	registerMemberRoutes(mux)
	registerShiftRoutes(mux)
	registerEventRoutes(mux)
	registerNewsRoutes(mux)
	registerJoinRequestRoutes(mux)
	registerFineRoutes(mux)
	registerFineTemplateRoutes(mux)

	registerUserRoutes(mux)

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
