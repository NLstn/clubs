package handlers

import (
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

func TestSearchRoutes(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Test unauthorized access
	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/search?q=test", nil)

		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		registerSearchRoutes(mux)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test with valid token
	t.Run("Authorized", func(t *testing.T) {
		user, accessToken := CreateTestUser(t, "test@example.com")
		club := CreateTestClub(t, user, "Test Club")

		// Test empty query
		t.Run("EmptyQuery", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/search?q=", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SearchResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(response.Clubs))
			assert.Equal(t, 0, len(response.Events))
		})

		// Test club search
		t.Run("ClubSearch", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/search?q=test", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SearchResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Should find the test club
			assert.GreaterOrEqual(t, len(response.Clubs), 1)

			// Verify club result structure
			if len(response.Clubs) > 0 {
				clubResult := response.Clubs[0]
				assert.Equal(t, "club", clubResult.Type)
				assert.NotEmpty(t, clubResult.ID)
				assert.NotEmpty(t, clubResult.Name)
			}
		})

		// Test case insensitive search
		t.Run("CaseInsensitiveSearch", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/search?q=TEST", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SearchResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Should still find the test club
			assert.GreaterOrEqual(t, len(response.Clubs), 1)
		})

		// Test no results
		t.Run("NoResults", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/search?q=nonexistentxyz", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SearchResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(response.Clubs))
			assert.Equal(t, 0, len(response.Events))
		})

		// Test event search
		t.Run("EventSearch", func(t *testing.T) {
			// Create a test event
			startTime := time.Now().Add(24 * time.Hour)
			endTime := startTime.Add(2 * time.Hour)

			event, err := club.CreateEvent("Test Event", "", "", startTime, endTime, user.ID)
			assert.NoError(t, err)

			req := httptest.NewRequest("GET", "/api/v1/search?q=event", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SearchResponse
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Should find the test event
			assert.GreaterOrEqual(t, len(response.Events), 1)

			// Verify event result structure
			if len(response.Events) > 0 {
				eventResult := response.Events[0]
				assert.Equal(t, "event", eventResult.Type)
				assert.Equal(t, event.ID, eventResult.ID)
				assert.Equal(t, club.ID, eventResult.ClubID)
				assert.NotEmpty(t, eventResult.ClubName)
				assert.NotEmpty(t, eventResult.StartTime)
				assert.NotEmpty(t, eventResult.EndTime)
			}
		})

		// Test that non-members can't see clubs
		t.Run("NonMemberAccess", func(t *testing.T) {
			// Create another user who is not a member of the club
			nonMember, nonMemberToken := CreateTestUser(t, "nonmember@example.com")

			req := httptest.NewRequest("GET", "/api/v1/search?q=test", nil)
			req.Header.Set("Authorization", "Bearer "+nonMemberToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, nonMember.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SearchResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Non-member should not see any results from clubs they're not a member of
			for _, clubResult := range response.Clubs {
				assert.NotEqual(t, club.ID, clubResult.ID, "Non-member should not see clubs they're not a member of")
			}
			for _, eventResult := range response.Events {
				assert.NotEqual(t, club.ID, eventResult.ClubID, "Non-member should not see events from clubs they're not a member of")
			}
		})

		_ = club // Avoid unused variable warning
	})
}

func TestSearchClubs(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)

	user, _ := CreateTestUser(t, "test@example.com")
	club := CreateTestClub(t, user, "Test Club")

	t.Run("ValidSearch", func(t *testing.T) {
		results, err := searchClubs(user, "test")
		assert.NoError(t, err)
		assert.Greater(t, len(results), 0)

		found := false
		for _, result := range results {
			if result.ID == club.ID {
				found = true
				assert.Equal(t, "club", result.Type)
				assert.Equal(t, club.Name, result.Name)
				break
			}
		}
		assert.True(t, found, "Expected to find the created test club in search results")
	})

	t.Run("DescriptionSearch", func(t *testing.T) {
		// Search by description (test club has "Test club description" from CreateTestClub)
		results, err := searchClubs(user, "description")
		assert.NoError(t, err)
		assert.Greater(t, len(results), 0, "Expected to find club by description")
	})

	t.Run("NoMatch", func(t *testing.T) {
		results, err := searchClubs(user, "nonexistentxyz")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(results))
	})
}

func TestSearchEvents(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)

	user, _ := CreateTestUser(t, "test@example.com")
	club := CreateTestClub(t, user, "Test Club")

	// Create a test event
	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	event, err := club.CreateEvent("Search Test Event", "", "", startTime, endTime, user.ID)
	assert.NoError(t, err)

	t.Run("ValidSearch", func(t *testing.T) {
		results, err := searchEvents(user, "Search")
		assert.NoError(t, err)
		assert.Greater(t, len(results), 0)

		found := false
		for _, result := range results {
			if result.ID == event.ID {
				found = true
				assert.Equal(t, "event", result.Type)
				assert.Equal(t, event.Name, result.Name)
				assert.Equal(t, club.ID, result.ClubID)
				break
			}
		}
		assert.True(t, found, "Expected to find the created test event in search results")
	})

	t.Run("NoMatch", func(t *testing.T) {
		results, err := searchEvents(user, "nonexistentxyz")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(results))
	})
}

func TestSearchHTTPMethods(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)

	user, accessToken := CreateTestUser(t, "test@example.com")
	_ = CreateTestClub(t, user, "Test Club")

	// Test unsupported methods
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(fmt.Sprintf("%s_Method", method), func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/search", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			registerSearchRoutes(mux)
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}
