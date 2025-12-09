package odata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testContext holds common test fixtures
type testContext struct {
	service    *Service
	handler    http.Handler
	testUser   *models.User
	testUser2  *models.User
	testClub   *models.Club
	testMember *models.Member
	token      string
}

// setupTestContext creates a fresh test environment for each test
func setupTestContext(t *testing.T) *testContext {
	// Set test environment
	os.Setenv("GO_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")

	// Initialize auth with test secret
	err := auth.Init()
	require.NoError(t, err, "Failed to initialize auth")

	// Set up in-memory SQLite database with custom configuration
	// SQLite doesn't support all PostgreSQL features, so we need to be careful with auto-migration
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Disable foreign key constraints in SQLite for simpler testing
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err, "Failed to connect to test database")
	database.Db = testDB

	// Create tables manually with SQLite-compatible SQL to avoid UUID function issues
	// This mirrors the approach used in handlers/test_helpers.go
	testDB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		email TEXT NOT NULL UNIQUE,
		keycloak_id TEXT UNIQUE,
		birth_date DATE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS clubs (
		id TEXT PRIMARY KEY,
		name TEXT,
		description TEXT,
		logo_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT,
		deleted BOOLEAN DEFAULT FALSE,
		deleted_at DATETIME,
		deleted_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS members (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT NOT NULL,
		role TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS teams (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		name TEXT,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT,
		deleted BOOLEAN DEFAULT FALSE,
		deleted_at DATETIME,
		deleted_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		team_id TEXT,
		name TEXT,
		description TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		location TEXT,
		recurrence_pattern TEXT,
		recurrence_end_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS event_rsvps (
		id TEXT PRIMARY KEY,
		event_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		status TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS shifts (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		event_id TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS shift_members (
		id TEXT PRIMARY KEY,
		shift_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS fines (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT NOT NULL,
		team_id TEXT,
		reason TEXT,
		amount REAL,
		paid BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS fine_templates (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		description TEXT,
		amount REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS invites (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS join_requests (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS news (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		title TEXT,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT,
		type TEXT,
		title TEXT,
		message TEXT,
		is_read BOOLEAN DEFAULT FALSE,
		related_entity_id TEXT,
		related_entity_type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS user_notification_preferences (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT,
		email_enabled BOOLEAN DEFAULT TRUE,
		push_enabled BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS club_settings (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL UNIQUE,
		enable_fines BOOLEAN DEFAULT TRUE,
		enable_shifts BOOLEAN DEFAULT TRUE,
		enable_events BOOLEAN DEFAULT TRUE,
		enable_teams BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS user_privacy_settings (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT,
		hide_email BOOLEAN DEFAULT FALSE,
		hide_phone BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	// Create OData service
	service, err := NewService(database.Db)
	require.NoError(t, err, "Failed to create OData service")

	// Wrap service with auth middleware
	handler := http.StripPrefix("/api/v2", AuthMiddleware(auth.GetJWTSecret())(service))

	// Create test users
	testUser := &models.User{
		ID:        uuid.New().String(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	require.NoError(t, database.Db.Create(testUser).Error, "Failed to create test user")

	testUser2 := &models.User{
		ID:        uuid.New().String(),
		Email:     "test2@example.com",
		FirstName: "Test",
		LastName:  "User 2",
	}
	require.NoError(t, database.Db.Create(testUser2).Error, "Failed to create test user 2")

	// Generate access token for testUser
	token, err := auth.GenerateAccessToken(testUser.ID)
	require.NoError(t, err, "Failed to generate access token")

	// Create test club with member
	description := "A test club"
	testClub := &models.Club{
		ID:          uuid.New().String(),
		Name:        "Test Club",
		Description: &description,
		CreatedBy:   testUser.ID,
		UpdatedBy:   testUser.ID,
		Deleted:     false,
	}
	require.NoError(t, database.Db.Create(testClub).Error, "Failed to create test club")

	testMember := &models.Member{
		ID:        uuid.New().String(),
		ClubID:    testClub.ID,
		UserID:    testUser.ID,
		Role:      "owner",
		CreatedBy: testUser.ID,
		UpdatedBy: testUser.ID,
	}
	require.NoError(t, database.Db.Create(testMember).Error, "Failed to create test member")

	return &testContext{
		service:    service,
		handler:    handler,
		testUser:   testUser,
		testUser2:  testUser2,
		testClub:   testClub,
		testMember: testMember,
		token:      token,
	}
}

// makeAuthenticatedRequest creates an HTTP request with JWT token
func (ctx *testContext) makeAuthenticatedRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	// Marshal body if provided
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		require.NoError(t, err, "Failed to marshal request body")
	}

	// Create request
	req := httptest.NewRequest(method, "/api/v2"+path, bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+ctx.token)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	return w.Result()
}

// parseJSONResponse parses the response body as JSON
func parseJSONResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "Failed to decode response body")
}

// TestClubCRUD tests basic CRUD operations for Club entity
func TestClubCRUD(t *testing.T) {
	ctx := setupTestContext(t)

	t.Run("GET collection - all clubs", func(t *testing.T) {
		resp := ctx.makeAuthenticatedRequest(t, "GET", "/Clubs", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		// Verify OData response structure
		assert.Contains(t, result, "value")
		values := result["value"].([]interface{})
		assert.GreaterOrEqual(t, len(values), 1, "Should have at least one club")

		// Verify club data
		club := values[0].(map[string]interface{})
		assert.Equal(t, ctx.testClub.ID, club["id"])
		assert.Equal(t, ctx.testClub.Name, club["name"])
	})

	t.Run("GET single club by key", func(t *testing.T) {
		// OData key format without quotes for UUID strings
		path := fmt.Sprintf("/Clubs(%s)", ctx.testClub.ID)
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var club map[string]interface{}
		parseJSONResponse(t, resp, &club)

		assert.Equal(t, ctx.testClub.ID, club["id"])
		assert.Equal(t, ctx.testClub.Name, club["name"])
		// Description is returned as string, not pointer
		if ctx.testClub.Description != nil {
			assert.Equal(t, *ctx.testClub.Description, club["description"])
		}
	})

	t.Run("POST create new club", func(t *testing.T) {
		newClub := map[string]interface{}{
			"name":        "New Test Club",
			"description": "A newly created club",
			// These fields will be set automatically by hooks in Phase 4
			// For now, we pass them explicitly for testing
			"createdBy": ctx.testUser.ID,
			"updatedBy": ctx.testUser.ID,
		}

		resp := ctx.makeAuthenticatedRequest(t, "POST", "/Clubs", newClub)

		// If status is not 201, print the error response
		if resp.StatusCode != http.StatusCreated {
			var errResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errResp)
			t.Logf("POST failed with status %d: %+v", resp.StatusCode, errResp)
		}

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created map[string]interface{}
		parseJSONResponse(t, resp, &created)

		assert.NotEmpty(t, created["id"])
		assert.Equal(t, "New Test Club", created["name"])
		assert.Equal(t, "A newly created club", created["description"])
		assert.Equal(t, ctx.testUser.ID, created["createdBy"])
		assert.Equal(t, ctx.testUser.ID, created["updatedBy"])
		assert.Equal(t, false, created["deleted"])
	})

	t.Run("PATCH update existing club", func(t *testing.T) {
		update := map[string]interface{}{
			"name":        "Updated Club Name",
			"description": "Updated description",
		}

		path := fmt.Sprintf("/Clubs(%s)", ctx.testClub.ID)
		resp := ctx.makeAuthenticatedRequest(t, "PATCH", path, update)

		// PATCH may return 204 No Content on success
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}

		// If 204, verify by fetching the entity
		if resp.StatusCode == http.StatusNoContent {
			getResp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
			assert.Equal(t, http.StatusOK, getResp.StatusCode)

			var updated map[string]interface{}
			parseJSONResponse(t, getResp, &updated)

			assert.Equal(t, "Updated Club Name", updated["name"])
			assert.Equal(t, "Updated description", updated["description"])
		} else {
			var updated map[string]interface{}
			parseJSONResponse(t, resp, &updated)

			assert.Equal(t, "Updated Club Name", updated["name"])
			assert.Equal(t, "Updated description", updated["description"])
			assert.Equal(t, ctx.testUser.ID, updated["updatedBy"])
		}
	})

	t.Run("DELETE soft delete club", func(t *testing.T) {
		// Create a new club specifically for deletion test
		desc := "Club for deletion test"
		clubToDelete := &models.Club{
			ID:          uuid.New().String(),
			Name:        "Club To Delete",
			Description: &desc,
			CreatedBy:   ctx.testUser.ID,
			UpdatedBy:   ctx.testUser.ID,
			Deleted:     false,
		}
		require.NoError(t, database.Db.Create(clubToDelete).Error)

		// Create membership
		memberForDelete := &models.Member{
			ID:        uuid.New().String(),
			ClubID:    clubToDelete.ID,
			UserID:    ctx.testUser.ID,
			Role:      "owner",
			CreatedBy: ctx.testUser.ID,
			UpdatedBy: ctx.testUser.ID,
		}
		require.NoError(t, database.Db.Create(memberForDelete).Error)

		path := fmt.Sprintf("/Clubs(%s)", clubToDelete.ID)
		resp := ctx.makeAuthenticatedRequest(t, "DELETE", path, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify club is soft deleted in database
		var club models.Club
		err := database.Db.Where("id = ?", clubToDelete.ID).First(&club).Error
		require.NoError(t, err)
		assert.True(t, club.Deleted)
		assert.NotNil(t, club.DeletedAt)
		assert.NotNil(t, club.DeletedBy)
		if club.DeletedBy != nil {
			assert.Equal(t, ctx.testUser.ID, *club.DeletedBy)
		}
	})
}

// TestMemberCRUD tests CRUD operations for Member entity
func TestMemberCRUD(t *testing.T) {
	ctx := setupTestContext(t)

	t.Run("GET members filtered by club", func(t *testing.T) {
		// OData string comparison requires quotes around UUID values
		path := fmt.Sprintf("/Members?$filter=clubId%%20eq%%20'%s'", ctx.testClub.ID)
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)

		if resp.StatusCode != http.StatusOK {
			var errResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errResp)
			t.Logf("GET Members filter failed with status %d: %+v", resp.StatusCode, errResp)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		assert.Equal(t, 1, len(values))

		member := values[0].(map[string]interface{})
		assert.Equal(t, ctx.testMember.ID, member["id"])
		assert.Equal(t, ctx.testClub.ID, member["clubId"])
		assert.Equal(t, ctx.testUser.ID, member["userId"])
	})

	t.Run("GET members with expanded user", func(t *testing.T) {
		path := fmt.Sprintf("/Members?$filter=clubId%%20eq%%20'%s'&$expand=User", ctx.testClub.ID)
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		member := values[0].(map[string]interface{})

		// Verify expanded user data
		assert.Contains(t, member, "User")
		user := member["User"].(map[string]interface{})
		assert.Equal(t, ctx.testUser.ID, user["id"])
		assert.Equal(t, ctx.testUser.Email, user["email"])
	})

	t.Run("POST create new member", func(t *testing.T) {
		newMember := map[string]interface{}{
			"clubId":    ctx.testClub.ID,
			"userId":    ctx.testUser2.ID,
			"role":      "member",
			"createdBy": ctx.testUser.ID,
			"updatedBy": ctx.testUser.ID,
		}

		resp := ctx.makeAuthenticatedRequest(t, "POST", "/Members", newMember)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created map[string]interface{}
		parseJSONResponse(t, resp, &created)

		assert.NotEmpty(t, created["id"])
		assert.Equal(t, ctx.testClub.ID, created["clubId"])
		assert.Equal(t, ctx.testUser2.ID, created["userId"])
		assert.Equal(t, "member", created["role"])
	})

	t.Run("PATCH update member role", func(t *testing.T) {
		update := map[string]interface{}{
			"role": "admin",
		}

		path := fmt.Sprintf("/Members(%s)", ctx.testMember.ID)
		resp := ctx.makeAuthenticatedRequest(t, "PATCH", path, update)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated map[string]interface{}
		parseJSONResponse(t, resp, &updated)

		assert.Equal(t, "admin", updated["role"])
	})
}

// TestEventCRUD tests CRUD operations for Event entity
func TestEventCRUD(t *testing.T) {
	ctx := setupTestContext(t)

	// Create test event
	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	t.Run("POST create event", func(t *testing.T) {
		newEvent := map[string]interface{}{
			"clubId":      ctx.testClub.ID,
			"name":        "Test Event",
			"description": "A test event",
			"startTime":   startTime.Format(time.RFC3339),
			"endTime":     endTime.Format(time.RFC3339),
			"location":    "Test Location",
			"createdBy":   ctx.testUser.ID,
			"updatedBy":   ctx.testUser.ID,
		}

		resp := ctx.makeAuthenticatedRequest(t, "POST", "/Events", newEvent)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created map[string]interface{}
		parseJSONResponse(t, resp, &created)

		assert.NotEmpty(t, created["id"])
		assert.Equal(t, "Test Event", created["name"])
		assert.Equal(t, ctx.testClub.ID, created["clubId"])
	})

	t.Run("GET events filtered by club and ordered by startTime", func(t *testing.T) {
		// Create another event
		desc2 := "Another event"
		loc2 := "Another Location"
		newEvent := &models.Event{
			ID:          uuid.New().String(),
			ClubID:      ctx.testClub.ID,
			Name:        "Second Event",
			Description: &desc2,
			StartTime:   time.Now().Add(48 * time.Hour),
			EndTime:     time.Now().Add(50 * time.Hour),
			Location:    &loc2,
			CreatedBy:   ctx.testUser.ID,
			UpdatedBy:   ctx.testUser.ID,
		}
		require.NoError(t, database.Db.Create(newEvent).Error)

		path := fmt.Sprintf("/Events?$filter=clubId%%20eq%%20'%s'&$orderby=startTime%%20asc", ctx.testClub.ID)
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		assert.GreaterOrEqual(t, len(values), 1)

		// Verify ordering
		if len(values) >= 2 {
			event1 := values[0].(map[string]interface{})
			event2 := values[1].(map[string]interface{})

			time1, err1 := time.Parse(time.RFC3339, event1["startTime"].(string))
			time2, err2 := time.Parse(time.RFC3339, event2["startTime"].(string))
			require.NoError(t, err1)
			require.NoError(t, err2)

			assert.True(t, time1.Before(time2), "Events should be ordered by startTime ascending")
		}
	})

	t.Run("PATCH update event", func(t *testing.T) {
		// First create an event
		desc3 := "Original description"
		loc3 := "Original Location"
		event := &models.Event{
			ID:          uuid.New().String(),
			ClubID:      ctx.testClub.ID,
			Name:        "Event to Update",
			Description: &desc3,
			StartTime:   time.Now().Add(72 * time.Hour),
			EndTime:     time.Now().Add(74 * time.Hour),
			Location:    &loc3,
			CreatedBy:   ctx.testUser.ID,
			UpdatedBy:   ctx.testUser.ID,
		}
		require.NoError(t, database.Db.Create(event).Error)

		update := map[string]interface{}{
			"name":        "Updated Event Name",
			"description": "Updated description",
			"location":    "Updated Location",
		}

		path := fmt.Sprintf("/Events(%s)", event.ID)
		resp := ctx.makeAuthenticatedRequest(t, "PATCH", path, update)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated map[string]interface{}
		parseJSONResponse(t, resp, &updated)

		assert.Equal(t, "Updated Event Name", updated["name"])
		assert.Equal(t, "Updated description", updated["description"])
		assert.Equal(t, "Updated Location", updated["location"])
	})
}

// TestODataQueryFeatures tests advanced OData query capabilities
func TestODataQueryFeatures(t *testing.T) {
	ctx := setupTestContext(t)

	// Create multiple clubs for testing
	desc4 := "A soccer club"
	club2 := &models.Club{
		ID:          uuid.New().String(),
		Name:        "Soccer Club",
		Description: &desc4,
		CreatedBy:   ctx.testUser.ID,
		UpdatedBy:   ctx.testUser.ID,
		Deleted:     false,
	}
	require.NoError(t, database.Db.Create(club2).Error)

	desc5 := "A chess club"
	club3 := &models.Club{
		ID:          uuid.New().String(),
		Name:        "Chess Club",
		Description: &desc5,
		CreatedBy:   ctx.testUser.ID,
		UpdatedBy:   ctx.testUser.ID,
		Deleted:     false,
	}
	require.NoError(t, database.Db.Create(club3).Error)

	// Create members for all clubs
	member2 := &models.Member{
		ID:        uuid.New().String(),
		ClubID:    club2.ID,
		UserID:    ctx.testUser.ID,
		Role:      "member",
		CreatedBy: ctx.testUser.ID,
		UpdatedBy: ctx.testUser.ID,
	}
	require.NoError(t, database.Db.Create(member2).Error)

	member3 := &models.Member{
		ID:        uuid.New().String(),
		ClubID:    club3.ID,
		UserID:    ctx.testUser.ID,
		Role:      "member",
		CreatedBy: ctx.testUser.ID,
		UpdatedBy: ctx.testUser.ID,
	}
	require.NoError(t, database.Db.Create(member3).Error)

	t.Run("$select - return only specified fields", func(t *testing.T) {
		path := "/Clubs?$select=id,name"
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		assert.GreaterOrEqual(t, len(values), 1)

		club := values[0].(map[string]interface{})
		assert.Contains(t, club, "id")
		assert.Contains(t, club, "name")
		// Description should not be included
		assert.NotContains(t, club, "description")
	})

	t.Run("$filter - filter by name", func(t *testing.T) {
		path := "/Clubs?$filter=name%%20eq%%20'Soccer%%20Club'"
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		assert.Equal(t, 1, len(values))

		club := values[0].(map[string]interface{})
		assert.Equal(t, "Soccer Club", club["name"])
	})

	t.Run("$orderby - sort clubs by name", func(t *testing.T) {
		path := "/Clubs?$orderby=name%%20asc"
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		assert.GreaterOrEqual(t, len(values), 3)

		// Verify alphabetical order
		for i := 0; i < len(values)-1; i++ {
			name1 := values[i].(map[string]interface{})["name"].(string)
			name2 := values[i+1].(map[string]interface{})["name"].(string)
			assert.LessOrEqual(t, name1, name2)
		}
	})

	t.Run("$top and $skip - pagination", func(t *testing.T) {
		path := "/Clubs?$orderby=name%%20asc&$top=2&$skip=1"
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		parseJSONResponse(t, resp, &result)

		values := result["value"].([]interface{})
		assert.LessOrEqual(t, len(values), 2)
	})

	t.Run("$count - return count of records", func(t *testing.T) {
		path := "/Clubs/$count"
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var count int
		err := json.NewDecoder(resp.Body).Decode(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 3)
	})

	t.Run("$expand - include related members", func(t *testing.T) {
		path := fmt.Sprintf("/Clubs(%s)?$expand=Members", ctx.testClub.ID)
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var club map[string]interface{}
		parseJSONResponse(t, resp, &club)

		assert.Contains(t, club, "Members")
		members := club["Members"].([]interface{})
		assert.GreaterOrEqual(t, len(members), 1)

		member := members[0].(map[string]interface{})
		assert.Equal(t, ctx.testClub.ID, member["clubId"])
	})
}

// TestSoftDeleteFiltering tests that deleted entities are filtered out
func TestSoftDeleteFiltering(t *testing.T) {
	ctx := setupTestContext(t)

	// Create a club that we'll soft delete
	desc6 := "This will be deleted"
	clubToDelete := &models.Club{
		ID:          uuid.New().String(),
		Name:        "Club to Delete",
		Description: &desc6,
		CreatedBy:   ctx.testUser.ID,
		UpdatedBy:   ctx.testUser.ID,
		Deleted:     false,
	}
	require.NoError(t, database.Db.Create(clubToDelete).Error)

	// Create member so user can access the club
	memberToDelete := &models.Member{
		ID:        uuid.New().String(),
		ClubID:    clubToDelete.ID,
		UserID:    ctx.testUser.ID,
		Role:      "owner",
		CreatedBy: ctx.testUser.ID,
		UpdatedBy: ctx.testUser.ID,
	}
	require.NoError(t, database.Db.Create(memberToDelete).Error)

	t.Run("deleted club is filtered from collection", func(t *testing.T) {
		// Get initial count
		resp := ctx.makeAuthenticatedRequest(t, "GET", "/Clubs", nil)
		var result1 map[string]interface{}
		parseJSONResponse(t, resp, &result1)
		initialCount := len(result1["value"].([]interface{}))

		// Soft delete the club
		path := fmt.Sprintf("/Clubs(%s)", clubToDelete.ID)
		deleteResp := ctx.makeAuthenticatedRequest(t, "DELETE", path, nil)
		assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

		// Get clubs again - deleted club should not appear
		resp2 := ctx.makeAuthenticatedRequest(t, "GET", "/Clubs", nil)
		var result2 map[string]interface{}
		parseJSONResponse(t, resp2, &result2)
		afterDeleteCount := len(result2["value"].([]interface{}))

		assert.Equal(t, initialCount-1, afterDeleteCount, "Deleted club should not appear in collection")

		// Verify we can't get the deleted club by ID
		getSingle := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusNotFound, getSingle.StatusCode, "Deleted club should return 404")
	})

	t.Run("club is actually soft deleted in database", func(t *testing.T) {
		var club models.Club
		err := database.Db.Unscoped().Where("id = ?", clubToDelete.ID).First(&club).Error
		require.NoError(t, err)

		assert.True(t, club.Deleted, "Club should be marked as deleted")
		assert.NotNil(t, club.DeletedAt, "DeletedAt should be set")
		assert.NotNil(t, club.DeletedBy, "DeletedBy should be set")
		assert.Equal(t, ctx.testUser.ID, *club.DeletedBy, "DeletedBy should be the authenticated user")
	})
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	ctx := setupTestContext(t)

	t.Run("unauthorized request without token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/Clubs", nil)
		w := httptest.NewRecorder()
		ctx.handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("not found - nonexistent entity", func(t *testing.T) {
		fakeID := uuid.New().String()
		path := fmt.Sprintf("/Clubs(%s)", fakeID)
		resp := ctx.makeAuthenticatedRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("bad request - invalid entity data", func(t *testing.T) {
		invalidClub := map[string]interface{}{
			"name": "", // Empty name should fail validation
		}

		resp := ctx.makeAuthenticatedRequest(t, "POST", "/Clubs", invalidClub)
		// The actual status code depends on validation implementation
		// Could be 400 Bad Request or 422 Unprocessable Entity
		assert.NotEqual(t, http.StatusCreated, resp.StatusCode)
	})
}
