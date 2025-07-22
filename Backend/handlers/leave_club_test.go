package handlers

import (
	"net/http"
	"strings"
	"testing"
)

func TestHandleLeaveClub(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	// Create test data
	user1, token1 := CreateTestUser(t, "user1@example.com")
	user2, token2 := CreateTestUser(t, "user2@example.com")
	club := CreateTestClub(t, user1, "Test Club")

	// Add user2 as a member
	CreateTestMember(t, user2, club, "member")

	t.Run("Member can leave club", func(t *testing.T) {
		// Create request
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/leave", nil, token2)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user is no longer a member
		if club.IsMember(user2) {
			t.Error("User should no longer be a member after leaving")
		}
	})

	t.Run("Owner cannot leave if they are the last owner", func(t *testing.T) {
		// Create request for owner to leave
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/leave", nil, token1)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)

		// Check error message
		if !strings.Contains(rr.Body.String(), "last owner") {
			t.Error("Expected error message about being the last owner")
		}

		// Verify owner is still a member
		if !club.IsMember(user1) {
			t.Error("Owner should still be a member after failed leave attempt")
		}
	})

	t.Run("Non-member cannot leave club", func(t *testing.T) {
		// Create a new user who is not a member
		_, token3 := CreateTestUser(t, "user3@example.com")

		// Create request
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/leave", nil, token3)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)

		// Check error message
		if !strings.Contains(rr.Body.String(), "not a member") {
			t.Error("Expected error message about not being a member")
		}
	})

	t.Run("Owner can leave if there are other owners", func(t *testing.T) {
		// Add user2 back as a member and promote to owner
		CreateTestMember(t, user2, club, "member")

		members, err := club.GetClubMembers()
		if err != nil {
			t.Fatalf("Failed to get club members: %v", err)
		}

		var user2MemberID string
		for _, member := range members {
			if member.UserID == user2.ID {
				user2MemberID = member.ID
				break
			}
		}

		err = club.UpdateMemberRole(user1, user2MemberID, "owner")
		if err != nil {
			t.Fatalf("Failed to promote user to owner: %v", err)
		}

		// Now user1 should be able to leave since user2 is also an owner
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/leave", nil, token1)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify user1 is no longer a member
		if club.IsMember(user1) {
			t.Error("User should no longer be a member after leaving")
		}

		// Verify user2 is still an owner
		role, err := club.GetMemberRole(user2)
		if err != nil || role != "owner" {
			t.Error("User2 should still be an owner")
		}
	})
}

func TestHandleLeaveClubInvalidClub(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()
	_, token := CreateTestUser(t, "user@example.com")

	t.Run("Invalid club ID format", func(t *testing.T) {
		req := MakeRequest(t, "POST", "/api/v1/clubs/invalid-id/leave", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Non-existent club", func(t *testing.T) {
		req := MakeRequest(t, "POST", "/api/v1/clubs/00000000-0000-0000-0000-000000000000/leave", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
	})
}
