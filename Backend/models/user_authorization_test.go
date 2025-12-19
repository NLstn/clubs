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

// TestUserReadAuthorizationMissing tests for missing authorization on User read operations
// This is a CRITICAL security vulnerability where users can read other users' private information
func TestUserReadAuthorizationMissing(t *testing.T) {
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

	t.Run("VULNERABILITY: User can read all users without authorization", func(t *testing.T) {
		user := User{}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		// Try to get authorization scopes for reading users
		// The User model should have ODataBeforeReadCollection but currently doesn't
		// This will fail compilation which is expected - we need to add the hooks
		scopes, err := user.ODataBeforeReadCollection(ctx, req, nil)
		if err != nil {
			t.Logf("User model does not have ODataBeforeReadCollection - SECURITY ISSUE CONFIRMED")
			t.Logf("Any authenticated user can query all users in the system!")
			
			// Currently, users can query ALL users without restriction
			var allUsers []User
			err = database.Db.Find(&allUsers).Error
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}

			if len(allUsers) == 3 {
				t.Error("CRITICAL VULNERABILITY: user1 can see ALL users (including user2 from different club)")
				t.Error("Expected: Users should only see users from clubs they share")
				t.Error("Actual: All 3 users are visible")
			}
		} else {
			// If hooks exist, verify they work correctly
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
			if len(results) > 2 {
				t.Error("User can see users outside of their shared clubs")
			}

			// Verify user2 is NOT in the results
			for _, u := range results {
				if u.ID == user2ID {
					t.Error("user1 should NOT be able to see user2 (no shared clubs)")
				}
			}
		}
	})
}

// TestUserUpdateAuthorizationMissing tests for missing authorization on User update operations
func TestUserUpdateAuthorizationMissing(t *testing.T) {
	setupUserAuthTestDB(t)

	// Create two users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")

	t.Run("VULNERABILITY: User might be able to update other users' information", func(t *testing.T) {
		user := User{
			ID:        user2ID,
			FirstName: "Modified",
			LastName:  "Name",
			Email:     "user2@test.com",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		// Try to update user2 while logged in as user1
		// The User model should have ODataBeforeUpdate but currently doesn't
		err := user.ODataBeforeUpdate(ctx, req)
		if err != nil {
			// This is expected - there's no such method yet
			t.Logf("User model does not have ODataBeforeUpdate - POTENTIAL SECURITY ISSUE")
			t.Logf("Need to verify if OData framework prevents this automatically")
		} else {
			// If the hook exists and returns no error, that would be a vulnerability
			t.Error("CRITICAL: user1 was able to update user2's information!")
		}
	})
}
