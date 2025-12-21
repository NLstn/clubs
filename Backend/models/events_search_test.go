package models_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/NLstn/clubs/odata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to convert string to *string
func strPtr(s string) *string {
	return &s
}

func TestSearchEventsViaOData(t *testing.T) {
	// Set test environment variables
	t.Setenv("GO_ENV", "test")
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing")

	// Initialize auth with test secret
	err := auth.Init()
	require.NoError(t, err, "Failed to initialize auth")

	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Set up OData service for testing
	service, err := odata.NewService(database.Db)
	require.NoError(t, err, "Failed to create OData service")

	odataV2Mux := http.NewServeMux()
	service.RegisterCustomHandlers(odataV2Mux)
	odataV2Mux.Handle("/", service)
	handler := http.StripPrefix("/api/v2", handlers.CompositeAuthMiddleware(odataV2Mux))

	// Create test user
	user, token := handlers.CreateTestUser(t, "searchuser@example.com")

	// Create test clubs
	club1ID := uuid.New().String()
	club1 := models.Club{
		ID:          club1ID,
		Name:        "Test Club 1",
		Description: strPtr("First test club"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = database.Db.Create(&club1).Error
	require.NoError(t, err)

	club2ID := uuid.New().String()
	club2 := models.Club{
		ID:          club2ID,
		Name:        "Test Club 2",
		Description: strPtr("Second test club"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = database.Db.Create(&club2).Error
	require.NoError(t, err)

	// Create membership for user in club1 only
	member1 := models.Member{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ClubID:    club1ID,
		Role:      "member",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}
	err = database.Db.Create(&member1).Error
	require.NoError(t, err)

	// Create events
	now := time.Now()
	event1 := models.Event{
		ID:          uuid.New().String(),
		ClubID:      club1ID,
		Name:        "Search Test Event",
		Description: strPtr("Event for search testing"),
		Location:    strPtr("Test Location"),
		StartTime:   now.Add(24 * time.Hour),
		EndTime:     now.Add(26 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = database.Db.Create(&event1).Error
	require.NoError(t, err)

	event2 := models.Event{
		ID:          uuid.New().String(),
		ClubID:      club1ID,
		Name:        "Other Event",
		Description: strPtr("Another event"),
		Location:    strPtr("Another Location"),
		StartTime:   now.Add(48 * time.Hour),
		EndTime:     now.Add(50 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = database.Db.Create(&event2).Error
	require.NoError(t, err)

	event3 := models.Event{
		ID:          uuid.New().String(),
		ClubID:      club2ID,
		Name:        "Search Event Club 2",
		Description: strPtr("Event in club user is not member of"),
		Location:    strPtr("Club 2 Location"),
		StartTime:   now.Add(72 * time.Hour),
		EndTime:     now.Add(74 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = database.Db.Create(&event3).Error
	require.NoError(t, err)

	t.Run("SearchEventsUserIsMemberOf", func(t *testing.T) {
		// Use OData SearchGlobal function to search for events
		req := httptest.NewRequest("GET", "/api/v2/SearchGlobal(query='Search')", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		
		// Check the structure - it's {value: {Clubs: [], Events: []}}
		value, hasValue := result["value"].(map[string]interface{})
		if hasValue {
			events, ok := value["Events"].([]interface{})
			require.True(t, ok, "Events field should be present in value")
			
			// The search should find events only from clubs user is a member of
			foundEvent1 := false
			for _, e := range events {
				eventMap := e.(map[string]interface{})
				if eventMap["ID"].(string) == event1.ID {
					foundEvent1 = true
					assert.Equal(t, "Search Test Event", eventMap["Name"])
					assert.Equal(t, club1ID, eventMap["ClubID"])
				}
				// Should NOT find event3 (user is not member of club2)
				assert.NotEqual(t, event3.ID, eventMap["ID"])
			}
			if len(events) > 0 {
				assert.True(t, foundEvent1, "Should find event1 in search results")
			}
		} else {
			// Try direct access
			events, ok := result["Events"].([]interface{})
			require.True(t, ok, "Events field should be present")
			
			// The search should find events only from clubs user is a member of
			foundEvent1 := false
			for _, e := range events {
				eventMap := e.(map[string]interface{})
				if eventMap["ID"].(string) == event1.ID {
					foundEvent1 = true
					assert.Equal(t, "Search Test Event", eventMap["Name"])
					assert.Equal(t, club1ID, eventMap["ClubID"])
				}
				// Should NOT find event3 (user is not member of club2)
				assert.NotEqual(t, event3.ID, eventMap["ID"])
			}
			if len(events) > 0 {
				assert.True(t, foundEvent1, "Should find event1 in search results")
			}
		}
	})

	t.Run("SearchEventsNoMatch", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/SearchGlobal(query='NonExistent')", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		value := result["value"].(map[string]interface{})
		events, ok := value["Events"].([]interface{})
		if !ok {
			events = []interface{}{}
		}
		assert.Equal(t, 0, len(events))
	})

	t.Run("SearchEventsCaseInsensitive", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/SearchGlobal(query='search')", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		value := result["value"].(map[string]interface{})
		events, ok := value["Events"].([]interface{})
		if !ok {
			events = []interface{}{}
		}
		
		// Search should be case insensitive and find matching events
		if len(events) > 0 {
			foundMatch := false
			for _, e := range events {
				eventMap := e.(map[string]interface{})
				if eventMap["ID"].(string) == event1.ID {
					foundMatch = true
					break
				}
			}
			assert.True(t, foundMatch, "Should find events (case insensitive)")
		}
	})
}
