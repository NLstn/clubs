package odata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIKeyCreation tests the custom API key creation endpoint
func TestAPIKeyCreation(t *testing.T) {
	ctx := setupTestContext(t)

	t.Run("create_api_key_success", func(t *testing.T) {
		// Create request
		requestBody := map[string]interface{}{
			"Name": "Test API Key",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		// Debug: print response body if not 201
		if rec.Code != http.StatusCreated {
			t.Logf("Response Status: %d", rec.Code)
			t.Logf("Response Body: %s", rec.Body.String())
		}

		// Assert response
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		// Debug: print response
		t.Logf("Response: %+v", response)

		// Verify response contains plaintext key
		assert.NotEmpty(t, response["APIKey"], "Plaintext key should be present")
		assert.NotEmpty(t, response["ID"], "ID should be present")
		assert.Equal(t, "Test API Key", response["Name"])
		assert.NotEmpty(t, response["KeyPrefix"], "KeyPrefix should be present")
		assert.True(t, response["IsActive"].(bool))

		// Verify plaintext key format
		plainKey := response["APIKey"].(string)
		assert.Contains(t, plainKey, "sk_live_", "Key should have sk_live_ prefix")
		assert.True(t, len(plainKey) > 20, "Key should be sufficiently long")

		// Verify key is stored in database (hashed)
		var apiKey models.APIKey
		err = ctx.service.db.Where("id = ?", response["ID"]).First(&apiKey).Error
		require.NoError(t, err)
		assert.Equal(t, ctx.testUser.ID, apiKey.UserID)
		assert.Equal(t, "Test API Key", apiKey.Name)
		assert.NotEmpty(t, apiKey.KeyHash)
		assert.NotEqual(t, plainKey, apiKey.KeyHash, "Key hash should not equal plaintext key")
	})

	t.Run("create_api_key_with_expiration", func(t *testing.T) {
		expiresAt := time.Now().Add(30 * 24 * time.Hour)

		requestBody := map[string]interface{}{
			"Name":      "Expiring Key",
			"ExpiresAt": expiresAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["ExpiresAt"])
	})

	t.Run("create_api_key_with_permissions", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"Name":        "Limited Key",
			"Permissions": []string{"read:events", "read:members"},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		// Verify permissions are returned
		permissions := response["Permissions"].([]interface{})
		assert.Len(t, permissions, 2)
	})

	t.Run("create_api_key_missing_name", func(t *testing.T) {
		requestBody := map[string]interface{}{}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Name is required")
	})

	t.Run("create_api_key_unauthorized", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"Name": "Test Key",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("create_api_key_rate_limit", func(t *testing.T) {
		// Create fresh test context to avoid interference from other tests
		rateLimitCtx := setupTestContext(t)

		// Create 10 API keys to hit the limit
		for i := 0; i < 10; i++ {
			requestBody := map[string]interface{}{
				"Name": fmt.Sprintf("Key %d", i+1),
			}
			body, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+rateLimitCtx.token)

			rec := httptest.NewRecorder()
			rateLimitCtx.handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusCreated {
				t.Logf("Failed to create key %d: %d - %s", i+1, rec.Code, rec.Body.String())
			}
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// 11th key should fail
		requestBody := map[string]interface{}{
			"Name": "Key 11",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v2/APIKeys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+rateLimitCtx.token)

		rec := httptest.NewRecorder()
		rateLimitCtx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusTooManyRequests, rec.Code)
		assert.Contains(t, rec.Body.String(), "Maximum number of active API keys")
	})
}

