package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupAPIKeyTestDB creates an in-memory SQLite database for testing
func setupAPIKeyTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create tables manually for SQLite compatibility (UUID is stored as TEXT)
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT NOT NULL UNIQUE,
			keycloak_id TEXT UNIQUE,
			birth_date DATE,
			created_at DATETIME,
			updated_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS api_keys (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			key_hash TEXT NOT NULL UNIQUE,
			key_hash_sha256 TEXT UNIQUE,
			key_prefix TEXT NOT NULL,
			permissions TEXT,
			last_used_at DATETIME,
			expires_at DATETIME,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		);
	`).Error
	require.NoError(t, err)

	database.Db = db
}

// TestAPIKeyReadAuthorization tests that users can only read their own API keys
func TestAPIKeyReadAuthorization(t *testing.T) {
	// Setup test database
	setupAPIKeyTestDB(t)

	// Create test users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	err := database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user1ID, "Test", "User1", "user1@test.com").Error
	require.NoError(t, err)

	err = database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user2ID, "Test", "User2", "user2@test.com").Error
	require.NoError(t, err)

	// Create API keys for both users
	apiKey1ID := uuid.New().String()
	apiKey2ID := uuid.New().String()

	err = database.Db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))",
		apiKey1ID, user1ID, "User1 Key", "hash1", "sk_test").Error
	require.NoError(t, err)

	err = database.Db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))",
		apiKey2ID, user2ID, "User2 Key", "hash2", "sk_test").Error
	require.NoError(t, err)

	// Test 1: User1 should only see their own API keys in collection
	t.Run("ReadCollection_UserCanOnlySeeOwnKeys", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/APIKeys", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)
		req = req.WithContext(ctx)

		apiKey := APIKey{} // Empty struct for calling the hook method
		scopes, err := apiKey.ODataBeforeReadCollection(ctx, req, nil)
		require.NoError(t, err)
		require.Len(t, scopes, 1)

		// Apply scope and query
		var keys []APIKey
		query := database.Db.Model(&APIKey{})
		for _, scope := range scopes {
			query = scope(query)
		}
		err = query.Find(&keys).Error
		require.NoError(t, err)

		// User1 should only see their own key
		assert.Len(t, keys, 1)
		assert.Equal(t, user1ID, keys[0].UserID)
		assert.Equal(t, "User1 Key", keys[0].Name)
	})

	// Test 2: User2 should only see their own API keys in collection
	t.Run("ReadCollection_DifferentUserCannotSeeOtherKeys", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/APIKeys", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user2ID)
		req = req.WithContext(ctx)

		apiKey := APIKey{} // Empty struct for calling the hook method
		scopes, err := apiKey.ODataBeforeReadCollection(ctx, req, nil)
		require.NoError(t, err)
		require.Len(t, scopes, 1)

		// Apply scope and query
		var keys []APIKey
		query := database.Db.Model(&APIKey{})
		for _, scope := range scopes {
			query = scope(query)
		}
		err = query.Find(&keys).Error
		require.NoError(t, err)

		// User2 should only see their own key
		assert.Len(t, keys, 1)
		assert.Equal(t, user2ID, keys[0].UserID)
		assert.Equal(t, "User2 Key", keys[0].Name)
	})

	// Test 3: User1 should be able to read their own API key entity
	t.Run("ReadEntity_UserCanReadOwnKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/APIKeys('"+apiKey1ID+"')", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)
		req = req.WithContext(ctx)

		apiKey := APIKey{} // Empty struct for calling the hook method
		scopes, err := apiKey.ODataBeforeReadEntity(ctx, req, nil)
		require.NoError(t, err)
		require.Len(t, scopes, 1)

		// Apply scope and query
		var key APIKey
		query := database.Db.Model(&APIKey{}).Where("id = ?", apiKey1ID)
		for _, scope := range scopes {
			query = scope(query)
		}
		err = query.First(&key).Error
		require.NoError(t, err)

		// Verify it's user1's key
		assert.Equal(t, user1ID, key.UserID)
		assert.Equal(t, "User1 Key", key.Name)
	})

	// Test 4: User1 should NOT be able to read User2's API key entity
	t.Run("ReadEntity_UserCannotReadOtherUsersKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/APIKeys('"+apiKey2ID+"')", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)
		req = req.WithContext(ctx)

		apiKey := APIKey{} // Empty struct for calling the hook method
		scopes, err := apiKey.ODataBeforeReadEntity(ctx, req, nil)
		require.NoError(t, err)
		require.Len(t, scopes, 1)

		// Apply scope and query - should return no results
		var key APIKey
		query := database.Db.Model(&APIKey{}).Where("id = ?", apiKey2ID)
		for _, scope := range scopes {
			query = scope(query)
		}
		err = query.First(&key).Error

		// Should fail with record not found
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Test 5: Unauthenticated request should be rejected
	t.Run("ReadCollection_UnauthenticatedRequestRejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/APIKeys", nil)
		ctx := req.Context() // No user ID in context

		apiKey := APIKey{} // Empty struct for calling the hook method
		_, err := apiKey.ODataBeforeReadCollection(ctx, req, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	// Test 6: Unauthenticated entity read should be rejected
	t.Run("ReadEntity_UnauthenticatedRequestRejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/APIKeys('"+apiKey1ID+"')", nil)
		ctx := req.Context() // No user ID in context

		apiKey := APIKey{} // Empty struct for calling the hook method
		_, err := apiKey.ODataBeforeReadEntity(ctx, req, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

// TestAPIKeyCreateAuthorization tests that users can only create API keys for themselves
func TestAPIKeyCreateAuthorization(t *testing.T) {
	// Setup test database
	setupAPIKeyTestDB(t)

	// Create test users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	err := database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user1ID, "Test", "User1", "user1@test.com").Error
	require.NoError(t, err)

	err = database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user2ID, "Test", "User2", "user2@test.com").Error
	require.NoError(t, err)

	// Test 1: User can create API key for themselves (no UserID set)
	t.Run("Create_UserCanCreateOwnKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v2/APIKeys", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)

		apiKey := &APIKey{
			Name:      "Test Key",
			KeyHash:   "testhash",
			KeyPrefix: "sk_test",
		}

		err := apiKey.ODataBeforeCreate(ctx, req)
		require.NoError(t, err)

		// Verify UserID was set automatically
		assert.Equal(t, user1ID, apiKey.UserID)
		assert.False(t, apiKey.CreatedAt.IsZero())
		assert.False(t, apiKey.UpdatedAt.IsZero())
	})

	// Test 2: User can create API key with their own UserID explicitly set
	t.Run("Create_UserCanCreateWithOwnUserID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v2/APIKeys", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)

		apiKey := &APIKey{
			UserID:    user1ID,
			Name:      "Test Key 2",
			KeyHash:   "testhash2",
			KeyPrefix: "sk_test",
		}

		err := apiKey.ODataBeforeCreate(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, user1ID, apiKey.UserID)
	})

	// Test 3: User CANNOT create API key for another user
	t.Run("Create_UserCannotCreateForOtherUser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v2/APIKeys", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)

		apiKey := &APIKey{
			UserID:    user2ID, // Trying to create for user2
			Name:      "Malicious Key",
			KeyHash:   "testhash3",
			KeyPrefix: "sk_test",
		}

		err := apiKey.ODataBeforeCreate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot create API keys for other users")
	})

	// Test 4: Unauthenticated request should be rejected
	t.Run("Create_UnauthenticatedRequestRejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v2/APIKeys", nil)
		ctx := req.Context() // No user ID in context

		apiKey := &APIKey{
			Name:      "Test Key",
			KeyHash:   "testhash",
			KeyPrefix: "sk_test",
		}

		err := apiKey.ODataBeforeCreate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

// TestAPIKeyUpdateAuthorization tests that users can only update their own API keys
func TestAPIKeyUpdateAuthorization(t *testing.T) {
	// Setup test database
	setupAPIKeyTestDB(t)

	// Create test users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	err := database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user1ID, "Test", "User1", "user1@test.com").Error
	require.NoError(t, err)

	err = database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user2ID, "Test", "User2", "user2@test.com").Error
	require.NoError(t, err)

	// Create API key for user2
	apiKeyID := uuid.New().String()
	err = database.Db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, created_at, updated_at) VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
		apiKeyID, user2ID, "User2 Key", "hash2", "sk_test").Error
	require.NoError(t, err)

	apiKey := &APIKey{
		ID:        apiKeyID,
		UserID:    user2ID,
		Name:      "User2 Key",
		KeyHash:   "hash2",
		KeyPrefix: "sk_test",
	}

	// Test 1: User1 cannot update User2's API key
	t.Run("Update_UserCannotUpdateOtherUsersKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v2/APIKeys('"+apiKeyID+"')", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)

		err := apiKey.ODataBeforeUpdate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot update API keys of other users")
	})

	// Test 2: User2 can update their own API key
	t.Run("Update_UserCanUpdateOwnKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v2/APIKeys('"+apiKeyID+"')", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user2ID)

		err := apiKey.ODataBeforeUpdate(ctx, req)
		assert.NoError(t, err)
	})

	// Test 3: Unauthenticated request should be rejected
	t.Run("Update_UnauthenticatedRequestRejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v2/APIKeys('"+apiKeyID+"')", nil)
		ctx := req.Context() // No user ID in context

		err := apiKey.ODataBeforeUpdate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

// TestAPIKeyDeleteAuthorization tests that users can only delete their own API keys
func TestAPIKeyDeleteAuthorization(t *testing.T) {
	// Setup test database
	setupAPIKeyTestDB(t)

	// Create test users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	err := database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user1ID, "Test", "User1", "user1@test.com").Error
	require.NoError(t, err)

	err = database.Db.Exec("INSERT INTO users (id, first_name, last_name, email, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		user2ID, "Test", "User2", "user2@test.com").Error
	require.NoError(t, err)

	// Create API key for user2
	apiKeyID := uuid.New().String()
	err = database.Db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, created_at, updated_at) VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
		apiKeyID, user2ID, "User2 Key", "hash2", "sk_test").Error
	require.NoError(t, err)

	apiKey := &APIKey{
		ID:        apiKeyID,
		UserID:    user2ID,
		Name:      "User2 Key",
		KeyHash:   "hash2",
		KeyPrefix: "sk_test",
	}

	// Test 1: User1 cannot delete User2's API key
	t.Run("Delete_UserCannotDeleteOtherUsersKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v2/APIKeys('"+apiKeyID+"')", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user1ID)

		err := apiKey.ODataBeforeDelete(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete API keys of other users")
	})

	// Test 2: User2 can delete their own API key
	t.Run("Delete_UserCanDeleteOwnKey", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v2/APIKeys('"+apiKeyID+"')", nil)
		ctx := context.WithValue(req.Context(), auth.UserIDKey, user2ID)

		err := apiKey.ODataBeforeDelete(ctx, req)
		assert.NoError(t, err)
	})

	// Test 3: Unauthenticated request should be rejected
	t.Run("Delete_UnauthenticatedRequestRejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v2/APIKeys('"+apiKeyID+"')", nil)
		ctx := req.Context() // No user ID in context

		err := apiKey.ODataBeforeDelete(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}
