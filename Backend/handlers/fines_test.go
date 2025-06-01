package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinesEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get My Fines - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get My Fines - No Fines", func(t *testing.T) {
		user, token := CreateTestUser(t, "fines1@example.com")
		_ = user // We don't need the user for this test

		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var fines []map[string]interface{}
		ParseJSONResponse(t, rr, &fines)
		assert.Equal(t, 0, len(fines))
	})

	t.Run("Get My Fines - With Fines", func(t *testing.T) {
		user, token := CreateTestUser(t, "fines2@example.com")
		club := CreateTestClub(t, user, "Test Club for Fines")

		// Create both paid and unpaid fines
		unpaidFine := CreateTestFine(t, user, club, "Late arrival", 25.0, false)
		_ = CreateTestFine(t, user, club, "Missed meeting", 10.0, true) // paid fine

		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var fines []map[string]interface{}
		ParseJSONResponse(t, rr, &fines)
		
		// After fix, this should only return unpaid fines
		assert.Equal(t, 1, len(fines))

		// Check that we only have the unpaid fine
		fine := fines[0]
		assert.Equal(t, unpaidFine.ID, fine["id"])
		assert.Equal(t, false, fine["paid"])
		assert.Equal(t, "Late arrival", fine["reason"])
		assert.Equal(t, 25.0, fine["amount"])
		assert.Equal(t, club.Name, fine["clubName"])
		assert.Equal(t, user.Name, fine["createdByName"]) // Verify creator name is included
	})

	t.Run("Get My Fines - Only Unpaid Fines Expected", func(t *testing.T) {
		user, token := CreateTestUser(t, "fines3@example.com")
		club := CreateTestClub(t, user, "Test Club for Unpaid Fines")

		// Create both paid and unpaid fines
		unpaidFine1 := CreateTestFine(t, user, club, "Late arrival", 25.0, false)
		unpaidFine2 := CreateTestFine(t, user, club, "No show", 50.0, false)
		_ = CreateTestFine(t, user, club, "Missed meeting", 10.0, true) // paid fine

		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var fines []map[string]interface{}
		ParseJSONResponse(t, rr, &fines)
		
		// After fix, this should only return unpaid fines (2 fines)
		assert.Equal(t, 2, len(fines))

		// Verify we get only the unpaid fines
		var foundUnpaid1, foundUnpaid2 bool
		for _, fine := range fines {
			assert.Equal(t, false, fine["paid"], "All returned fines should be unpaid")
			if fine["id"] == unpaidFine1.ID {
				foundUnpaid1 = true
				assert.Equal(t, "Late arrival", fine["reason"])
				assert.Equal(t, 25.0, fine["amount"])
			}
			if fine["id"] == unpaidFine2.ID {
				foundUnpaid2 = true
				assert.Equal(t, "No show", fine["reason"])
				assert.Equal(t, 50.0, fine["amount"])
			}
		}
		assert.True(t, foundUnpaid1, "Should find first unpaid fine")
		assert.True(t, foundUnpaid2, "Should find second unpaid fine")
	})

	t.Run("Get Club Fines - Admin sees all fines", func(t *testing.T) {
		adminUser, adminToken := CreateTestUser(t, "admin@example.com")
		club := CreateTestClub(t, adminUser, "Test Club for Admin Fines")

		// Create another user to assign fines to
		memberUser, _ := CreateTestUser(t, "member@example.com")
		
		// Add member to the club
		CreateTestMember(t, memberUser, club, "member")

		// Create both paid and unpaid fines for the member
		unpaidFine := CreateTestFine(t, memberUser, club, "Late arrival", 25.0, false)
		paidFine := CreateTestFine(t, memberUser, club, "Missed meeting", 10.0, true)

		req := MakeRequest(t, "GET", "/api/v1/clubs/" + club.ID + "/fines", nil, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var fines []map[string]interface{}
		ParseJSONResponse(t, rr, &fines)
		
		// Admin should see ALL fines (both paid and unpaid)
		assert.Equal(t, 2, len(fines))

		// Verify we get both fines
		var foundUnpaid, foundPaid bool
		for _, fine := range fines {
			if fine["id"] == unpaidFine.ID {
				foundUnpaid = true
				assert.Equal(t, false, fine["paid"])
				assert.Equal(t, "Late arrival", fine["reason"])
				assert.Equal(t, 25.0, fine["amount"])
			}
			if fine["id"] == paidFine.ID {
				foundPaid = true
				assert.Equal(t, true, fine["paid"])
				assert.Equal(t, "Missed meeting", fine["reason"])
				assert.Equal(t, 10.0, fine["amount"])
			}
		}
		assert.True(t, foundUnpaid, "Admin should see unpaid fine")
		assert.True(t, foundPaid, "Admin should see paid fine")
	})
}