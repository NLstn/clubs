package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateFieldFormat(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get My Fines - Should have camelCase date fields", func(t *testing.T) {
		user, token := CreateTestUser(t, "test@example.com")
		club := CreateTestClub(t, user, "Test Club")
		fine := CreateTestFine(t, user, club, "Test fine", 25.0, false)
		_ = fine // Create the fine

		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		// Parse response as raw JSON to check field names
		var responseData []map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &responseData)
		assert.NoError(t, err, "Should be able to parse JSON response")

		// Should have at least one fine
		assert.GreaterOrEqual(t, len(responseData), 1, "Should have at least one fine")

		fine_response := responseData[0]

		// Check that camelCase fields exist and snake_case fields don't
		assert.Contains(t, fine_response, "createdAt", "Should have camelCase createdAt field")
		assert.Contains(t, fine_response, "updatedAt", "Should have camelCase updatedAt field")
		assert.NotContains(t, fine_response, "created_at", "Should not have snake_case created_at field")
		assert.NotContains(t, fine_response, "updated_at", "Should not have snake_case updated_at field")

		// Check that the date values are valid strings
		createdAt, ok := fine_response["createdAt"].(string)
		assert.True(t, ok, "createdAt should be a string")
		assert.NotEmpty(t, createdAt, "createdAt should not be empty")

		updatedAt, ok := fine_response["updatedAt"].(string)
		assert.True(t, ok, "updatedAt should be a string")
		assert.NotEmpty(t, updatedAt, "updatedAt should not be empty")
	})

	t.Run("Get Club Fines - Should have camelCase date fields", func(t *testing.T) {
		adminUser, adminToken := CreateTestUser(t, "admin@example.com")
		club := CreateTestClub(t, adminUser, "Test Club for Admin")
		fine := CreateTestFine(t, adminUser, club, "Admin test fine", 30.0, false)
		_ = fine // Create the fine

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/fines", nil, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		// Parse response as raw JSON to check field names
		var responseData []map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &responseData)
		assert.NoError(t, err, "Should be able to parse JSON response")

		// Should have at least one fine
		assert.GreaterOrEqual(t, len(responseData), 1, "Should have at least one fine")

		fine_response := responseData[0]

		// Check that camelCase fields exist and snake_case fields don't
		assert.Contains(t, fine_response, "createdAt", "Should have camelCase createdAt field")
		assert.Contains(t, fine_response, "updatedAt", "Should have camelCase updatedAt field")
		assert.NotContains(t, fine_response, "created_at", "Should not have snake_case created_at field")
		assert.NotContains(t, fine_response, "updated_at", "Should not have snake_case updated_at field")

		// Check that the date values are valid strings
		createdAt, ok := fine_response["createdAt"].(string)
		assert.True(t, ok, "createdAt should be a string")
		assert.NotEmpty(t, createdAt, "createdAt should not be empty")

		updatedAt, ok := fine_response["updatedAt"].(string)
		assert.True(t, ok, "updatedAt should be a string")
		assert.NotEmpty(t, updatedAt, "updatedAt should not be empty")
	})
}
