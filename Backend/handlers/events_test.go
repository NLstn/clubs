package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestEventRoutes(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create test user and club
	user, accessToken := CreateTestUser(t, "test@example.com")
	club := CreateTestClub(t, user, "Test Club")

	t.Run("Create Event", func(t *testing.T) {
		payload := map[string]string{
			"name":       "Test Event",
			"start_time": "2024-06-01T10:00:00Z",
			"end_time":   "2024-06-01T12:00:00Z",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/clubs/"+club.ID+"/events", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Event
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Event", response.Name)
		assert.Equal(t, club.ID, response.ClubID)
	})

	t.Run("Create Event as Admin", func(t *testing.T) {
		adminUser, adminToken := CreateTestUser(t, "admin@example.com")
		CreateTestMember(t, adminUser, club, "admin")

		payload := map[string]string{
			"name":       "Admin Event",
			"start_time": "2024-07-01T10:00:00Z",
			"end_time":   "2024-07-01T12:00:00Z",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/clubs/"+club.ID+"/events", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, adminUser.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Event
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, club.ID, response.ClubID)
	})

	t.Run("Get Events", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/clubs/"+club.ID+"/events", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var events []models.Event
		err := json.NewDecoder(w.Body).Decode(&events)
		assert.NoError(t, err)
		assert.Greater(t, len(events), 0)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		// Create another user who is not a member
		otherUser, otherToken := CreateTestUser(t, "other@example.com")

		req := httptest.NewRequest("GET", "/api/v1/clubs/"+club.ID+"/events", nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, otherUser.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
