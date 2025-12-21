package models_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/NLstn/clubs/odata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClubSettingsCreatedWithClub(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Set up OData service for testing
	service, err := odata.NewService(handlers.GetDB())
	require.NoError(t, err, "Failed to create OData service")

	odataV2Mux := http.NewServeMux()
	service.RegisterCustomHandlers(odataV2Mux)
	odataV2Mux.Handle("/", service)
	handler := http.StripPrefix("/api/v2", handlers.CompositeAuthMiddleware(odataV2Mux))

	t.Run("settings created automatically when club is created via OData", func(t *testing.T) {
		user, token := handlers.CreateTestUser(t, "clubsettingstest@example.com")

		// Create a club via OData API
		clubData := map[string]interface{}{
			"Name":        "Test Club",
			"Description": "Test Description",
		}
		body, err := json.Marshal(clubData)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v2/Clubs", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoError(t, err)
		clubID := created["ID"].(string)
		assert.NotEmpty(t, clubID)

		// Verify that settings were created automatically
		settings, err := models.GetClubSettings(clubID)
		assert.NoError(t, err)
		assert.NotEmpty(t, settings.ID)
		assert.Equal(t, clubID, settings.ClubID)

		// Verify all settings are disabled by default
		assert.False(t, settings.FinesEnabled, "FinesEnabled should be false by default")
		assert.False(t, settings.ShiftsEnabled, "ShiftsEnabled should be false by default")
		assert.False(t, settings.TeamsEnabled, "TeamsEnabled should be false by default")
		assert.False(t, settings.MembersListVisible, "MembersListVisible should be false by default")
		assert.False(t, settings.DiscoverableByNonMembers, "DiscoverableByNonMembers should be false by default")

		// Verify audit fields
		assert.Equal(t, user.ID, settings.CreatedBy)
		assert.Equal(t, user.ID, settings.UpdatedBy)
		assert.NotZero(t, settings.CreatedAt)
		assert.NotZero(t, settings.UpdatedAt)
	})

	t.Run("test helper creates settings with all disabled", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "testhelperclub@example.com")

		// Create a club using test helper
		club := handlers.CreateTestClub(t, user, "Helper Test Club")

		// Verify settings exist
		settings, err := models.GetClubSettings(club.ID)
		assert.NoError(t, err)
		assert.Equal(t, club.ID, settings.ClubID)

		// Verify all settings are disabled
		assert.False(t, settings.FinesEnabled)
		assert.False(t, settings.ShiftsEnabled)
		assert.False(t, settings.TeamsEnabled)
		assert.False(t, settings.MembersListVisible)
		assert.False(t, settings.DiscoverableByNonMembers)
	})
}

func TestCreateDefaultClubSettings(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("creates settings with all features disabled", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "defaultsettings@example.com")
		club := handlers.CreateTestClub(t, user, "Default Settings Club")

		// Delete the settings created by test helper and ensure deletion succeeds
		db := handlers.GetDB().Exec("DELETE FROM club_settings WHERE club_id = ?", club.ID)
		assert.NoError(t, db.Error)
		assert.NotZero(t, db.RowsAffected, "expected at least one club_settings row to be deleted")

		// Create default settings with proper user ID for audit trail
		settings, err := models.CreateDefaultClubSettings(club.ID, user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, settings.ID)

		// Verify all settings are disabled
		assert.False(t, settings.FinesEnabled)
		assert.False(t, settings.ShiftsEnabled)
		assert.False(t, settings.TeamsEnabled)
		assert.False(t, settings.MembersListVisible)
		assert.False(t, settings.DiscoverableByNonMembers)
		
		// Verify audit fields are correct
		assert.Equal(t, user.ID, settings.CreatedBy)
		assert.Equal(t, user.ID, settings.UpdatedBy)
	})
}
