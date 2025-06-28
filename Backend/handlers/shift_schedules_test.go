package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get Shifts - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner1@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/shifts", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Shifts - Valid", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner2@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/shifts", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var shifts []map[string]interface{}
		ParseJSONResponse(t, rr, &shifts)
		// Should return empty array initially
		assert.Equal(t, 0, len(shifts))
	})

	t.Run("Get Shifts - Invalid Club ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "user3@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/invalid-uuid/shifts", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Create Shift - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner4@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		shiftData := map[string]interface{}{
			"startTime": "2024-01-01T09:00:00Z",
			"endTime":   "2024-01-01T17:00:00Z",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts", shiftData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Create Shift - No Longer Supported", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner5@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		shiftData := map[string]interface{}{
			"startTime": "2024-01-01T09:00:00Z",
			"endTime":   "2024-01-01T17:00:00Z",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts", shiftData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Create Shift - No Longer Supported (Missing Start Time)", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner6@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		shiftData := map[string]interface{}{
			"endTime": "2024-01-01T17:00:00Z",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts", shiftData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
		AssertContains(t, rr.Body.String(), "Method not allowed")
	})

	t.Run("Create Shift - No Longer Supported (Missing End Time)", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner7@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		shiftData := map[string]interface{}{
			"startTime": "2024-01-01T09:00:00Z",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts", shiftData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
		AssertContains(t, rr.Body.String(), "Method not allowed")
	})

	t.Run("Create Shift - No Longer Supported (Invalid Club ID)", func(t *testing.T) {
		_, token := CreateTestUser(t, "user8@example.com")

		shiftData := map[string]interface{}{
			"startTime": "2024-01-01T09:00:00Z",
			"endTime":   "2024-01-01T17:00:00Z",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/invalid-uuid/shifts", shiftData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
		AssertContains(t, rr.Body.String(), "Method not allowed")
	})

	t.Run("Get Shift Members - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner9@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/shifts/test-shift-id/members", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Shift Members - Invalid Club ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "user10@example.com")

		req := MakeRequest(t, "GET", "/api/v1/clubs/invalid-uuid/shifts/test-shift-id/members", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Get Shift Members - Invalid Shift ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner11@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/shifts/invalid-uuid/members", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid shift ID format")
	})

	t.Run("Add Member to Shift - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner12@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		memberData := map[string]string{
			"userId": "550e8400-e29b-41d4-a716-446655440000",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts/test-shift-id/members", memberData, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Add Member to Shift - Invalid Club ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "user13@example.com")

		memberData := map[string]string{
			"userId": "550e8400-e29b-41d4-a716-446655440000",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/invalid-uuid/shifts/test-shift-id/members", memberData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Add Member to Shift - Invalid Shift ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner14@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		memberData := map[string]string{
			"userId": "550e8400-e29b-41d4-a716-446655440000",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts/invalid-uuid/members", memberData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid shift ID format")
	})

	t.Run("Add Member to Shift - Invalid User ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner15@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		memberData := map[string]string{
			"userId": "invalid-uuid",
		}

		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/shifts/550e8400-e29b-41d4-a716-446655440000/members", memberData, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid user ID format")
	})

	t.Run("Remove Member from Shift - Unauthorized", func(t *testing.T) {
		owner, _ := CreateTestUser(t, "owner16@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/shifts/test-shift-id/members/test-member-id", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Remove Member from Shift - Invalid Club ID", func(t *testing.T) {
		_, token := CreateTestUser(t, "user17@example.com")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/invalid-uuid/shifts/test-shift-id/members/test-member-id", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid club ID format")
	})

	t.Run("Remove Member from Shift - Invalid Shift ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner18@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/shifts/invalid-uuid/members/test-member-id", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid shift ID format")
	})

	t.Run("Remove Member from Shift - Invalid Member ID", func(t *testing.T) {
		owner, ownerToken := CreateTestUser(t, "owner19@example.com")
		club := CreateTestClub(t, owner, "Test Club")

		req := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/shifts/550e8400-e29b-41d4-a716-446655440000/members/invalid-uuid", nil, ownerToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		AssertContains(t, rr.Body.String(), "Invalid member ID format")
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		_, token := CreateTestUser(t, "test20@example.com")

		endpoints := []string{
			"/api/v1/clubs/test-id/shifts",
			"/api/v1/clubs/test-id/shifts/shift-id/members",
			"/api/v1/clubs/test-id/shifts/shift-id/members/member-id",
		}

		invalidMethods := []string{"PUT", "PATCH"}

		for _, endpoint := range endpoints {
			for _, method := range invalidMethods {
				req := MakeRequest(t, method, endpoint, nil, token)
				rr := ExecuteRequest(t, handler, req)
				CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
			}
		}
	})
}
