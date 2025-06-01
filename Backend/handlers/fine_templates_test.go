package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFineTemplateEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get Fine Templates - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/clubs/test-club/fine-templates", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Fine Templates - No Templates", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates1@example.com")
		club := CreateTestClub(t, user, "Test Club for Templates")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/fine-templates", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var templates []map[string]interface{}
		ParseJSONResponse(t, rr, &templates)
		assert.Equal(t, 0, len(templates))
	})

	t.Run("Create Fine Template - Unauthorized", func(t *testing.T) {
		payload := map[string]interface{}{
			"description": "Late arrival",
			"amount":      25.0,
		}
		req := MakeRequest(t, "POST", "/api/v1/clubs/test-club/fine-templates", payload, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Create Fine Template - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates2@example.com")
		club := CreateTestClub(t, user, "Test Club for Template Creation")

		payload := map[string]interface{}{
			"description": "Late arrival",
			"amount":      25.0,
		}
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", payload, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusCreated, rr.Code)

		var template map[string]interface{}
		ParseJSONResponse(t, rr, &template)
		assert.Equal(t, "Late arrival", template["description"])
		assert.Equal(t, 25.0, template["amount"])
		assert.Equal(t, club.ID, template["club_id"])
		assert.NotEmpty(t, template["id"])
	})

	t.Run("Create Fine Template - Missing Description", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates3@example.com")
		club := CreateTestClub(t, user, "Test Club for Template Validation")

		payload := map[string]interface{}{
			"amount": 25.0,
		}
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", payload, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Create Fine Template - Invalid Amount", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates4@example.com")
		club := CreateTestClub(t, user, "Test Club for Template Validation")

		payload := map[string]interface{}{
			"description": "Test fine",
			"amount":      -5.0,
		}
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", payload, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Get Fine Templates - With Templates", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates5@example.com")
		club := CreateTestClub(t, user, "Test Club with Templates")

		// Create multiple templates
		template1 := map[string]interface{}{
			"description": "Late arrival",
			"amount":      25.0,
		}
		req1 := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", template1, token)
		rr1 := ExecuteRequest(t, handler, req1)
		CheckResponseCode(t, http.StatusCreated, rr1.Code)

		template2 := map[string]interface{}{
			"description": "No show",
			"amount":      50.0,
		}
		req2 := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", template2, token)
		rr2 := ExecuteRequest(t, handler, req2)
		CheckResponseCode(t, http.StatusCreated, rr2.Code)

		// Get all templates
		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/fine-templates", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var templates []map[string]interface{}
		ParseJSONResponse(t, rr, &templates)
		assert.Equal(t, 2, len(templates))

		// Verify templates
		descriptions := []string{}
		amounts := []float64{}
		for _, template := range templates {
			descriptions = append(descriptions, template["description"].(string))
			amounts = append(amounts, template["amount"].(float64))
		}
		assert.Contains(t, descriptions, "Late arrival")
		assert.Contains(t, descriptions, "No show")
		assert.Contains(t, amounts, 25.0)
		assert.Contains(t, amounts, 50.0)
	})

	t.Run("Update Fine Template - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates6@example.com")
		club := CreateTestClub(t, user, "Test Club for Template Update")

		// Create template
		createPayload := map[string]interface{}{
			"description": "Late arrival",
			"amount":      25.0,
		}
		createReq := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", createPayload, token)
		createRr := ExecuteRequest(t, handler, createReq)
		CheckResponseCode(t, http.StatusCreated, createRr.Code)

		var createdTemplate map[string]interface{}
		ParseJSONResponse(t, createRr, &createdTemplate)
		templateID := createdTemplate["id"].(string)

		// Update template
		updatePayload := map[string]interface{}{
			"description": "Very late arrival",
			"amount":      30.0,
		}
		updateReq := MakeRequest(t, "PUT", "/api/v1/clubs/"+club.ID+"/fine-templates/"+templateID, updatePayload, token)
		updateRr := ExecuteRequest(t, handler, updateReq)
		CheckResponseCode(t, http.StatusOK, updateRr.Code)

		var updatedTemplate map[string]interface{}
		ParseJSONResponse(t, updateRr, &updatedTemplate)
		assert.Equal(t, "Very late arrival", updatedTemplate["description"])
		assert.Equal(t, 30.0, updatedTemplate["amount"])
		assert.Equal(t, templateID, updatedTemplate["id"])
	})

	t.Run("Delete Fine Template - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "templates7@example.com")
		club := CreateTestClub(t, user, "Test Club for Template Delete")

		// Create template
		createPayload := map[string]interface{}{
			"description": "Late arrival",
			"amount":      25.0,
		}
		createReq := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", createPayload, token)
		createRr := ExecuteRequest(t, handler, createReq)
		CheckResponseCode(t, http.StatusCreated, createRr.Code)

		var createdTemplate map[string]interface{}
		ParseJSONResponse(t, createRr, &createdTemplate)
		templateID := createdTemplate["id"].(string)

		// Delete template
		deleteReq := MakeRequest(t, "DELETE", "/api/v1/clubs/"+club.ID+"/fine-templates/"+templateID, nil, token)
		deleteRr := ExecuteRequest(t, handler, deleteReq)
		CheckResponseCode(t, http.StatusNoContent, deleteRr.Code)

		// Verify template is deleted
		getReq := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/fine-templates", nil, token)
		getRr := ExecuteRequest(t, handler, getReq)
		CheckResponseCode(t, http.StatusOK, getRr.Code)

		var templates []map[string]interface{}
		ParseJSONResponse(t, getRr, &templates)
		assert.Equal(t, 0, len(templates))
	})

	t.Run("Non-Admin Cannot Create Template", func(t *testing.T) {
		admin, _ := CreateTestUser(t, "template_admin2@example.com")
		member, memberToken := CreateTestUser(t, "template_member@example.com")
		club := CreateTestClub(t, admin, "Test Club for Admin Check")
		CreateTestMember(t, member, club, "member")

		payload := map[string]interface{}{
			"description": "Late arrival",
			"amount":      25.0,
		}
		req := MakeRequest(t, "POST", "/api/v1/clubs/"+club.ID+"/fine-templates", payload, memberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Non-Member Cannot Access Templates", func(t *testing.T) {
		admin, _ := CreateTestUser(t, "template_admin@example.com")
		_, nonMemberToken := CreateTestUser(t, "template_nonmember@example.com")
		club := CreateTestClub(t, admin, "Test Club for Member Check")

		req := MakeRequest(t, "GET", "/api/v1/clubs/"+club.ID+"/fine-templates", nil, nonMemberToken)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusForbidden, rr.Code)
	})
}