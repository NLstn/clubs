package odata

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCreateClubForNewUser tests creating a club for a user who has no existing club memberships
// This is the scenario that was failing in production
func TestCreateClubForNewUser(t *testing.T) {
	// Set test environment
	os.Setenv("GO_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")

	// Initialize auth with test secret
	err := auth.Init()
	require.NoError(t, err, "Failed to initialize auth")

	// Set up file-based SQLite database
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err, "Failed to connect to test database")
	database.Db = testDB

	// Create tables
	testDB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		email TEXT NOT NULL UNIQUE,
		keycloak_id TEXT UNIQUE,
		birth_date DATE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted BOOLEAN DEFAULT FALSE,
		deleted_at DATETIME,
		deleted_by TEXT
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

	// Clean up
	testDB.Exec("DELETE FROM users")
	testDB.Exec("DELETE FROM clubs")
	testDB.Exec("DELETE FROM members")

	// Create OData service
	service, err := NewService(database.Db)
	require.NoError(t, err, "Failed to create OData service")

	// Create a submux to handle both OData and custom routes
	odataV2Mux := http.NewServeMux()
	service.RegisterCustomHandlers(odataV2Mux)
	odataV2Mux.Handle("/", service)

	// Wrap service with auth middleware
	handler := http.StripPrefix("/api/v2", handlers.CompositeAuthMiddleware(odataV2Mux))

	// Create a NEW user who has NO club memberships
	newUser := &models.User{
		ID:        uuid.New().String(),
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}
	require.NoError(t, database.Db.Create(newUser).Error, "Failed to create new user")

	// Generate access token for the new user
	token, err := auth.GenerateAccessToken(newUser.ID)
	require.NoError(t, err, "Failed to generate access token")

	// Try to create a club for this new user (who has no existing memberships)
	newClub := map[string]interface{}{
		"Name":        "First Club",
		"Description": "This is my first club",
	}

	body, err := json.Marshal(newClub)
	require.NoError(t, err, "Failed to marshal club data")

	req := httptest.NewRequest("POST", "/api/v2/Clubs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()

	// Log the response if it's not successful
	if resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		t.Logf("POST failed with status %d: %+v", resp.StatusCode, errResp)
	}

	// This should succeed, but currently fails with:
	// "unauthorized: only admins and owners can add members"
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Creating a club should succeed")

	// Verify the club was created
	if resp.StatusCode == http.StatusCreated {
		var created map[string]interface{}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoError(t, err, "Failed to decode response")

		assert.NotEmpty(t, created["ID"])
		assert.Equal(t, "First Club", created["Name"])
		assert.Equal(t, newUser.ID, created["CreatedBy"])

		// Verify the user is now a member with owner role
		var member models.Member
		err = database.Db.Where("user_id = ? AND club_id = ?", newUser.ID, created["ID"]).First(&member).Error
		assert.NoError(t, err, "Owner member should be created")
		assert.Equal(t, "owner", member.Role, "Creator should be owner")
		assert.NotEmpty(t, member.ID, "Member ID should not be empty")
		t.Logf("Created member: ID=%s, UserID=%s, ClubID=%s, Role=%s", member.ID, member.UserID, member.ClubID, member.Role)
	}
}
