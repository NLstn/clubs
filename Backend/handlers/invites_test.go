package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestInviteEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Create Invite - Admin Can Invite", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		inviteData := map[string]string{
			"email": "newmember@example.com",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/invites", inviteData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		// Verify invite was created
		var invites []models.Invite
		err := database.Db.Where("club_id = ? AND email = ?", club.ID, "newmember@example.com").Find(&invites).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(invites))

		invite := invites[0]
		assert.Equal(t, "newmember@example.com", invite.Email)
		assert.Equal(t, club.ID, invite.ClubID)
		assert.Equal(t, owner.ID, invite.CreatedBy)
	})

	t.Run("Get Club Invites - Admin Can View", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner2@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		// Create an invite
		err := club.CreateInvite("test@example.com", owner.ID)
		assert.NoError(t, err)

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/invites", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var invites []map[string]interface{}
		ParseJSONResponse(t, rr, &invites)
		assert.GreaterOrEqual(t, len(invites), 1)
	})

	t.Run("Get User Invites", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner3@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "invitee@example.com")

		// Create an invite for the user
		err := club.CreateInvite(user.Email, owner.ID)
		assert.NoError(t, err)

		req := MakeRequest(t, "GET", "/api/v1/invites", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var invites []map[string]interface{}
		ParseJSONResponse(t, rr, &invites)
		assert.GreaterOrEqual(t, len(invites), 1)
	})

	t.Run("Accept Invite - User Can Accept", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner4@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "invitee2@example.com")

		// Create an invite for the user
		err := club.CreateInvite(user.Email, owner.ID)
		assert.NoError(t, err)

		// Get the invite ID
		var invites []models.Invite
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&invites).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(invites))
		inviteID := invites[0].ID

		// User accepts the invite
		req := MakeRequest(t, "POST", "/api/v1/invites/"+inviteID+"/accept", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user was added to club and invite was deleted
		isMember := club.IsMember(user)
		assert.True(t, isMember, "User should be a member after accepting invite")

		// Verify invite was deleted
		err = database.Db.Where("id = ?", inviteID).First(&invites[0]).Error
		assert.Error(t, err, "Invite should be deleted after acceptance")
	})

	t.Run("Reject Invite - User Can Reject", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner5@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user, userToken := CreateTestUser(t, "invitee3@example.com")

		// Create an invite for the user
		err := club.CreateInvite(user.Email, owner.ID)
		assert.NoError(t, err)

		// Get the invite ID
		var invites []models.Invite
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user.Email).Find(&invites).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(invites))
		inviteID := invites[0].ID

		// User rejects the invite
		req := MakeRequest(t, "POST", "/api/v1/invites/"+inviteID+"/reject", nil, userToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user was not added to club and invite was deleted
		isMember := club.IsMember(user)
		assert.False(t, isMember, "User should not be a member after rejecting invite")

		// Verify invite was deleted
		err = database.Db.Where("id = ?", inviteID).First(&invites[0]).Error
		assert.Error(t, err, "Invite should be deleted after rejection")
	})

	t.Run("Create Invite - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner6@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		_, nonOwnerToken := CreateTestUser(t, "notowner@example.com")

		inviteData := map[string]string{
			"email": "newmember@example.com",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/invites", inviteData, nonOwnerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Accept Invite - Wrong User", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner7@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		user1, _ := CreateTestUser(t, "invitee4@example.com")
		_, user2Token := CreateTestUser(t, "different@example.com")

		// Create an invite for user1
		err := club.CreateInvite(user1.Email, owner.ID)
		assert.NoError(t, err)

		// Get the invite ID
		var invites []models.Invite
		err = database.Db.Where("club_id = ? AND email = ?", club.ID, user1.Email).Find(&invites).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(invites))
		inviteID := invites[0].ID

		// User2 tries to accept user1's invite
		req := MakeRequest(t, "POST", "/api/v1/invites/"+inviteID+"/accept", nil, user2Token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid UUID Formats", func(t *testing.T) {
		_, token := CreateTestUser(t, "test@example.com")

		endpoints := []string{
			"/api/v1/clubs/invalid-uuid/invites",
			"/api/v1/invites/invalid-uuid/accept",
			"/api/v1/invites/invalid-uuid/reject",
		}

		for _, endpoint := range endpoints {
			req := MakeRequest(t, "POST", endpoint, nil, token)
			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		}
	})
}
