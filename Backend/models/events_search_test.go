package models

import (
	"testing"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create tables manually for SQLite compatibility
	err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT UNIQUE NOT NULL,
			keycloak_id TEXT,
			birth_date DATE,
			created_at DATETIME,
			updated_at DATETIME
		)
	`).Error
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE clubs (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT,
			deleted BOOLEAN DEFAULT FALSE,
			deleted_at DATETIME,
			deleted_by TEXT
		)
	`).Error
	if err != nil {
		t.Fatalf("Failed to create clubs table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE members (
			id TEXT PRIMARY KEY,
			club_id TEXT,
			user_id TEXT,
			role TEXT DEFAULT 'member',
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
		)
	`).Error
	if err != nil {
		t.Fatalf("Failed to create members table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE events (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			name TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
		)
	`).Error
	if err != nil {
		t.Fatalf("Failed to create events table: %v", err)
	}

	// Set the global database reference
	database.Db = db

	return db
}

func TestSearchEventsForUser(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	userID := uuid.New().String()
	club1ID := uuid.New().String()
	club2ID := uuid.New().String()

	// Create test user
	user := User{
		ID:        userID,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(&user).Error
	assert.NoError(t, err)

	// Create test clubs
	club1 := Club{
		ID:          club1ID,
		Name:        "Test Club 1",
		Description: "First test club",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	club2 := Club{
		ID:          club2ID,
		Name:        "Test Club 2",
		Description: "Second test club",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	err = db.Create(&club1).Error
	assert.NoError(t, err)
	err = db.Create(&club2).Error
	assert.NoError(t, err)

	// Add user as member of club1 only
	member := Member{
		ID:        uuid.New().String(),
		UserID:    userID,
		ClubID:    club1ID,
		Role:      "member",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	err = db.Create(&member).Error
	assert.NoError(t, err)

	// Create events in both clubs
	event1 := Event{
		ID:        uuid.New().String(),
		ClubID:    club1ID,
		Name:      "Search Test Event",
		StartTime: time.Now().Add(24 * time.Hour),
		EndTime:   time.Now().Add(26 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	event2 := Event{
		ID:        uuid.New().String(),
		ClubID:    club1ID,
		Name:      "Another Event",
		StartTime: time.Now().Add(48 * time.Hour),
		EndTime:   time.Now().Add(50 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	event3 := Event{
		ID:        uuid.New().String(),
		ClubID:    club2ID,
		Name:      "Search Private Event",
		StartTime: time.Now().Add(72 * time.Hour),
		EndTime:   time.Now().Add(74 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err = db.Create(&event1).Error
	assert.NoError(t, err)
	err = db.Create(&event2).Error
	assert.NoError(t, err)
	err = db.Create(&event3).Error
	assert.NoError(t, err)

	t.Run("SearchEventsUserIsMemberOf", func(t *testing.T) {
		results, err := SearchEventsForUser(userID, "Search")
		assert.NoError(t, err)

		// Should find event1 (user is member of club1)
		// Should NOT find event3 (user is not member of club2)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, event1.ID, results[0].ID)
		assert.Equal(t, "Search Test Event", results[0].Name)
		assert.Equal(t, club1ID, results[0].ClubID)
		assert.Equal(t, "Test Club 1", results[0].ClubName)
	})

	t.Run("SearchEventsNoMatch", func(t *testing.T) {
		results, err := SearchEventsForUser(userID, "NonExistent")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("SearchEventsCaseInsensitive", func(t *testing.T) {
		results, err := SearchEventsForUser(userID, "search")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, event1.ID, results[0].ID)
	})

	t.Run("SearchEventsPartialMatch", func(t *testing.T) {
		results, err := SearchEventsForUser(userID, "Test")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, event1.ID, results[0].ID)
	})

	t.Run("SearchEventsUserNotMemberOfAnyClub", func(t *testing.T) {
		// Create a new user who is not a member of any club
		newUserID := uuid.New().String()
		newUser := User{
			ID:        newUserID,
			FirstName: "New",
			LastName:  "User",
			Email:     "new@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := db.Create(&newUser).Error
		assert.NoError(t, err)

		results, err := SearchEventsForUser(newUserID, "Search")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("SearchEventsFromDeletedClubs", func(t *testing.T) {
		// Mark club1 as deleted
		err := db.Model(&club1).Update("deleted", true).Error
		assert.NoError(t, err)

		results, err := SearchEventsForUser(userID, "Search")
		assert.NoError(t, err)

		// Should not find events from deleted clubs
		assert.Equal(t, 0, len(results))

		// Clean up - mark club as not deleted for other tests
		err = db.Model(&club1).Update("deleted", false).Error
		assert.NoError(t, err)
	})
}
