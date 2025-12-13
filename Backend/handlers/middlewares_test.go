package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// dbMutex ensures tests don't interfere with each other's database.Db swapping
var dbMutex sync.Mutex

func setupMiddlewareTestDB(t *testing.T) *gorm.DB {
	// Use temporary file instead of :memory: to avoid connection issues
	tmpfile := t.TempDir() + "/test.db"
	db, err := gorm.Open(sqlite.Open(tmpfile), &gorm.Config{})
	assert.NoError(t, err)

	// Create tables without UUID constraints for SQLite
	err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT NOT NULL,
			keycloak_id TEXT,
			birth_date DATE,
			created_at DATETIME,
			updated_at DATETIME
		)
	`).Error
	assert.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE api_keys (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			key_hash TEXT UNIQUE NOT NULL,
			key_prefix TEXT NOT NULL,
			permissions TEXT,
			last_used_at DATETIME,
			expires_at DATETIME,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`).Error
	assert.NoError(t, err)

	return db
}

func TestAPIKeyAuthMiddleware(t *testing.T) {
	// Lock to ensure tests don't interfere with database.Db swapping
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Setup test database
	db := setupMiddlewareTestDB(t)
	oldDb := database.Db
	database.Db = db
	defer func() {
		// Wait for any async goroutines to complete
		// API key validation updates last_used_at in a goroutine
		time.Sleep(100 * time.Millisecond)
		database.Db = oldDb
	}()

	// Create a test user
	user := models.User{
		ID:    "test-user-middleware",
		Email: "middleware@example.com",
	}
	err := db.Create(&user).Error
	assert.NoError(t, err)

	// Generate a valid API key
	plainKey, keyHash, keyPrefix, err := auth.GenerateAPIKey("sk_live")
	assert.NoError(t, err)

	// Create API key in database using raw SQL
	err = db.Exec(`
		INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?)
	`, "test-key-middleware", user.ID, "Test Middleware Key", keyHash, keyPrefix, time.Now(), time.Now()).Error
	assert.NoError(t, err)

	// Test handler that checks if user ID is in context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey)
		if userID == nil {
			t.Error("Expected user ID in context, got nil")
		}
		assert.Equal(t, user.ID, userID.(string))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	t.Run("Valid API key in X-API-Key header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", plainKey)
		rr := httptest.NewRecorder()

		handler := APIKeyAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authenticated", rr.Body.String())
	})

	t.Run("Valid API key in Authorization: ApiKey header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "ApiKey "+plainKey)
		rr := httptest.NewRecorder()

		handler := APIKeyAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authenticated", rr.Body.String())
	})

	t.Run("Missing API key", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler := APIKeyAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid API key", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", "sk_live_invalid_key")
		rr := httptest.NewRecorder()

		handler := APIKeyAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Inactive API key", func(t *testing.T) {
		// Create an inactive key
		inactiveKey, inactiveHash, inactivePrefix, err := auth.GenerateAPIKey("sk_live")
		assert.NoError(t, err)

		// Insert directly with SQL to ensure boolean is set correctly in SQLite
		err = db.Exec(`
			INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 0, ?, ?)
		`, "test-key-inactive-mw", user.ID, "Inactive Key", inactiveHash, inactivePrefix, time.Now(), time.Now()).Error
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", inactiveKey)
		rr := httptest.NewRecorder()

		handler := APIKeyAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Expired API key", func(t *testing.T) {
		// Create an expired key
		expiredKey, expiredHash, expiredPrefix, err := auth.GenerateAPIKey("sk_live")
		assert.NoError(t, err)

		pastExpiry := time.Now().Add(-24 * time.Hour)
		err = db.Exec(`
			INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, expires_at, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?)
		`, "test-key-expired-mw", user.ID, "Expired Key", expiredHash, expiredPrefix, pastExpiry, time.Now(), time.Now()).Error
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", expiredKey)
		rr := httptest.NewRecorder()

		handler := APIKeyAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestCompositeAuthMiddleware(t *testing.T) {
	// Lock to ensure tests don't interfere with database.Db swapping
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Setup JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-composite")
	err := auth.Init()
	assert.NoError(t, err)

	// Setup test database
	db := setupMiddlewareTestDB(t)
	oldDb := database.Db
	database.Db = db
	defer func() {
		// Wait for any async goroutines to complete
		time.Sleep(10 * time.Millisecond)
		database.Db = oldDb
	}()

	// Create a test user
	user := models.User{
		ID:    "test-user-composite",
		Email: "composite@example.com",
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Generate a valid API key
	plainKey, keyHash, keyPrefix, err := auth.GenerateAPIKey("sk_live")
	assert.NoError(t, err)

	// Create API key in database using raw SQL
	err = db.Exec(`
		INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?)
	`, "test-key-composite", user.ID, "Test Composite Key", keyHash, keyPrefix, time.Now(), time.Now()).Error
	assert.NoError(t, err)

	// Generate a valid JWT token
	jwtToken, err := auth.GenerateAccessToken(user.ID)
	assert.NoError(t, err)

	// Test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey)
		if userID == nil {
			t.Error("Expected user ID in context, got nil")
		}
		assert.Equal(t, user.ID, userID.(string))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	t.Run("Valid JWT Bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authenticated", rr.Body.String())
	})

	t.Run("Valid API key in X-API-Key header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", plainKey)
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authenticated", rr.Body.String())
	})

	t.Run("Valid API key in Authorization: ApiKey header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "ApiKey "+plainKey)
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authenticated", rr.Body.String())
	})

	t.Run("No authentication provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JWT token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid API key", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", "sk_live_invalid_key")
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("JWT token takes precedence when both provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		req.Header.Set("X-API-Key", plainKey)
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authenticated", rr.Body.String())
	})
}

