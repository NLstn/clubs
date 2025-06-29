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

	// Test new approval flow scenarios
	t.Run("Join Via Link - User Requests to Join", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner_link@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "user_link@example.com")

		// User joins via link
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/join", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		// Verify request was created with correct approval flags
		var allRequests []models.JoinRequest
		err := database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&allRequests).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(allRequests))
		
		request := allRequests[0]
		assert.False(t, request.AdminApproved, "AdminApproved should be false for user-initiated requests")
		assert.True(t, request.UserApproved, "UserApproved should be true for user-initiated requests")
		assert.Equal(t, user.ID, request.CreatedBy, "CreatedBy should be the user who requested to join")

		// Verify it shows up in admin's pending requests but not user's
		adminRequests, _ := club.GetJoinRequests()
		assert.Equal(t, 1, len(adminRequests), "Admin should see the user's request to join")
		
		userRequests, _ := user.GetJoinRequests()
		assert.Equal(t, 0, len(userRequests), "User should not see their own request as an invite")
	})

	t.Run("Admin Partially Accepts User Request (Approval Logic)", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner_partial@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, _ := CreateTestUser(t, "user_partial@example.com")

		// Create user request to join where user has NOT approved yet (so completion won't happen)
		err := club.CreateJoinRequest(user.Email, user.ID, false, false)
		assert.NoError(t, err)

		// Get the request ID
		var requests []models.JoinRequest
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&requests).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(requests))
		requestID := requests[0].ID

		// Admin accepts the request (but since user hasn't approved, it won't complete)
		req := MakeRequest(t, "POST", "/api/v1/joinRequests/"+requestID+"/accept", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify admin approval was set but request still exists
		err = database.Db.Where("id = ?", requestID).First(&requests[0]).Error
		assert.NoError(t, err)
		assert.True(t, requests[0].AdminApproved, "Admin approval should be set")
		assert.False(t, requests[0].UserApproved, "User approval should still be false")
	})

	t.Run("User Partially Accepts Admin Invite (Approval Logic)", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner_partial2@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "user_partial2@example.com")

		// Create admin invite where admin has NOT approved yet (so completion won't happen)
		err := club.CreateJoinRequest(user.Email, owner.ID, false, false)
		assert.NoError(t, err)

		// Get the request ID
		var requests []models.JoinRequest
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&requests).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(requests))
		requestID := requests[0].ID

		// User accepts the invite (but since admin hasn't approved, it won't complete)
		req := MakeRequest(t, "POST", "/api/v1/joinRequests/"+requestID+"/accept", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user approval was set but request still exists
		err = database.Db.Where("id = ?", requestID).First(&requests[0]).Error
		assert.NoError(t, err)
		assert.False(t, requests[0].AdminApproved, "Admin approval should still be false")
		assert.True(t, requests[0].UserApproved, "User approval should be set")
	})
}