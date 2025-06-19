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

	t.Run("RSVP History", func(t *testing.T) {
		// First create an event to test with
		event := CreateTestEvent(t, club, user, "RSVP Test Event")
		assert.NotEmpty(t, event.ID)
		t.Logf("Created event with ID: %s", event.ID)

		// Create initial RSVP
		err := user.CreateOrUpdateRSVP(event.ID, "yes")
		assert.NoError(t, err)

		// Update RSVP to create history
		err = user.CreateOrUpdateRSVP(event.ID, "no")
		assert.NoError(t, err)

		// Update again
		err = user.CreateOrUpdateRSVP(event.ID, "yes")
		assert.NoError(t, err)

		// Test getting RSVP history
		testURL := "/api/v1/clubs/"+club.ID+"/rsvp-history?eventid="+event.ID+"&userid="+user.ID
		t.Logf("Test URL: %s", testURL)
		req := httptest.NewRequest("GET", testURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))
		
		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		// Debug: print the response body for non-200 responses
		if w.Code != 200 {
			t.Logf("Response code: %d", w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		history, ok := response["history"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, history, 3) // Should have 3 entries in history

		// Only check history entries if we have them
		if len(history) >= 3 {
			// Verify the history entries are in chronological order
			firstEntry := history[0].(map[string]interface{})
			assert.Equal(t, "yes", firstEntry["response"])
			
			secondEntry := history[1].(map[string]interface{})
			assert.Equal(t, "no", secondEntry["response"])
			
			thirdEntry := history[2].(map[string]interface{})
			assert.Equal(t, "yes", thirdEntry["response"])
		}
	})

	t.Run("RSVP History - Unauthorized", func(t *testing.T) {
		// Create another user who is not an admin
		otherUser, otherToken := CreateTestUser(t, "other2@example.com")
		
		// Create an event for testing
		event := CreateTestEvent(t, club, user, "RSVP Test Event 2")

		req := httptest.NewRequest("GET", "/api/v1/clubs/"+club.ID+"/rsvp-history?eventid="+event.ID+"&userid="+user.ID, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, otherUser.ID))
		
		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}