func TestCompositeAuthMiddleware_ContextPropagation(t *testing.T) {
	// Lock to ensure tests don't interfere with database.Db swapping
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Setup
	os.Setenv("JWT_SECRET", "test-secret-context")
	err := auth.Init()
	assert.NoError(t, err)

	db := setupMiddlewareTestDB(t)
	oldDb := database.Db
	database.Db = db
	defer func() {
		// Wait for any async goroutines to complete
		// API key validation updates last_used_at in a goroutine
		time.Sleep(100 * time.Millisecond)
		database.Db = oldDb
	}()

	user := models.User{
		ID:    "test-user-context",
		Email: "context@example.com",
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	plainKey, keyHash, keyPrefix, err := auth.GenerateAPIKey("sk_live")
	assert.NoError(t, err)

	// Create API key in database using raw SQL
	err = db.Exec(`
		INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?)
	`, "test-key-context", user.ID, "Context Test Key", keyHash, keyPrefix, time.Now(), time.Now()).Error
	assert.NoError(t, err)

	t.Run("User ID correctly set in context with API key", func(t *testing.T) {
		var capturedUserID string
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(auth.UserIDKey)
			if userID != nil {
				capturedUserID = userID.(string)
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", plainKey)
		rr := httptest.NewRecorder()

		handler := CompositeAuthMiddleware(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, user.ID, capturedUserID)
	})

	t.Run("Context can be accessed by downstream handlers", func(t *testing.T) {
		middleware1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Add custom value to context
				ctx := context.WithValue(r.Context(), "middleware1", "value1")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}

		finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check both auth and custom context values
			userID := r.Context().Value(auth.UserIDKey)
			customValue := r.Context().Value("middleware1")

			assert.NotNil(t, userID)
			assert.Equal(t, user.ID, userID.(string))
			assert.Equal(t, "value1", customValue.(string))
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", plainKey)
		rr := httptest.NewRecorder()

		// Chain middlewares
		handler := CompositeAuthMiddleware(middleware1(finalHandler))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
