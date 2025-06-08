package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/clubs/database"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Set up test environment
	MockEnvironmentVariables(t)

	// Store original database state and reset it for this test
	originalDb := database.Db
	database.Db = nil
	defer func() {
		database.Db = originalDb
	}()

	// Create a request to the health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the health check handler
	HealthCheck(rr, req)

	// In test environment without real database, health check should return 503
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code, "Health check should return 503 when database is unavailable")

	// Check the content type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Content-Type should be application/json")

	// Parse the response body
	var response HealthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Check the response structure
	assert.Equal(t, "unhealthy", response.Status, "Status should be unhealthy when database is unavailable")
	assert.Contains(t, response.Services, "api", "Response should contain api service status")
	assert.Equal(t, "healthy", response.Services["api"], "API service should be healthy")
	assert.Contains(t, response.Services, "database", "Response should contain database service status")
	assert.Equal(t, "unavailable", response.Services["database"], "Database service should be unavailable in test environment")
}

func TestHealthCheckResponseStructure(t *testing.T) {
	// Set up test environment
	MockEnvironmentVariables(t)

	// Store original database state and reset it for this test
	originalDb := database.Db
	database.Db = nil
	defer func() {
		database.Db = originalDb
	}()

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	HealthCheck(rr, req)

	var response HealthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response has required fields
	assert.NotEmpty(t, response.Status, "Status field should not be empty")
	assert.NotNil(t, response.Services, "Services field should not be nil")
	assert.IsType(t, map[string]string{}, response.Services, "Services should be a map of strings")
}