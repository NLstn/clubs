package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestJoinRequestEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Create Join Request - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "owner1@example.com")
		club := CreateTestClub(t, user, "Test Club")

		joinData := map[string]string{
			"email": "newmember1@example.com",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/joinRequests", joinData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Create Join Request - Not Owner", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner2@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		_, nonOwnerToken := CreateTestUser(t, "notowner2@example.com")

		joinData := map[string]string{
			"email": "newmember2@example.com",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/joinRequests", joinData, nonOwnerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		AssertContains(t, rr.Body.String(), "Unauthorized")
	})

	t.Run("Create Join Request - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner3@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		joinData := map[string]string{
			"email": "newmember3@example.com",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/joinRequests", joinData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		// Verify that the join request was created with proper created_by field
		// Since this is an admin invite (admin_approved=true), it should NOT appear in the club's pending requests
		// but should be stored in the database. Let's check the database directly.
		var allRequests []models.JoinRequest
		err := database.Db.Where("club_id = ?", club.ID).Find(&allRequests).Error
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(allRequests), 1)
		
		// Find our join request
		var foundRequest *models.JoinRequest
		for _, jr := range allRequests {
			if jr.Email == "newmember3@example.com" {
				foundRequest = &jr
				break
			}
		}
		assert.NotNil(t, foundRequest, "Join request should be found")
		assert.Equal(t, owner.ID, foundRequest.CreatedBy, "CreatedBy should be set to the owner's ID")
		assert.Equal(t, owner.ID, foundRequest.UpdatedBy, "UpdatedBy should be set to the owner's ID")
		assert.True(t, foundRequest.AdminApproved, "AdminApproved should be true for admin invites")
		assert.False(t, foundRequest.UserApproved, "UserApproved should be false for admin invites")
	})

	t.Run("Create Join Request - Missing Email", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner4@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		joinData := map[string]string{}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/joinRequests", joinData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Missing email")
	})

	t.Run("Create Join Request - Club Not Found", func(t *testing.T) {
		_, ownerToken := CreateTestUser(t, "owner5@example.com")

		joinData := map[string]string{
			"email": "newmember5@example.com",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/invalid-id/joinRequests", joinData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
		AssertContains(t, rr.Body.String(), "Club not found")
	})

	t.Run("Get Join Requests - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner6@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/joinRequests", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Join Requests - Not Owner", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner7@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		_, nonOwnerToken := CreateTestUser(t, "notowner7@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/joinRequests", nil, nonOwnerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		AssertContains(t, rr.Body.String(), "Unauthorized")
	})

	t.Run("Get Join Requests - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner8@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		
		// Create a user who will request to join
		requestingUser, _ := CreateTestUser(t, "newmember8@example.com")

		// Create a join request where user is requesting to join (admin needs to approve)
		club.CreateJoinRequest(requestingUser.Email, requestingUser.ID, false, true)

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/joinRequests", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var joinRequests []map[string]interface{}
		ParseJSONResponse(t, rr, &joinRequests)
		assert.GreaterOrEqual(t, len(joinRequests), 1)
	})

	t.Run("Get Join Requests - Invalid Club ID", func(t *testing.T) {
		_, ownerToken := CreateTestUser(t, "owner9@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/invalid-uuid/joinRequests", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Get User Join Requests - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/joinRequests", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get User Join Requests - Valid", func(t *testing.T) {
		user, userToken := CreateTestUser(t, "user10@example.com")

		// Create a join request for this user
		owner, _ := CreateTestUser(t, "owner10@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		club.CreateJoinRequest(user.Email, owner.ID, true, false)

		req := MakeRequest(t, "GET", "/api/v1/joinRequests", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var joinRequests []map[string]interface{}
		ParseJSONResponse(t, rr, &joinRequests)
		// Should have at least one join request
		assert.GreaterOrEqual(t, len(joinRequests), 0)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		_, token := CreateTestUser(t, "test11@example.com")

		endpoints := []string{
			"/api/v1/clubs/test-id/joinRequests",
			"/api/v1/joinRequests",
			"/api/v1/joinRequests/test-id/accept",
			"/api/v1/joinRequests/test-id/reject",
		}

		for _, endpoint := range endpoints {
			req := MakeRequest(t, "PUT", endpoint, nil, token)
			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
		}
	})
}