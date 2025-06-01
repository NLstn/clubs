package models

import (
	"testing"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForJoinRequest(t *testing.T) {
	// Create in-memory SQLite database for testing
	var err error
	database.Db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create SQLite-compatible tables manually
	database.Db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT,
			email TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	database.Db.Exec(`
		CREATE TABLE IF NOT EXISTS clubs (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	database.Db.Exec(`
		CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			role TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	database.Db.Exec(`
		CREATE TABLE IF NOT EXISTS join_requests (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
}

func TestAcceptJoinRequestSetsCorrectRole(t *testing.T) {
	// Setup test database
	setupTestDBForJoinRequest(t)

	// Create test data
	clubID := uuid.New().String()
	userID := uuid.New().String()
	requestID := uuid.New().String()
	
	// Create test club
	club := Club{
		ID:          clubID,
		Name:        "Test Club",
		Description: "Test Description",
		CreatedBy:   userID,
	}
	err := database.Db.Create(&club).Error
	assert.NoError(t, err)

	// Create test user
	user := User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}
	err = database.Db.Create(&user).Error
	assert.NoError(t, err)

	// Create test join request
	joinRequest := JoinRequest{
		ID:     requestID,
		ClubID: clubID,
		Email:  "test@example.com",
	}
	err = database.Db.Create(&joinRequest).Error
	assert.NoError(t, err)

	// Accept the join request
	err = AcceptJoinRequest(requestID, userID)
	assert.NoError(t, err)

	// Verify that the member was added with the correct role
	var member Member
	err = database.Db.Where("club_id = ? AND user_id = ?", clubID, userID).First(&member).Error
	assert.NoError(t, err)
	
	// This is the key assertion - the role should be "member", not "initial"
	assert.Equal(t, "member", member.Role, "Member role should be 'member' when accepting a join request")
	
	// Verify join request was deleted
	var deletedRequest JoinRequest
	err = database.Db.Where("id = ?", requestID).First(&deletedRequest).Error
	assert.Error(t, err) // Should be not found
}