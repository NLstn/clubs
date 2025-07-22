package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get Me - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/me", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Me - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "getme@example.com")

		req := MakeRequest(t, "GET", "/api/v1/me", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var userResponse map[string]interface{}
		ParseJSONResponse(t, rr, &userResponse)
		assert.Equal(t, user.Email, userResponse["Email"])
		assert.Equal(t, user.FirstName, userResponse["FirstName"])
		assert.Equal(t, user.LastName, userResponse["LastName"])
		assert.Equal(t, user.ID, userResponse["ID"])
	})

	t.Run("Update Me - Unauthorized", func(t *testing.T) {
		updateData := map[string]string{
			"firstName": "Updated",
			"lastName":  "Name",
		}

		req := MakeRequest(t, "PUT", "/api/v1/me", updateData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Update Me - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "updateme@example.com")
		updateData := map[string]string{
			"firstName": "Updated",
			"lastName":  "TestUser",
		}

		req := MakeRequest(t, "PUT", "/api/v1/me", updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify the user was updated by fetching again
		req2 := MakeRequest(t, "GET", "/api/v1/me", nil, token)
		rr2 := ExecuteRequest(t, handler, req2)
		CheckResponseCode(t, http.StatusOK, rr2.Code)

		var userResponse map[string]interface{}
		ParseJSONResponse(t, rr2, &userResponse)
		assert.Equal(t, "Updated", userResponse["FirstName"])
		assert.Equal(t, "TestUser", userResponse["LastName"])
		assert.Equal(t, user.Email, userResponse["Email"])
	})

	t.Run("Update Me - Missing FirstName", func(t *testing.T) {
		_, token := CreateTestUser(t, "updateme2@example.com")
		updateData := map[string]string{
			"lastName": "User",
		}

		req := MakeRequest(t, "PUT", "/api/v1/me", updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "First name and last name are required")
	})

	t.Run("Update Me - Missing LastName", func(t *testing.T) {
		_, token := CreateTestUser(t, "updateme3@example.com")
		updateData := map[string]string{
			"firstName": "Test",
		}

		req := MakeRequest(t, "PUT", "/api/v1/me", updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "First name and last name are required")
	})

	t.Run("Update Me - Empty FirstName", func(t *testing.T) {
		_, token := CreateTestUser(t, "updateme4@example.com")
		updateData := map[string]string{
			"firstName": "",
			"lastName":  "User",
		}

		req := MakeRequest(t, "PUT", "/api/v1/me", updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "First name and last name are required")
	})

	t.Run("Update Me - Empty LastName", func(t *testing.T) {
		_, token := CreateTestUser(t, "updateme5@example.com")
		updateData := map[string]string{
			"firstName": "Test",
			"lastName":  "",
		}

		req := MakeRequest(t, "PUT", "/api/v1/me", updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "First name and last name are required")
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		_, token := CreateTestUser(t, "methodtest@example.com")

		// Test unsupported methods
		methods := []string{"POST", "DELETE", "PATCH"}

		for _, method := range methods {
			req := MakeRequest(t, method, "/api/v1/me", nil, token)
			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
		}
	})
}