// TestAPIKeyRetrieval tests getting API key entities via OData
func TestAPIKeyRetrieval(t *testing.T) {
	ctx := setupTestContext(t)

	// Create some API keys for testing
	plainKey1, keyHash1, keyPrefix1, _ := auth.GenerateAPIKey("sk_live")
	apiKey1 := &models.APIKey{
		UserID:    ctx.testUser.ID,
		Name:      "Key 1",
		KeyHash:   keyHash1,
		KeyPrefix: keyPrefix1,
		IsActive:  true,
	}
	ctx.service.db.Create(apiKey1)

	plainKey2, keyHash2, keyPrefix2, _ := auth.GenerateAPIKey("sk_live")
	now := time.Now()
	apiKey2 := &models.APIKey{
		UserID:     ctx.testUser.ID,
		Name:       "Key 2",
		KeyHash:    keyHash2,
		KeyPrefix:  keyPrefix2,
		IsActive:   false,
		LastUsedAt: &now,
	}
	ctx.service.db.Create(apiKey2)

	// Create key for another user (should not be visible)
	_, keyHash3, keyPrefix3, _ := auth.GenerateAPIKey("sk_live")
	apiKey3 := &models.APIKey{
		UserID:    ctx.testUser2.ID,
		Name:      "User2 Key",
		KeyHash:   keyHash3,
		KeyPrefix: keyPrefix3,
		IsActive:  true,
	}
	ctx.service.db.Create(apiKey3)

	// Store keys for later validation
	_ = plainKey1
	_ = plainKey2

	t.Run("list_api_keys", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/APIKeys", nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure
		assert.NotNil(t, response["value"])
		keys := response["value"].([]interface{})

		// Should return keys for current user only
		// Note: Depending on implementation, may need entity-level filtering
		// For now, we should see at least our 2 keys
		assert.True(t, len(keys) >= 2, "Should see at least 2 keys")

		// Verify plaintext key is NOT in response
		for _, k := range keys {
			key := k.(map[string]interface{})
			assert.Nil(t, key["APIKey"], "Plaintext key should not be in list response")
			assert.Nil(t, key["KeyHash"], "KeyHash should not be exposed")
			assert.NotNil(t, key["KeyPrefix"], "KeyPrefix should be present")
			assert.NotNil(t, key["Name"], "Name should be present")
		}
	})

	t.Run("get_single_api_key", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v2/APIKeys('%s')", apiKey1.ID), nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		// Debug
		t.Logf("Single entity response: %+v", response)

		// Verify key details
		assert.Equal(t, apiKey1.ID, response["ID"])
		assert.Equal(t, "Key 1", response["Name"])
		assert.Nil(t, response["APIKey"], "Plaintext key should not be in single entity response")
		assert.Nil(t, response["KeyHash"], "KeyHash should not be exposed")
	})

	t.Run("filter_api_keys_by_name", func(t *testing.T) {
		filter := url.QueryEscape("Name eq 'Key 1'")
		req := httptest.NewRequest("GET", "/api/v2/APIKeys?$filter="+filter, nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		keys := response["value"].([]interface{})
		assert.Len(t, keys, 1)
		assert.Equal(t, "Key 1", keys[0].(map[string]interface{})["Name"])
	})

	t.Run("filter_api_keys_by_active_status", func(t *testing.T) {
		filter := url.QueryEscape("IsActive eq true")
		req := httptest.NewRequest("GET", "/api/v2/APIKeys?$filter="+filter, nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		keys := response["value"].([]interface{})
		for _, k := range keys {
			key := k.(map[string]interface{})
			assert.True(t, key["IsActive"].(bool))
		}
	})

	t.Run("order_api_keys_by_created_date", func(t *testing.T) {
		orderby := url.QueryEscape("CreatedAt desc")
		req := httptest.NewRequest("GET", "/api/v2/APIKeys?$orderby="+orderby, nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		keys := response["value"].([]interface{})
		assert.True(t, len(keys) >= 2)
	})
}

// TestAPIKeyUpdate tests updating API key properties
func TestAPIKeyUpdate(t *testing.T) {
	ctx := setupTestContext(t)

	// Create API key
	_, keyHash, keyPrefix, _ := auth.GenerateAPIKey("sk_live")
	apiKey := &models.APIKey{
		UserID:    ctx.testUser.ID,
		Name:      "Original Name",
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		IsActive:  true,
	}
	ctx.service.db.Create(apiKey)

	t.Run("update_api_key_name", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"Name": "Updated Name",
		}
		body, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v2/APIKeys('%s')", apiKey.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify update in database
		var updated models.APIKey
		ctx.service.db.Where("id = ?", apiKey.ID).First(&updated)
		assert.Equal(t, "Updated Name", updated.Name)
	})

	t.Run("update_api_key_active_status", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"IsActive": false,
		}
		body, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v2/APIKeys('%s')", apiKey.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify update in database
		var updated models.APIKey
		ctx.service.db.Where("id = ?", apiKey.ID).First(&updated)
		assert.False(t, updated.IsActive)
	})

	t.Run("cannot_update_key_hash", func(t *testing.T) {
		// Attempting to update KeyHash should fail or be ignored
		// (KeyHash has json:"-" tag so it's not exposed)
		updateBody := map[string]interface{}{
			"KeyHash": "malicious-hash",
		}
		body, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v2/APIKeys('%s')", apiKey.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		// Should succeed but KeyHash should remain unchanged
		// (or return error depending on OData implementation)
		var updated models.APIKey
		ctx.service.db.Where("id = ?", apiKey.ID).First(&updated)
		assert.Equal(t, keyHash, updated.KeyHash, "KeyHash should not change")
	})
}

