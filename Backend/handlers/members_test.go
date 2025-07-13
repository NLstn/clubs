package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestMemberEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get Club Members - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner1@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/members", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Club Members - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner2@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		// Add another member to the club
		member, _ := CreateTestUser(t, "member2@example.com")
		CreateTestMember(t, member, club, "member")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/members", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var members []map[string]interface{}
		ParseJSONResponse(t, rr, &members)
		assert.GreaterOrEqual(t, len(members), 2) // At least owner and member
	})

	t.Run("Get Club Members - Invalid Club ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "user3@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/invalid-uuid/members", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Get Club Members - Club Not Found", func(t *testing.T) {
		_, token := CreateTestUser(t, "user4@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/550e8400-e29b-41d4-a716-446655440000/members", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
		AssertContains(t, rr.Body.String(), "Club not found")
	})

	t.Run("Delete Club Member - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner5@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, _ := CreateTestUser(t, "member5@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Delete Club Member - Not Owner", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner6@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, memberToken := CreateTestUser(t, "member6@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, nil, memberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		AssertContains(t, rr.Body.String(), "Unauthorized")
	})

	t.Run("Delete Club Member - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner7@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, _ := CreateTestUser(t, "member7@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Delete Club Member - Invalid Club ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "user8@example.com")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/invalid-uuid/members/test-id", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Delete Club Member - Invalid Member ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner9@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/members/invalid-uuid", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid member ID format")
	})

	t.Run("Delete Club Member - Member Not Found", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner10@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/members/550e8400-e29b-41d4-a716-446655440000", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
		AssertContains(t, rr.Body.String(), "Member not found")
	})

	t.Run("Update Member Role - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner11@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, _ := CreateTestUser(t, "member11@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		roleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, roleData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Update Member Role - Not Admin", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner12@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member1, member1Token := CreateTestUser(t, "member12a@example.com")
		CreateTestMember(t, member1, club, "member")
		member2, _ := CreateTestUser(t, "member12b@example.com")
		member2Record := CreateTestMember(t, member2, club, "member")

		roleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+member2Record.ID, roleData, member1Token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		AssertContains(t, rr.Body.String(), "Unauthorized")
	})

	t.Run("Update Member Role - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner13@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, _ := CreateTestUser(t, "member13@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		roleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, roleData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Update Member Role - Notification Sent", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner-notif@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, _ := CreateTestUser(t, "member-notif@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		roleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, roleData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify notification was created for the member whose role changed
		var notifications []models.Notification
		err := database.Db.Where("user_id = ? AND type = ?", member.ID, "role_changed").Find(&notifications).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, len(notifications))

		notification := notifications[0]
		assert.Equal(t, "Role Updated in Test Club", notification.Title)
		assert.Contains(t, notification.Message, "Your role in Test Club has been changed from member to admin")
		assert.Equal(t, club.ID, *notification.ClubID)
		assert.False(t, notification.Read)
	})

	t.Run("Update Member Role - Invalid Role", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner14@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, _ := CreateTestUser(t, "member14@example.com")
		memberRecord := CreateTestMember(t, member, club, "member")

		roleData := map[string]string{
			"role": "invalid-role",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+memberRecord.ID, roleData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid role")
	})

	t.Run("Check Admin Rights - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner15@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/isAdmin", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		ParseJSONResponse(t, rr, &response)
		assert.Equal(t, true, response["isAdmin"])
		assert.Equal(t, true, response["isOwner"])
	})

	t.Run("Check Admin Rights - Admin (Non-Owner)", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner15b@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		admin, adminToken := CreateTestUser(t, "admin15b@example.com")
		CreateTestMember(t, admin, club, "admin")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/isAdmin", nil, adminToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		ParseJSONResponse(t, rr, &response)
		assert.Equal(t, true, response["isAdmin"])
		assert.Equal(t, false, response["isOwner"])
	})

	t.Run("Check Admin Rights - Non Admin", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner16@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, memberToken := CreateTestUser(t, "member16@example.com")
		CreateTestMember(t, member, club, "member")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/isAdmin", nil, memberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		ParseJSONResponse(t, rr, &response)
		assert.Equal(t, false, response["isAdmin"])
		assert.Equal(t, false, response["isOwner"])
	})

	t.Run("Check Admin Rights - Club Not Found", func(t *testing.T) {
		_, token := CreateTestUser(t, "user17@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/550e8400-e29b-41d4-a716-446655440000/isAdmin", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
		AssertContains(t, rr.Body.String(), "Club not found")
	})

	t.Run("Prevent Last Owner From Demoting Themselves", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "single-owner@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		
		// Get the owner's member record that was created by CreateTestClub
		var ownerMember models.Member
		err := database.Db.Where("club_id = ? AND user_id = ? AND role = ?", club.ID, owner.ID, "owner").First(&ownerMember).Error
		assert.NoError(t, err)

		roleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+ownerMember.ID, roleData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Allow Owner To Demote Themselves When Multiple Owners Exist", func(t *testing.T) {
		owner1, owner1Token := CreateTestUser(t, "owner1-multiple@example.com")
		club := CreateTestClub(t, owner1, "Test Club")
		
		// Get the first owner's member record that was created by CreateTestClub
		var owner1Member models.Member
		err := database.Db.Where("club_id = ? AND user_id = ? AND role = ?", club.ID, owner1.ID, "owner").First(&owner1Member).Error
		assert.NoError(t, err)

		owner2, _ := CreateTestUser(t, "owner2-multiple@example.com")
		CreateTestMember(t, owner2, club, "owner")

		roleData := map[string]string{
			"role": "admin",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID+"/members/"+owner1Member.ID, roleData, owner1Token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Get Owner Count - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner-count@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		// Add another owner
		owner2, _ := CreateTestUser(t, "owner2-count@example.com")
		CreateTestMember(t, owner2, club, "owner")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/ownerCount", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		ParseJSONResponse(t, rr, &response)
		assert.Equal(t, float64(2), response["ownerCount"]) // JSON numbers are float64
	})

	t.Run("Get Owner Count - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner-count-unauth@example.com")
		club := CreateTestClub(t, owner, "Test Club")
		member, memberToken := CreateTestUser(t, "member-count-unauth@example.com")
		CreateTestMember(t, member, club, "member")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/ownerCount", nil, memberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
	})

	// ...existing code...
}
