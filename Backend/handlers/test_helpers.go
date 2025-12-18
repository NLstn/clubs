package handlers

import (
	"os"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDatabase holds the test database instance
var testDB *gorm.DB

// SetupTestDB initializes an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) {
	// Set test environment variable
	os.Setenv("GO_ENV", "test")

	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set the global database reference for the application
	database.Db = testDB

	// Set up SQLite-compatible tables
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS magic_links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT NOT NULL UNIQUE,
			keycloak_id TEXT UNIQUE,
			birth_date DATE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME,
			user_agent TEXT,
			ip_address TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS clubs (
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
		)
	`)
	testDB.Exec(`
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
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS join_requests (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fines (
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
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fine_templates (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			description TEXT,
			amount REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS shifts (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			event_id TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS shift_members (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			shift_id TEXT NOT NULL,
			user_id TEXT NOT NULL
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			team_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			location TEXT,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT,
			is_recurring BOOLEAN DEFAULT FALSE,
			recurrence_pattern TEXT,
			recurrence_interval INTEGER DEFAULT 1,
			recurrence_end DATETIME,
			parent_event_id TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS event_rsvps (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			response TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS news (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS club_settings (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL UNIQUE,
			fines_enabled BOOLEAN DEFAULT TRUE,
			shifts_enabled BOOLEAN DEFAULT TRUE,
			teams_enabled BOOLEAN DEFAULT TRUE,
			members_list_visible BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			read BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			club_id TEXT,
			event_id TEXT,
			fine_id TEXT,
			invite_id TEXT,
			join_request_id TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS activities (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			actor_id TEXT,
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			metadata TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS user_notification_preferences (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			member_added_in_app BOOLEAN DEFAULT TRUE,
			member_added_email BOOLEAN DEFAULT TRUE,
			invite_received_in_app BOOLEAN DEFAULT TRUE,
			invite_received_email BOOLEAN DEFAULT TRUE,
			event_created_in_app BOOLEAN DEFAULT TRUE,
			event_created_email BOOLEAN DEFAULT FALSE,
			fine_assigned_in_app BOOLEAN DEFAULT TRUE,
			fine_assigned_email BOOLEAN DEFAULT TRUE,
			news_created_in_app BOOLEAN DEFAULT TRUE,
			news_created_email BOOLEAN DEFAULT FALSE,
			role_changed_in_app BOOLEAN DEFAULT TRUE,
			role_changed_email BOOLEAN DEFAULT TRUE,
			join_request_in_app BOOLEAN DEFAULT TRUE,
			join_request_email BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS user_privacy_settings (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			share_birth_date BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS member_privacy_settings (
			id TEXT PRIMARY KEY,
			member_id TEXT NOT NULL UNIQUE,
			share_birth_date BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT,
			deleted BOOLEAN DEFAULT FALSE,
			deleted_at DATETIME,
			deleted_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT DEFAULT 'member',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T) {
	if testDB != nil {
		// Clear all data from tables to ensure clean state
		testDB.Exec("DELETE FROM activities")
		testDB.Exec("DELETE FROM refresh_tokens")
		testDB.Exec("DELETE FROM magic_links")
		testDB.Exec("DELETE FROM user_notification_preferences")
		testDB.Exec("DELETE FROM user_privacy_settings")
		testDB.Exec("DELETE FROM member_privacy_settings")
		testDB.Exec("DELETE FROM notifications")
		testDB.Exec("DELETE FROM fines")
		testDB.Exec("DELETE FROM members")
		testDB.Exec("DELETE FROM events")
		testDB.Exec("DELETE FROM clubs")
		testDB.Exec("DELETE FROM users")

		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// CreateTestUser creates a test user and returns it with an access token
func CreateTestUser(t *testing.T, email string) (models.User, string) {
	// Generate a UUID-like string for SQLite
	userID := uuid.New().String()
	keycloakID := uuid.New().String() // Generate unique KeycloakID for test users

	// Create user directly in database
	keycloakPtr := &keycloakID
	user := models.User{
		ID:         userID,
		Email:      email,
		FirstName:  "Test",
		LastName:   "User",
		KeycloakID: keycloakPtr,
	}

	if err := testDB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	return user, accessToken
}

// CreateTestClub creates a test club with the given user as owner
func CreateTestClub(t *testing.T, user models.User, clubName string) models.Club {
	clubID := uuid.New().String()
	description := "Test club description"

	club := models.Club{
		ID:          clubID,
		Name:        clubName,
		Description: &description,
	}

	if err := testDB.Create(&club).Error; err != nil {
		t.Fatalf("Failed to create test club: %v", err)
	}

	// Add the owner as a member with owner role
	memberID := uuid.New().String()
	member := models.Member{
		ID:     memberID,
		UserID: user.ID,
		ClubID: club.ID,
		Role:   "owner",
	}
	if err := testDB.Create(&member).Error; err != nil {
		t.Fatalf("Failed to add owner as member: %v", err)
	}

	return club
}

// CreateTestMember creates a test member directly in the database without notifications
func CreateTestMember(t *testing.T, user models.User, club models.Club, role string) models.Member {
	memberID := uuid.New().String()

	member := models.Member{
		ID:        memberID,
		UserID:    user.ID,
		ClubID:    club.ID,
		Role:      role,
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}

	if err := testDB.Create(&member).Error; err != nil {
		t.Fatalf("Failed to create test member: %v", err)
	}

	return member
}
