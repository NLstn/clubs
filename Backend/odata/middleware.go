package odata

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware wraps the OData service with JWT authentication
func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for metadata and service document endpoints
			// When wrapped with http.StripPrefix, the path is relative to the mount point
			path := r.URL.Path
			if path == "" || path == "/" || path == "$metadata" || path == "/$metadata" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract and validate JWT token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Expected format: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			// Parse and validate JWT token
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["user_id"] == nil {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			// Add user ID to context for use in read/write hooks
			ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)

			// Phase 5: Parse includeDeleted query parameter
			ctx = ParseIncludeDeletedFromQuery(ctx, r)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Context keys for additional OData features
type contextKey string

const (
	// IncludeDeletedKey enables visibility of soft-deleted items (owners only)
	IncludeDeletedKey contextKey = "includeDeleted"
)

// ParseIncludeDeletedFromQuery checks if the request has ?includeDeleted=true
// and sets it in the context for use in authorization hooks
//
// Phase 5: Complex Scenarios - Soft Delete Visibility
//
// Usage: GET /api/v2/Clubs?includeDeleted=true
//
// Authorization: Only club owners can see their deleted clubs
func ParseIncludeDeletedFromQuery(ctx context.Context, r *http.Request) context.Context {
	// Check query parameter
	if r.URL.Query().Get("includeDeleted") == "true" {
		return context.WithValue(ctx, IncludeDeletedKey, true)
	}

	return ctx
}

// GetIncludeDeletedFromContext retrieves the includeDeleted flag from context
func GetIncludeDeletedFromContext(ctx context.Context) bool {
	includeDeleted, ok := ctx.Value(IncludeDeletedKey).(bool)
	return ok && includeDeleted
}
