package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/stretchr/testify/assert"
)

func TestHandleDevLogin_Success(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Initialize JWT secret for test
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-purposes-only")
	auth.Init()

	// Enable dev auth for this test
	os.Setenv("ENABLE_DEV_AUTH", "true")
	defer os.Unsetenv("ENABLE_DEV_AUTH")

	// Create request body
	requestBody := map[string]string{
		"email": "test@example.com",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/dev-login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handleDevLogin(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify tokens are present
	assert.NotEmpty(t, response["access"], "Access token should be present")
	assert.NotEmpty(t, response["refresh"], "Refresh token should be present")
	assert.Contains(t, response, "profileComplete", "profileComplete should be present")
}

func TestHandleDevLogin_DisabledByDefault(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Make sure ENABLE_DEV_AUTH is not set
	os.Unsetenv("ENABLE_DEV_AUTH")

	// Create request body
	requestBody := map[string]string{
		"email": "test@example.com",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/dev-login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handleDevLogin(w, req)

	// Should return 404 when not enabled
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleDevLogin_MissingEmail(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Initialize JWT secret for test
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-purposes-only")
	auth.Init()

	// Enable dev auth for this test
	os.Setenv("ENABLE_DEV_AUTH", "true")
	defer os.Unsetenv("ENABLE_DEV_AUTH")

	// Create request body without email
	requestBody := map[string]string{}
	bodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/dev-login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handleDevLogin(w, req)

	// Should return 400 for missing email
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleDevLogin_InvalidJSON(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Enable dev auth for this test
	os.Setenv("ENABLE_DEV_AUTH", "true")
	defer os.Unsetenv("ENABLE_DEV_AUTH")

	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/dev-login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handleDevLogin(w, req)

	// Should return 400 for invalid JSON
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleDevLogin_CreatesUserIfNotExists(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Initialize JWT secret for test
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-purposes-only")
	auth.Init()

	// Enable dev auth for this test
	os.Setenv("ENABLE_DEV_AUTH", "true")
	defer os.Unsetenv("ENABLE_DEV_AUTH")

	email := "newuser@example.com"

	// Create request body
	requestBody := map[string]string{
		"email": email,
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/dev-login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handleDevLogin(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify user was created in database
	var count int64
	testDB.Raw("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	assert.Equal(t, int64(1), count, "User should be created in database")
}
