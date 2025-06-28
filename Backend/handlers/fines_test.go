package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/models"
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

	t.Run("Get My Fines - Test Caching with Multiple Creators", func(t *testing.T) {
		// Create test users
		userWithFines, token := CreateTestUser(t, "user_with_fines@example.com")
		creator1, _ := CreateTestUser(t, "creator1@example.com")
		creator2, _ := CreateTestUser(t, "creator2@example.com")

		// Create clubs
		club1 := CreateTestClub(t, creator1, "Club 1")
		club2 := CreateTestClub(t, creator2, "Club 2")

		// Add userWithFines as member to both clubs
		CreateTestMember(t, userWithFines, club1, "member")
		CreateTestMember(t, userWithFines, club2, "member")

		// Create multiple unpaid fines where some share the same creator
		// This tests the caching of creator names when multiple fines have the same creator
		fine1 := CreateTestFineWithCreator(t, userWithFines, club1, creator1, "Late to meeting", 10.0, false)
		fine2 := CreateTestFineWithCreator(t, userWithFines, club1, creator1, "Forgot materials", 15.0, false) // same creator as fine1
		fine3 := CreateTestFineWithCreator(t, userWithFines, club2, creator2, "No show", 25.0, false)
		fine4 := CreateTestFineWithCreator(t, userWithFines, club2, creator2, "Disruptive behavior", 20.0, false) // same creator as fine3

		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var fines []map[string]interface{}
		ParseJSONResponse(t, rr, &fines)

		// Should get all 4 unpaid fines
		assert.Equal(t, 4, len(fines))

		// Verify creator names are correctly populated from cache
		for _, fine := range fines {
			assert.NotEmpty(t, fine["createdByName"], "Creator name should be populated")
			assert.NotEmpty(t, fine["clubName"], "Club name should be populated")

			if fine["id"] == fine1.ID || fine["id"] == fine2.ID {
				assert.Equal(t, creator1.Name, fine["createdByName"], "Creator1 name should match")
				assert.Equal(t, club1.Name, fine["clubName"], "Club1 name should match")
			}
			if fine["id"] == fine3.ID || fine["id"] == fine4.ID {
				assert.Equal(t, creator2.Name, fine["createdByName"], "Creator2 name should match")
				assert.Equal(t, club2.Name, fine["clubName"], "Club2 name should match")
			}
		}
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

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/fines", nil, adminToken)
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

	t.Run("Get My Fines - Early Return for No Fines", func(t *testing.T) {
		// Create a user with no fines to verify early return behavior
		user, token := CreateTestUser(t, "no_fines_early_return@example.com")
		_ = user // User created but no fines assigned

		req := MakeRequest(t, "GET", "/api/v1/me/fines", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var fines []map[string]interface{}
		ParseJSONResponse(t, rr, &fines)

		// Verify empty array is returned
		assert.Equal(t, 0, len(fines))
		assert.NotNil(t, fines) // Should be empty array, not nil

		// Verify response content type is set correctly
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("Delete Fine - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "finedelete1@example.com")
		club := CreateTestClub(t, user, "Test Club")
		fine := CreateTestFine(t, user, club, "Test Fine", 25.0, false)

		// Try to delete with no token
		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/fines/"+fine.ID, nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Delete Fine - Forbidden (Not Admin)", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "fineowner@example.com")
		member, memberToken := CreateTestUser(t, "finemember@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		// Add member as regular member (not admin)
		CreateTestMember(t, member, club, "member")
		fine := CreateTestFine(t, member, club, "Test Fine", 25.0, false)

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/fines/"+fine.ID, nil, memberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		AssertContains(t, rr.Body.String(), "admin access required")
	})

	t.Run("Delete Fine - Success by Owner", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "fineowner2@example.com")
		member, _ := CreateTestUser(t, "finemember2@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		CreateTestMember(t, member, club, "member")
		fine := CreateTestFine(t, member, club, "Test Fine", 25.0, false)

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/fines/"+fine.ID, nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify fine is deleted
		var deletedFine models.Fine
		err := testDB.First(&deletedFine, "id = ?", fine.ID).Error
		assert.Error(t, err) // Should get "record not found" error
	})

	t.Run("Delete Fine - Success by Admin", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "fineowner3@example.com")
		admin, adminToken := CreateTestUser(t, "fineadmin@example.com")
		member, _ := CreateTestUser(t, "finemember3@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		CreateTestMember(t, admin, club, "admin")
		CreateTestMember(t, member, club, "member")
		fine := CreateTestFine(t, member, club, "Test Fine", 25.0, false)

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/fines/"+fine.ID, nil, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify fine is deleted
		var deletedFine models.Fine
		err := testDB.First(&deletedFine, "id = ?", fine.ID).Error
		assert.Error(t, err) // Should get "record not found" error
	})

	t.Run("Delete Fine - Invalid Fine ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "fineowner4@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/fines/invalid-id", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid fine ID format")
	})

	t.Run("Delete Fine - Club Not Found", func(t *testing.T) {
		_, ownerToken := CreateTestUser(t, "fineowner5@example.com")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/invalid-club-id/fines/invalid-fine-id", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Create Fine - Response Uses CamelCase Fields", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "fineowner6@example.com")
		member, _ := CreateTestUser(t, "finecamel@example.com")
		club := CreateTestClub(t, owner, "Camel Club")
		CreateTestMember(t, member, club, "member")

		payload := map[string]interface{}{
			"userId": member.ID,
			"reason": "Late arrival",
			"amount": 15.0,
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fines", payload, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var resp map[string]interface{}
		ParseJSONResponse(t, rr, &resp)
		assert.Contains(t, resp, "createdAt")
		assert.Contains(t, resp, "updatedAt")
		assert.NotContains(t, resp, "created_at")
		assert.NotContains(t, resp, "updated_at")
	})
}
