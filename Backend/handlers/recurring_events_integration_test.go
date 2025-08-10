package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/stretchr/testify/assert"
)

// Manual integration test to verify recurring events work end-to-end
func TestRecurringEventsIntegration(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	user, accessToken := CreateTestUser(t, "admin@example.com")
	club := CreateTestClub(t, user, "Test Club")

	t.Run("Integration test - create and verify weekly recurring events", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		recurrenceEnd := startTime.Add(28 * 24 * time.Hour) // 4 weeks

		requestBody := map[string]interface{}{
			"name":                "Weekly Stand-up",
			"description":         "Every Monday team stand-up meeting",
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

		eventCount := int(response["event_count"].(float64))
		assert.True(t, eventCount >= 4, fmt.Sprintf("Should create at least 4 events for 4 weeks, got %d", eventCount))

		// Verify events were created by fetching all events
		getReq := httptest.NewRequest("GET", "/api/v1/clubs/"+club.ID+"/events", nil)
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getReq = getReq.WithContext(context.WithValue(getReq.Context(), auth.UserIDKey, user.ID))

		getW := httptest.NewRecorder()
		mux.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusOK, getW.Code)

		var events []map[string]interface{}
		err = json.Unmarshal(getW.Body.Bytes(), &events)
		assert.NoError(t, err)

		assert.Equal(t, eventCount, len(events), "Number of fetched events should match created count")

		// Verify event properties
		parentEventFound := false
		childEventsFound := 0

		for _, event := range events {
			assert.Equal(t, "Weekly Stand-up", event["name"])
			assert.Equal(t, "Every Monday team stand-up meeting", event["description"])
			assert.Equal(t, "Conference Room A", event["location"])

			if event["is_recurring"].(bool) {
				parentEventFound = true
				assert.Equal(t, "weekly", event["recurrence_pattern"])
				assert.Equal(t, 1.0, event["recurrence_interval"])
				assert.Nil(t, event["parent_event_id"])
			} else {
				childEventsFound++
				assert.NotNil(t, event["parent_event_id"])
				// recurrence_pattern is nil/null for child events
				if pattern := event["recurrence_pattern"]; pattern != nil {
					assert.Equal(t, "", pattern)
				}
			}
		}

		assert.True(t, parentEventFound, "Should have one parent recurring event")
		assert.True(t, childEventsFound >= 3, fmt.Sprintf("Should have at least 3 child events, got %d", childEventsFound))

		t.Logf("✅ Successfully created %d events (%d child events + 1 parent)", len(events), childEventsFound)
	})

	t.Run("Integration test - create daily recurring events", func(t *testing.T) {
		startTime := time.Now().Add(48 * time.Hour) // Start later to avoid overlap
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.Add(5 * 24 * time.Hour) // 5 days

		requestBody := map[string]interface{}{
			"name":                "Daily Check-in",
			"description":         "Daily team check-in",
			"location":            "Online",
			"start_time":          startTime.Format(time.RFC3339),
			"end_time":            endTime.Format(time.RFC3339),
			"recurrence_pattern":  "daily",
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

		eventCount := int(response["event_count"].(float64))
		assert.Equal(t, 6, eventCount, "Should create 6 events for 5 days + parent")

		t.Logf("✅ Successfully created %d daily recurring events", eventCount)
	})
}