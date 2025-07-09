package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestHardDeleteClub(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Hard delete non-soft-deleted club should fail", func(t *testing.T) {
		user, token := CreateTestUser(t, "testuser@example.com")
		club := CreateTestClub(t, user, "Test Club")

		// Try to hard delete a club that's not soft deleted (should fail)
		req := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s/hard-delete", club.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Club must be soft deleted before permanent deletion")
	})

	t.Run("Hard delete after soft delete should succeed", func(t *testing.T) {
		user, token := CreateTestUser(t, "testuser2@example.com")
		club := CreateTestClub(t, user, "Test Club 2")

		// First soft delete the club
		req := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s", club.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Now try to hard delete
		req = MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s/hard-delete", club.ID), nil, token)
		rr = ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify the club is completely gone from the database
		var dbClub models.Club
		err := testDB.Unscoped().First(&dbClub, "id = ?", club.ID).Error
		assert.Error(t, err) // Should not be found
	})
}

func TestHardDeleteClubUnauthorized(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Hard delete as non-owner should fail", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner@example.com")
		_, nonOwnerToken := CreateTestUser(t, "nonowner@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		// First soft delete the club as owner
		ownerToken, err := auth.GenerateAccessToken(owner.ID)
		if err != nil {
			t.Fatalf("Failed to generate owner token: %v", err)
		}

		req := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s", club.ID), nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Try to hard delete as non-owner (should fail)
		req = MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s/hard-delete", club.ID), nil, nonOwnerToken)
		rr = ExecuteRequest(t, handler, req)

		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized - only owners can permanently delete clubs")
	})
}
