package handlers

import (
	"testing"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersByIDsEdgeCases(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	t.Run("GetUsersByIDs with empty slice", func(t *testing.T) {
		users, err := models.GetUsersByIDs([]string{})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(users))
		// Note: GORM returns nil slice when no results, which is normal behavior
	})

	t.Run("GetUsersByIDs with slice containing empty string", func(t *testing.T) {
		users, err := models.GetUsersByIDs([]string{""})
		// This might cause the error mentioned in the issue
		if err != nil {
			t.Logf("Error with empty string ID: %v", err)
		}
		assert.Equal(t, 0, len(users))
	})

	t.Run("GetClubsByIDs with empty slice", func(t *testing.T) {
		clubs, err := models.GetClubsByIDs([]string{})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(clubs))
		// Note: GORM returns nil slice when no results, which is normal behavior
	})

	t.Run("GetClubsByIDs with slice containing empty string", func(t *testing.T) {
		clubs, err := models.GetClubsByIDs([]string{""})
		// This might cause the error mentioned in the issue
		if err != nil {
			t.Logf("Error with empty string ID: %v", err)
		}
		assert.Equal(t, 0, len(clubs))
	})

	t.Run("Mixed valid and empty IDs", func(t *testing.T) {
		// Test the scenario where we might have mixed valid and empty IDs
		testUser, _ := CreateTestUser(t, "test_empty_ids@example.com")
		
		// Test with mixed IDs including empty strings
		users, err := models.GetUsersByIDs([]string{testUser.ID, "", "invalid-id"})
		assert.NoError(t, err)
		// Should find only the valid user
		assert.Equal(t, 1, len(users))
		if len(users) > 0 {
			assert.Equal(t, testUser.ID, users[0].ID)
		}
	})
}