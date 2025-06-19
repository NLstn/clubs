package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestClubEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get All Clubs - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/clubs", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get All Clubs - Authorized", func(t *testing.T) {
		user, token := CreateTestUser(t, "clubtest@example.com")
		CreateTestClub(t, user, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var clubs []map[string]interface{}
		ParseJSONResponse(t, rr, &clubs)
		assert.Len(t, clubs, 1)
		assert.Equal(t, "Test Club", clubs[0]["name"])
	})

	t.Run("Create Club - Unauthorized", func(t *testing.T) {
		clubData := map[string]string{
			"name":        "New Club",
			"description": "A new test club",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs", clubData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Create Club - Valid", func(t *testing.T) {
		_, token := CreateTestUser(t, "clubcreator@example.com")
		clubData := map[string]string{
			"name":        "New Club",
			"description": "A new test club",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs", clubData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		var club map[string]interface{}
		ParseJSONResponse(t, rr, &club)
		assert.Equal(t, "New Club", club["name"])
		assert.Equal(t, "A new test club", club["description"])
	})

	t.Run("Create Club - Missing Name", func(t *testing.T) {
		_, token := CreateTestUser(t, "clubcreator2@example.com")
		clubData := map[string]string{
			"description": "A club without name",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs", clubData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Name required")
	})

	t.Run("Get Club By ID - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "clubbyid@example.com")
		club := CreateTestClub(t, user, "Club By ID")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID, nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Club By ID - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "clubbyid2@example.com")
		club := CreateTestClub(t, user, "Club By ID Test")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID, nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var returnedClub map[string]interface{}
		ParseJSONResponse(t, rr, &returnedClub)
		assert.Equal(t, "Club By ID Test", returnedClub["name"])
	})

	t.Run("Get Club By ID - Invalid ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "clubbyid3@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/invalid-id", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
		AssertContains(t, rr.Body.String(), "Club not found")
	})

	t.Run("Update Club - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "clubupdate@example.com")
		club := CreateTestClub(t, user, "Club To Update")

		updateData := map[string]string{
			"name": "Updated Club Name",
		}

		req := MakeRequest(t, "PATCH", "/api/v1/clubs/"+club.ID, updateData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		_, token := CreateTestUser(t, "clubmethod@example.com")

		// Test unsupported methods
		methods := []string{"PUT"}
		endpoints := []string{"/api/v1/clubs", "/api/v1/clubs/test-id"}

		for _, method := range methods {
			for _, endpoint := range endpoints {
				req := MakeRequest(t, method, endpoint, nil, token)
				rr := ExecuteRequest(t, handler, req)
				CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
			}
		}
	})

	t.Run("Delete Club - Unauthorized", func(t *testing.T) {
		user, _ := CreateTestUser(t, "clubdelete1@example.com")
		club := CreateTestClub(t, user, "Club To Delete")

		// Try to delete with no token
		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID, nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Delete Club - Forbidden (Not Owner)", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "clubowner@example.com")
		nonOwner, nonOwnerToken := CreateTestUser(t, "nonowner@example.com")
		club := CreateTestClub(t, owner, "Club To Delete")

		// Add non-owner as regular member using test helper
		CreateTestMember(t, nonOwner, club, "member")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID, nil, nonOwnerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
		AssertContains(t, rr.Body.String(), "only owners can delete clubs")
	})

	t.Run("Delete Club - Success", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "clubowner2@example.com")
		club := CreateTestClub(t, owner, "Club To Delete")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID, nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNoContent, rr.Code)

		// Verify club is soft deleted
		var deletedClub models.Club
		err := testDB.First(&deletedClub, "id = ?", club.ID).Error
		assert.NoError(t, err)
		assert.True(t, deletedClub.Deleted)
		assert.NotNil(t, deletedClub.DeletedAt)
		assert.Equal(t, owner.ID, deletedClub.DeletedBy)
	})

	t.Run("Deleted Club Visibility", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "clubowner3@example.com")
		member, memberToken := CreateTestUser(t, "member3@example.com")
		club := CreateTestClub(t, owner, "Club To Delete")
		CreateTestMember(t, member, club, "member")

		// Delete the club
		club.SoftDelete(owner.ID)

		// Owner should still see the club
		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID, nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		// Member should not see the club
		req = MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID, nil, memberToken)
		rr = ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusNotFound, rr.Code)
	})
}