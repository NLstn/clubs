package odata

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/stretchr/testify/assert"
)

// TestAuthMiddleware tests the OData authentication middleware
func TestAuthMiddleware(t *testing.T) {
	// Set up test JWT secret
	testSecret := "test-secret-key-for-odata-tests"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Initialize auth with test secret
	err := auth.Init()
	assert.NoError(t, err)

	jwtSecret := []byte(testSecret)

	t.Run("valid_bearer_token", func(t *testing.T) {
		// Create a valid token
		token, err := auth.GenerateAccessToken("user-123")
		assert.NoError(t, err)

		middleware := AuthMiddleware(jwtSecret)

		// Create a test handler that checks if user is in context
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(auth.UserIDKey)
			if userID != "user-123" {
				http.Error(w, "User ID not in context", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		})

		wrapped := middleware(testHandler)

		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	})

	t.Run("missing_authorization_header", func(t *testing.T) {
		middleware := AuthMiddleware(jwtSecret)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		wrapped := middleware(testHandler)

		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		// No Authorization header

		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid_bearer_format", func(t *testing.T) {
		middleware := AuthMiddleware(jwtSecret)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		wrapped := middleware(testHandler)

		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("Authorization", "InvalidFormat")

		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid_token", func(t *testing.T) {
		middleware := AuthMiddleware(jwtSecret)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		wrapped := middleware(testHandler)

		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("skip_metadata_endpoint", func(t *testing.T) {
		middleware := AuthMiddleware(jwtSecret)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "metadata")
		})
		wrapped := middleware(testHandler)

		// Test metadata endpoint without token
		req := httptest.NewRequest("GET", "/$metadata", nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "metadata", rec.Body.String())
	})

	t.Run("skip_service_document", func(t *testing.T) {
		middleware := AuthMiddleware(jwtSecret)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "service")
		})
		wrapped := middleware(testHandler)

		// Test service document without token
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "service", rec.Body.String())
	})
}

// TestGetUserIDFromContext tests context user extraction
func TestGetUserIDFromContext(t *testing.T) {
	t.Run("valid_user_id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, "user-123")
		userID, err := getUserIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		ctx := context.Background()
		userID, err := getUserIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "", userID)
	})

	t.Run("invalid_user_id_type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, 123)
		userID, err := getUserIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "", userID)
	})

	t.Run("empty_user_id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, "")
		userID, err := getUserIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "", userID)
	})
}