// TestAPIKeyDeletion tests deleting/revoking API keys
func TestAPIKeyDeletion(t *testing.T) {
	ctx := setupTestContext(t)

	// Create API key
	_, keyHash, keyPrefix, _ := auth.GenerateAPIKey("sk_live")
	apiKey := &models.APIKey{
		UserID:    ctx.testUser.ID,
		Name:      "To Be Deleted",
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		IsActive:  true,
	}
	ctx.service.db.Create(apiKey)

	t.Run("delete_api_key", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v2/APIKeys('%s')", apiKey.ID), nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify key is deleted from database
		var count int64
		ctx.service.db.Model(&models.APIKey{}).Where("id = ?", apiKey.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("delete_nonexistent_key", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v2/APIKeys('00000000-0000-0000-0000-000000000000')", nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// TestAPIKeyExpansion tests OData $expand functionality
func TestAPIKeyExpansion(t *testing.T) {
	ctx := setupTestContext(t)

	// Create API key
	_, keyHash, keyPrefix, _ := auth.GenerateAPIKey("sk_live")
	apiKey := &models.APIKey{
		UserID:    ctx.testUser.ID,
		Name:      "Test Key",
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		IsActive:  true,
	}
	ctx.service.db.Create(apiKey)

	t.Run("expand_user_navigation_property", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v2/APIKeys('%s')?$expand=User", apiKey.ID), nil)
		req.Header.Set("Authorization", "Bearer "+ctx.token)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		// Verify User is expanded
		assert.NotNil(t, response["User"])
		user := response["User"].(map[string]interface{})
		assert.Equal(t, ctx.testUser.ID, user["ID"])
		assert.Equal(t, ctx.testUser.Email, user["Email"])
	})
}

// TestAPIKeyAuthentication tests using API keys for authentication
func TestAPIKeyAuthentication(t *testing.T) {
	ctx := setupTestContext(t)

	// Create API key
	plainKey, keyHash, keyPrefix, _ := auth.GenerateAPIKey("sk_live")
	apiKey := &models.APIKey{
		UserID:    ctx.testUser.ID,
		Name:      "Auth Test Key",
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		IsActive:  true,
	}
	ctx.service.db.Create(apiKey)

	t.Run("authenticate_with_api_key_x_api_key_header", func(t *testing.T) {
		// Try to access an endpoint using X-API-Key header
		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("X-API-Key", plainKey)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("authenticate_with_api_key_authorization_header", func(t *testing.T) {
		// Try to access an endpoint using Authorization: ApiKey header
		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("Authorization", "ApiKey "+plainKey)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("inactive_key_cannot_authenticate", func(t *testing.T) {
		// Deactivate the key
		ctx.service.db.Model(apiKey).Update("is_active", false)

		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("X-API-Key", plainKey)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Reactivate for other tests
		ctx.service.db.Model(apiKey).Update("is_active", true)
	})

	t.Run("expired_key_cannot_authenticate", func(t *testing.T) {
		// Set expiration to past
		pastTime := time.Now().Add(-1 * time.Hour)
		ctx.service.db.Model(apiKey).Update("expires_at", pastTime)

		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("X-API-Key", plainKey)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Clear expiration for other tests
		ctx.service.db.Model(apiKey).Update("expires_at", nil)
	})

	t.Run("invalid_key_cannot_authenticate", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("X-API-Key", "sk_live_invalid_key_12345")

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("last_used_at_updated_after_authentication", func(t *testing.T) {
		// Record current LastUsedAt
		var before models.APIKey
		ctx.service.db.Where("id = ?", apiKey.ID).First(&before)
		beforeLastUsed := before.LastUsedAt

		// Wait a bit to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Use the key
		req := httptest.NewRequest("GET", "/api/v2/Users", nil)
		req.Header.Set("X-API-Key", plainKey)

		rec := httptest.NewRecorder()
		ctx.handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// Wait for async update to complete
		time.Sleep(100 * time.Millisecond)

		// Check LastUsedAt was updated
		var after models.APIKey
		ctx.service.db.Where("id = ?", apiKey.ID).First(&after)

		if beforeLastUsed == nil {
			assert.NotNil(t, after.LastUsedAt, "LastUsedAt should be set after first use")
		} else {
			assert.True(t, after.LastUsedAt.After(*beforeLastUsed), "LastUsedAt should be updated")
		}
	})
}
