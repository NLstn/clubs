package models

import (
	"context"
	"net/http"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupUserAuthTestDB creates a test database for user authorization tests
func setupUserAuthTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT NOT NULL UNIQUE,
			keycloak_id TEXT UNIQUE,
			birth_date DATE,
			created_at DATETIME,
			updated_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS clubs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			logo_url TEXT,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT,
			deleted BOOLEAN DEFAULT 0,
			deleted_at DATETIME,
			deleted_by TEXT
		);
		CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT DEFAULT 'member',
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
		);
	`).Error
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	database.Db = db
}

// TestUserReadAuthorization verifies that authorization is enforced on User read operations
// Ensures users cannot read other users' private information without proper permissions
func TestUserReadAuthorization(t *testing.T) {
	setupUserAuthTestDB(t)

	// Create three users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	user3ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user3ID, "User", "Three", "user3@test.com")

	// Create two clubs
	club1ID := uuid.New().String()
	club2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		club1ID, "Club 1", user1ID, user1ID)
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		club2ID, "Club 2", user2ID, user2ID)

	// user1 is in club1, user2 is in club2, user3 is in both clubs
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), club1ID, user1ID, "owner", user1ID, user1ID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), club2ID, user2ID, "owner", user2ID, user2ID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), club1ID, user3ID, "member", user1ID, user1ID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), club2ID, user3ID, "member", user2ID, user2ID)

	t.Run("Authorization properly restricts user visibility", func(t *testing.T) {
		user := User{}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		// Try to get authorization scopes for reading users using the OData hook
		// The User model now has ODataBeforeReadCollection; this test verifies it is present
		// and correctly restricts which users can be read by the current user.
		scopes, err := user.ODataBeforeReadCollection(ctx, req, nil)
		if err != nil {
			t.Fatalf("User model is missing or failed ODataBeforeReadCollection - SECURITY ISSUE: %v", err)
		}
		
		// Apply the authorization scopes to restrict the query
			var allUsers []User
			err = database.Db.Find(&allUsers).Error
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}

		// Apply scopes to the database query
		db := database.Db
		for _, scope := range scopes {
			db = scope(db)
		}

		var results []User
		err = db.Find(&results).Error
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}

		// user1 should only see user1 and user3 (users in club1)
		// user1 should NOT see user2 (only in club2)
		if len(results) != 2 {
			t.Errorf("Expected user1 to see exactly 2 users (self + user3 in shared club), got %d", len(results))
		}

		// Verify user2 is NOT in the results
		for _, u := range results {
			if u.ID == user2ID {
				t.Error("CRITICAL: user1 should NOT be able to see user2 (no shared clubs)")
			}
		}
		
		// Verify user1 and user3 ARE in the results
		foundUser1 := false
		foundUser3 := false
		for _, u := range results {
			if u.ID == user1ID {
				foundUser1 = true
			}
			if u.ID == user3ID {
				foundUser3 = true
			}
		}
		if !foundUser1 || !foundUser3 {
			t.Error("user1 should be able to see themselves and user3 (shared club member)")
		}
	})
}

// TestUserUpdateAuthorization tests that authorization is enforced on User update operations
func TestUserUpdateAuthorization(t *testing.T) {
	setupUserAuthTestDB(t)

	// Create two users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")

	t.Run("User cannot update another user's information", func(t *testing.T) {
		user := User{
			ID:        user2ID,
			FirstName: "Modified",
			LastName:  "Name",
			Email:     "user2@test.com",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		// Try to update user2 while logged in as user1
		// The User model's ODataBeforeUpdate hook should prevent this unauthorized update
		err := user.ODataBeforeUpdate(ctx, req)
		if err == nil {
			// If the hook allows this, it's a critical vulnerability
			t.Error("CRITICAL: user1 was able to update user2's information!")
		} else {
			// Expected: authorization hook rejects updates to other users' data
			t.Logf("ODataBeforeUpdate correctly prevented unauthorized update: %v", err)
		}
	})
}

// TestUserCreateAuthorization tests that user creation via OData API is prevented
func TestUserCreateAuthorization(t *testing.T) {
	setupUserAuthTestDB(t)

	// Create an authenticated user
	user1ID := uuid.New().String()
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")

	t.Run("Direct user creation via OData API is prevented", func(t *testing.T) {
		newUser := User{
			FirstName: "New",
			LastName:  "User",
			Email:     "newuser@test.com",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		// Try to create a new user directly via OData API
		// This should be prevented - user creation must go through auth endpoints
		err := newUser.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("CRITICAL: Direct user creation via OData API was allowed!")
		} else {
			t.Logf("ODataBeforeCreate correctly prevented direct user creation: %v", err)
		}
	})
}
