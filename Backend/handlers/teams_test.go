package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestTeamEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get Club Teams - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "teamtest@example.com")
		club := CreateTestClub(t, user, "Test Club")

		req := MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/teams", club.ID), nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Club Teams - Authorized Member", func(t *testing.T) {
		user, token := CreateTestUser(t, "teamtest2@example.com")
		club := CreateTestClub(t, user, "Test Club")

		req := MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/teams", club.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var teams []map[string]interface{}
		ParseJSONResponse(t, rr, &teams)
		assert.Len(t, teams, 0) // No teams initially
	})

	t.Run("Create Team - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "teamtest3@example.com")
		club := CreateTestClub(t, user, "Test Club")

		teamData := map[string]string{
			"name":        "Development Team",
			"description": "Team for development tasks",
		}

		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/teams", club.ID), teamData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Create Team - Club Admin", func(t *testing.T) {
		user, token := CreateTestUser(t, "teamtest4@example.com")
		club := CreateTestClub(t, user, "Test Club")

		teamData := map[string]string{
			"name":        "Development Team",
			"description": "Team for development tasks",
		}

		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/teams", club.ID), teamData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		var team map[string]interface{}
		ParseJSONResponse(t, rr, &team)
		assert.Equal(t, "Development Team", team["name"])
		assert.Equal(t, "Team for development tasks", team["description"])
		assert.Equal(t, club.ID, team["clubId"])
	})

	t.Run("Create Team - Regular Member Cannot Create", func(t *testing.T) {
		admin, _ := CreateTestUser(t, "teamadmin@example.com")
		club := CreateTestClub(t, admin, "Test Club")

		member, memberToken := CreateTestUser(t, "teammember@example.com")
		// Add member to club as regular member
		err := club.AddMember(member.ID, "member")
		assert.NoError(t, err)

		teamData := map[string]string{
			"name":        "Member Team",
			"description": "Team created by member",
		}

		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/teams", club.ID), teamData, memberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Create Team - Invalid Data", func(t *testing.T) {
		user, token := CreateTestUser(t, "teamtest5@example.com")
		club := CreateTestClub(t, user, "Test Club")

		// Test with empty name
		teamData := map[string]string{
			"name":        "",
			"description": "Team with no name",
		}

		req := MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/teams", club.ID), teamData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Get Team Details", func(t *testing.T) {
		user, token := CreateTestUser(t, "teamtest6@example.com")
		club := CreateTestClub(t, user, "Test Club")

		// Create a team first
		team, err := club.CreateTeam("Test Team", "A test team", user.ID)
		assert.NoError(t, err)

		req := MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/teams/%s", club.ID, team.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var responseTeam map[string]interface{}
		ParseJSONResponse(t, rr, &responseTeam)
		assert.Equal(t, "Test Team", responseTeam["name"])
		assert.Equal(t, "A test team", responseTeam["description"])
	})

	t.Run("Update Team", func(t *testing.T) {
		user, token := CreateTestUser(t, "teamtest7@example.com")
		club := CreateTestClub(t, user, "Test Club")

		// Create a team first
		team, err := club.CreateTeam("Original Team", "Original description", user.ID)
		assert.NoError(t, err)

		updateData := map[string]string{
			"name":        "Updated Team",
			"description": "Updated description",
		}

		req := MakeRequest(t, "PUT", fmt.Sprintf("/api/v1/clubs/%s/teams/%s", club.ID, team.ID), updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify the update
		updatedTeam, err := models.GetTeamByID(team.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Team", updatedTeam.Name)
		assert.Equal(t, "Updated description", updatedTeam.Description)
	})

	t.Run("Delete Team", func(t *testing.T) {
		user, token := CreateTestUser(t, "teamtest8@example.com")
		club := CreateTestClub(t, user, "Test Club")

		// Create a team first
		team, err := club.CreateTeam("Team to Delete", "This team will be deleted", user.ID)
		assert.NoError(t, err)

		req := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s/teams/%s", club.ID, team.ID), nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify the team is soft deleted
		_, err = models.GetTeamByID(team.ID)
		assert.Error(t, err) // Should not be found since GetTeamByID only returns non-deleted teams
	})

	t.Run("Team Members Management", func(t *testing.T) {
		admin, adminToken := CreateTestUser(t, "teamadmin2@example.com")
		club := CreateTestClub(t, admin, "Test Club")

		member1, _ := CreateTestUser(t, "teammember1@example.com")
		member2, _ := CreateTestUser(t, "teammember2@example.com")

		// Add members to club
		err := club.AddMember(member1.ID, "member")
		assert.NoError(t, err)
		err = club.AddMember(member2.ID, "member")
		assert.NoError(t, err)

		// Create a team
		team, err := club.CreateTeam("Test Team", "Team for member testing", admin.ID)
		assert.NoError(t, err)

		// Test getting team members (should be empty initially)
		req := MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/teams/%s/members", club.ID, team.ID), nil, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var members []map[string]interface{}
		ParseJSONResponse(t, rr, &members)
		assert.Len(t, members, 0)

		// Add member1 to team as member
		addMemberData := map[string]string{
			"userId": member1.ID,
			"role":   "member",
		}

		req = MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/teams/%s/members", club.ID, team.ID), addMemberData, adminToken)
		rr = ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		// Add member2 to team as admin
		addAdminData := map[string]string{
			"userId": member2.ID,
			"role":   "admin",
		}

		req = MakeRequest(t, "POST", fmt.Sprintf("/api/v1/clubs/%s/teams/%s/members", club.ID, team.ID), addAdminData, adminToken)
		rr = ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		// Verify team members
		req = MakeRequest(t, "GET", fmt.Sprintf("/api/v1/clubs/%s/teams/%s/members", club.ID, team.ID), nil, adminToken)
		rr = ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		ParseJSONResponse(t, rr, &members)
		assert.Len(t, members, 2)

		// Find the members in the response
		var memberRole, adminRole string
		for _, member := range members {
			if member["userId"].(string) == member1.ID {
				memberRole = member["role"].(string)
			}
			if member["userId"].(string) == member2.ID {
				adminRole = member["role"].(string)
			}
		}
		assert.Equal(t, "member", memberRole)
		assert.Equal(t, "admin", adminRole)
	})

	t.Run("Team Member Role Update", func(t *testing.T) {
		admin, adminToken := CreateTestUser(t, "teamadmin3@example.com")
		club := CreateTestClub(t, admin, "Test Club")

		member, _ := CreateTestUser(t, "teammember3@example.com")
		err := club.AddMember(member.ID, "member")
		assert.NoError(t, err)

		// Create team and add member
		team, err := club.CreateTeam("Role Test Team", "Team for role testing", admin.ID)
		assert.NoError(t, err)

		err = team.AddMember(member.ID, "member", admin.ID)
		assert.NoError(t, err)

		// Get the team member record
		teamMembers, err := team.GetTeamMembersWithUserInfo()
		assert.NoError(t, err)
		assert.Len(t, teamMembers, 1)

		memberID := teamMembers[0]["id"].(string)

		// Promote member to admin
		updateRoleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", fmt.Sprintf("/api/v1/clubs/%s/teams/%s/members/%s", club.ID, team.ID, memberID), updateRoleData, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify role change
		updatedMembers, err := team.GetTeamMembersWithUserInfo()
		assert.NoError(t, err)
		assert.Equal(t, "admin", updatedMembers[0]["role"].(string))
	})

	t.Run("Remove Team Member", func(t *testing.T) {
		admin, adminToken := CreateTestUser(t, "teamadmin4@example.com")
		club := CreateTestClub(t, admin, "Test Club")

		member, _ := CreateTestUser(t, "teammember4@example.com")
		err := club.AddMember(member.ID, "member")
		assert.NoError(t, err)

		// Create team and add member
		team, err := club.CreateTeam("Remove Test Team", "Team for removal testing", admin.ID)
		assert.NoError(t, err)

		err = team.AddMember(member.ID, "member", admin.ID)
		assert.NoError(t, err)

		// Get the team member record
		teamMembers, err := team.GetTeamMembersWithUserInfo()
		assert.NoError(t, err)
		assert.Len(t, teamMembers, 1)

		memberID := teamMembers[0]["id"].(string)

		// Remove member from team
		req := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/clubs/%s/teams/%s/members/%s", club.ID, team.ID, memberID), nil, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify member is removed
		updatedMembers, err := team.GetTeamMembersWithUserInfo()
		assert.NoError(t, err)
		assert.Len(t, updatedMembers, 0)
	})
}
