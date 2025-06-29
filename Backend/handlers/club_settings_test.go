package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestClubSettingsEndpoints(t *testing.T) {
	// Set up test environment
	MockEnvironmentVariables(t)
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create test user and club
	testUser, token := CreateTestUser(t, "test@example.com")
	club := CreateTestClub(t, testUser, "Test Club")

	handler := GetTestHandler()

	t.Run("Get Club Settings - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/settings", club.ID), nil, "")
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Club Settings - Valid", func(t *testing.T) {
		req := MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/settings", club.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusOK, rr.Code)

		var settings models.ClubSettings
		ParseJSONResponse(t, rr, &settings)

		// Default settings should have both enabled
		assert.True(t, settings.FinesEnabled)
		assert.True(t, settings.ShiftsEnabled)
		assert.Equal(t, club.ID, settings.ClubID)
	})

	t.Run("Update Club Settings - Unauthorized", func(t *testing.T) {
		body := map[string]bool{
			"finesEnabled":  false,
			"shiftsEnabled": true,
		}
		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/settings", club.ID), body, "")
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Update Club Settings - Valid", func(t *testing.T) {
		body := map[string]bool{
			"finesEnabled":  false,
			"shiftsEnabled": true,
		}
		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/settings", club.ID), body, token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify settings were updated
		settings, err := models.GetClubSettings(club.ID)
		assert.NoError(t, err)
		assert.False(t, settings.FinesEnabled)
		assert.True(t, settings.ShiftsEnabled)
	})

	t.Run("Update Club Settings - Invalid JSON", func(t *testing.T) {
		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/settings", club.ID), "invalid json", token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Get Club Settings - Club Not Found", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/clubs/invalid-id/settings", nil, token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Update Club Settings - Club Not Found", func(t *testing.T) {
		body := map[string]bool{
			"finesEnabled":  true,
			"shiftsEnabled": false,
		}
		req := MakeRequest(t, "POST", "/api/v1/clubs/invalid-id/settings", body, token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		req := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s/settings", club.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
	})
}
