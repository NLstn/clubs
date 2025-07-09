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

	t.Run("Join Club Via Link - Creates Join Request", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "user@example.com")

		// User joins via link
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/join", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		// Verify request was created
		var requests []models.JoinRequest
		err := database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&requests).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(requests))

		request := requests[0]
		assert.Equal(t, user.ID, request.UserID)
		assert.Equal(t, user.Email, request.Email)
		assert.Equal(t, club.ID, request.ClubID)
	})

	t.Run("Get Join Requests - Admin Can View", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner2@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, _ := CreateTestUser(t, "user2@example.com")

		// Create a join request
		err := club.CreateJoinRequest(user.ID, user.Email)
		assert.NoError(t, err)

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/joinRequests", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var joinRequests []map[string]interface{}
		ParseJSONResponse(t, rr, &joinRequests)
		assert.GreaterOrEqual(t, len(joinRequests), 1)
	})

	t.Run("Accept Join Request - Admin Can Approve", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner3@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, _ := CreateTestUser(t, "user3@example.com")

		// Create a join request
		err := club.CreateJoinRequest(user.ID, user.Email)
		assert.NoError(t, err)

		// Get the request ID
		var requests []models.JoinRequest
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&requests).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(requests))
		requestID := requests[0].ID

		// Admin accepts the request
		req := MakeRequest(t, "POST", "/api/v1/joinRequests/"+requestID+"/accept", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user was added to club and request was deleted
		isMember := club.IsMember(user)
		assert.True(t, isMember, "User should be a member after request acceptance")

		// Verify request was deleted
		err = database.Db.Where("id = ?", requestID).First(&requests[0]).Error
		assert.Error(t, err, "Request should be deleted after acceptance")
	})

	t.Run("Reject Join Request - Admin Can Reject", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner4@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, _ := CreateTestUser(t, "user4@example.com")

		// Create a join request
		err := club.CreateJoinRequest(user.ID, user.Email)
		assert.NoError(t, err)

		// Get the request ID
		var requests []models.JoinRequest
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&requests).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(requests))
		requestID := requests[0].ID

		// Admin rejects the request
		req := MakeRequest(t, "POST", "/api/v1/joinRequests/"+requestID+"/reject", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user was not added to club and request was deleted
		isMember := club.IsMember(user)
		assert.False(t, isMember, "User should not be a member after request rejection")

		// Verify request was deleted
		err = database.Db.Where("id = ?", requestID).First(&requests[0]).Error
		assert.Error(t, err, "Request should be deleted after rejection")
	})

	t.Run("Join Club Via Link - Already Member", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner5@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "user5@example.com")

		// Add user as member first
		err := club.AddMember(user.ID, "member")
		assert.NoError(t, err)

		// User tries to join via link
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/join", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusConflict, rr.Code)
		AssertContains(t, rr.Body.String(), "already a member")
	})

	t.Run("Get Join Requests - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner6@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		_, nonOwnerToken := CreateTestUser(t, "notowner@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/joinRequests", nil, nonOwnerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Invalid UUID Formats", func(t *testing.T) {
		_, token := CreateTestUser(t, "test@example.com")

		// Test GET endpoints with invalid UUID
		getEndpoints := []string{
			"/api/v1/clubs/invalid-uuid/joinRequests",
			"/api/v1/clubs/invalid-uuid/inviteLink",
			"/api/v1/clubs/invalid-uuid/info",
		}

		for _, endpoint := range getEndpoints {
			req := MakeRequest(t, "GET", endpoint, nil, token)
			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		}

		// Test POST endpoint with invalid UUID
		req := MakeRequest(t, "POST", "/api/v1/clubs/invalid-uuid/join", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})
}
