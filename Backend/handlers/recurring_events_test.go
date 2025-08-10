package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/stretchr/testify/assert"
)

func TestRecurringEventHandler(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	user, accessToken := CreateTestUser(t, "eventcreator@example.com")
	club := CreateTestClub(t, user, "Test Club")

	t.Run("create weekly recurring event", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		recurrenceEnd := startTime.Add(30 * 24 * time.Hour)

		requestBody := map[string]interface{}{
			"name":                "Weekly Team Meeting",
			"description":         "Every Monday team meeting",
			"location":            "Conference Room A",
			"start_time":          startTime.Format(time.RFC3339),
			"end_time":            endTime.Format(time.RFC3339),
			"recurrence_pattern":  "weekly",
			"recurrence_interval": 1,
			"recurrence_end":      recurrenceEnd.Format(time.RFC3339),
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/api/v1/clubs/"+club.ID+"/events/recurring", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "event_count")
		assert.Contains(t, response, "events")

		eventCount := response["event_count"].(float64)
		assert.True(t, eventCount > 1, "Should create multiple events")
	})

	t.Run("create recurring event - unauthorized", func(t *testing.T) {
		// Create a regular member (not admin)
		member, memberToken := CreateTestUser(t, "member@example.com")
		CreateTestMember(t, member, club, "member")

		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		recurrenceEnd := startTime.Add(7 * 24 * time.Hour)

		requestBody := map[string]interface{}{
			"name":                "Weekly Meeting",
			"description":         "Test",
			"location":            "Room",
			"start_time":          startTime.Format(time.RFC3339),
			"end_time":            endTime.Format(time.RFC3339),
			"recurrence_pattern":  "weekly",
			"recurrence_interval": 1,
			"recurrence_end":      recurrenceEnd.Format(time.RFC3339),
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/api/v1/clubs/"+club.ID+"/events/recurring", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+memberToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, member.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("create recurring event - invalid pattern", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		recurrenceEnd := startTime.Add(7 * 24 * time.Hour)

		requestBody := map[string]interface{}{
			"name":                "Invalid Pattern Event",
			"description":         "Test invalid pattern",
			"location":            "Room",
			"start_time":          startTime.Format(time.RFC3339),
			"end_time":            endTime.Format(time.RFC3339),
			"recurrence_pattern":  "invalid",
			"recurrence_interval": 1,
			"recurrence_end":      recurrenceEnd.Format(time.RFC3339),
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/api/v1/clubs/"+club.ID+"/events/recurring", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid recurrence pattern")
	})

	t.Run("create recurring event - method not allowed", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/clubs/"+club.ID+"/events/recurring", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerEventRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}