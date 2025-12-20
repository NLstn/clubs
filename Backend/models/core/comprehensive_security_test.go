package core

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

// setupComprehensiveSecurityTestDB creates a test database for comprehensive security tests
func setupComprehensiveSecurityTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create all necessary tables
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
		CREATE TABLE IF NOT EXISTS news (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
		);
		CREATE TABLE IF NOT EXISTS api_keys (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			key_hash TEXT NOT NULL UNIQUE,
			key_prefix TEXT NOT NULL,
			permissions TEXT,
			last_used_at DATETIME,
			expires_at DATETIME,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS user_privacy_settings (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			share_birth_date BOOLEAN DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS member_privacy_settings (
			id TEXT PRIMARY KEY,
			member_id TEXT NOT NULL UNIQUE,
			share_birth_date BOOLEAN DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		);
	`).Error
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	database.Db = db
}

// TestIDORVulnerabilityInNews tests for Insecure Direct Object Reference in News model
func TestIDORVulnerabilityInNews(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create test users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")

	// Create two clubs
	club1ID := uuid.New().String()
	club2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		club1ID, "Club 1", user1ID, user1ID)
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		club2ID, "Club 2", user2ID, user2ID)

	// Add users to their respective clubs
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), club1ID, user1ID, "owner", user1ID, user1ID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), club2ID, user2ID, "owner", user2ID, user2ID)

	// Create news in Club 1
	newsID := uuid.New().String()
	database.Db.Exec("INSERT INTO news (id, club_id, title, content, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		newsID, club1ID, "News 1", "Content 1", user1ID, user1ID)

	// Test 1: User2 tries to read news from Club 1 (should fail)
	t.Run("User cannot access news from club they don't belong to", func(t *testing.T) {
		news := News{ID: newsID, ClubID: club1ID}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		scopes, err := news.ODataBeforeReadEntity(ctx, req, nil)
		if err != nil {
			t.Errorf("Expected no error from hook, got: %v", err)
		}

		// Apply scopes and try to fetch
		db := database.Db
		for _, scope := range scopes {
			db = scope(db)
		}

		var result News
		err = db.Where("id = ?", newsID).First(&result).Error
		if err == nil {
			t.Error("User2 should NOT be able to access news from Club 1")
		}
	})

	// Test 2: User2 tries to update news from Club 1 (should fail)
	t.Run("User cannot update news from club they don't belong to", func(t *testing.T) {
		news := News{
			ID:     newsID,
			ClubID: club1ID,
			Title:  "Modified Title",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		err := news.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to update news from Club 1")
		}
	})

	// Test 3: User2 tries to delete news from Club 1 (should fail)
	t.Run("User cannot delete news from club they don't belong to", func(t *testing.T) {
		news := News{
			ID:     newsID,
			ClubID: club1ID,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		err := news.ODataBeforeDelete(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to delete news from Club 1")
		}
	})
}

// TestAPIKeyAuthorizationIsolation tests that API keys properly isolate user access
func TestAPIKeyAuthorizationIsolation(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create two users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")

	// Create API key for user1
	apiKey1ID := uuid.New().String()
	database.Db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active) VALUES (?, ?, ?, ?, ?, ?)", 
		apiKey1ID, user1ID, "Test Key 1", "hash1", "sk_test", true)

	// Create API key for user2
	apiKey2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, is_active) VALUES (?, ?, ?, ?, ?, ?)", 
		apiKey2ID, user2ID, "Test Key 2", "hash2", "sk_prod", true)

	// Test: User1 tries to update User2's API key (should fail)
	t.Run("User cannot update another user's API key", func(t *testing.T) {
		apiKey := APIKey{
			ID:     apiKey2ID,
			UserID: user2ID, // This is user2's key
			Name:   "Modified Name",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		err := apiKey.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("User1 should NOT be able to update User2's API key")
		}
	})

	// Test: User1 tries to delete User2's API key (should fail)
	t.Run("User cannot delete another user's API key", func(t *testing.T) {
		apiKey := APIKey{
			ID:     apiKey2ID,
			UserID: user2ID,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		err := apiKey.ODataBeforeDelete(ctx, req)
		if err == nil {
			t.Error("User1 should NOT be able to delete User2's API key")
		}
	})
}

// TestPrivacySettingsIsolation tests that users cannot access other users' privacy settings
func TestPrivacySettingsIsolation(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create two users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")

	// Create privacy settings for user1
	privacy1ID := uuid.New().String()
	database.Db.Exec("INSERT INTO user_privacy_settings (id, user_id, share_birth_date) VALUES (?, ?, ?)", 
		privacy1ID, user1ID, false)

	// Test 1: User2 tries to read User1's privacy settings (should fail)
	t.Run("User cannot access another user's privacy settings", func(t *testing.T) {
		settings := UserPrivacySettings{ID: privacy1ID, UserID: user1ID}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		scopes, err := settings.ODataBeforeReadEntity(ctx, req, nil)
		if err != nil {
			t.Errorf("Expected no error from hook, got: %v", err)
		}

		// Apply scopes and try to fetch
		db := database.Db
		for _, scope := range scopes {
			db = scope(db)
		}

		var result UserPrivacySettings
		err = db.Where("id = ?", privacy1ID).First(&result).Error
		if err == nil {
			t.Error("User2 should NOT be able to access User1's privacy settings")
		}
	})

	// Test 2: User2 tries to update User1's privacy settings (should fail)
	t.Run("User cannot update another user's privacy settings", func(t *testing.T) {
		settings := UserPrivacySettings{
			ID:             privacy1ID,
			UserID:         user1ID,
			ShareBirthDate: true,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		err := settings.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to update User1's privacy settings")
		}
	})

	// Test 3: User2 tries to create privacy settings for User1 (should fail)
	t.Run("User cannot create privacy settings for another user", func(t *testing.T) {
		settings := UserPrivacySettings{
			UserID:         user1ID, // Trying to create for user1
			ShareBirthDate: true,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID) // but logged in as user2
		req := &http.Request{}

		err := settings.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to create privacy settings for User1")
		}
	})
}

// TestMemberPrivacySettingsIsolation tests member-specific privacy settings isolation
func TestMemberPrivacySettingsIsolation(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create two users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user1ID, "User", "One", "user1@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		user2ID, "User", "Two", "user2@test.com")

	// Create a club
	clubID := uuid.New().String()
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		clubID, "Test Club", user1ID, user1ID)

	// Add both users as members
	member1ID := uuid.New().String()
	member2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		member1ID, clubID, user1ID, "owner", user1ID, user1ID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		member2ID, clubID, user2ID, "member", user1ID, user1ID)

	// Create member privacy settings for member1
	privacy1ID := uuid.New().String()
	database.Db.Exec("INSERT INTO member_privacy_settings (id, member_id, share_birth_date) VALUES (?, ?, ?)", 
		privacy1ID, member1ID, false)

	// Test: User2 tries to update User1's member privacy settings (should fail)
	t.Run("User cannot update another member's privacy settings", func(t *testing.T) {
		settings := MemberPrivacySettings{
			ID:             privacy1ID,
			MemberID:       member1ID, // This is user1's member record
			ShareBirthDate: true,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		err := settings.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to update User1's member privacy settings")
		}
	})

	// Test: User2 tries to create member privacy settings for User1's member record (should fail)
	t.Run("User cannot create member privacy settings for another user's member", func(t *testing.T) {
		settings := MemberPrivacySettings{
			MemberID:       member1ID, // Trying to create for user1's member
			ShareBirthDate: true,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID) // but logged in as user2
		req := &http.Request{}

		err := settings.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to create member privacy settings for User1's member")
		}
	})
}

// TestMemberCannotPromoteThemselves tests that regular members cannot change their own role
func TestMemberCannotPromoteThemselves(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create users
	ownerID := uuid.New().String()
	memberID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		ownerID, "Owner", "User", "owner@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		memberID, "Member", "User", "member@test.com")

	// Create club
	clubID := uuid.New().String()
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		clubID, "Test Club", ownerID, ownerID)

	// Add members
	ownerMemberID := uuid.New().String()
	regularMemberID := uuid.New().String()
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		ownerMemberID, clubID, ownerID, "owner", ownerID, ownerID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		regularMemberID, clubID, memberID, "member", ownerID, ownerID)

	// Test: Regular member tries to promote themselves to admin (should fail)
	t.Run("Regular member cannot promote themselves", func(t *testing.T) {
		member := Member{
			ID:     regularMemberID,
			ClubID: clubID,
			UserID: memberID,
			Role:   "admin", // Trying to promote to admin
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, memberID)
		req := &http.Request{}

		err := member.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("Regular member should NOT be able to promote themselves to admin")
		}
	})
}

// TestClubIsolationInMemberQueries tests that users cannot query members across club boundaries
func TestClubIsolationInMemberQueries(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create users
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

	// user1 is member of club1
	member1ID := uuid.New().String()
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		member1ID, club1ID, user1ID, "owner", user1ID, user1ID)

	// user2 is member of club2
	member2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		member2ID, club2ID, user2ID, "owner", user2ID, user2ID)

	// user3 is member of club2
	member3ID := uuid.New().String()
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		member3ID, club2ID, user3ID, "member", user2ID, user2ID)

	// Test: user1 tries to read member list (should only see club1 members)
	t.Run("User can only see members of their own clubs", func(t *testing.T) {
		member := Member{}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		scopes, err := member.ODataBeforeReadCollection(ctx, req, nil)
		if err != nil {
			t.Fatalf("Expected no error from hook, got: %v", err)
		}

		// Apply scopes and try to fetch all members
		db := database.Db
		for _, scope := range scopes {
			db = scope(db)
		}

		var results []Member
		err = db.Find(&results).Error
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}

		// user1 should only see members from club1 (just themselves)
		if len(results) != 1 {
			t.Errorf("Expected 1 member visible to user1, got %d", len(results))
		}

		if len(results) > 0 && results[0].ClubID != club1ID {
			t.Error("user1 can see members from other clubs!")
		}
	})
}

// TestNewsCreationClubAuthorization tests that news can only be created by authorized club members
func TestNewsCreationClubAuthorization(t *testing.T) {
	setupComprehensiveSecurityTestDB(t)

	// Create users
	ownerID := uuid.New().String()
	adminID := uuid.New().String()
	memberID := uuid.New().String()
	outsiderID := uuid.New().String()
	
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		ownerID, "Owner", "User", "owner@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		adminID, "Admin", "User", "admin@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		memberID, "Member", "User", "member@test.com")
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)", 
		outsiderID, "Outsider", "User", "outsider@test.com")

	// Create club
	clubID := uuid.New().String()
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)", 
		clubID, "Test Club", ownerID, ownerID)

	// Add members with different roles
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), clubID, ownerID, "owner", ownerID, ownerID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), clubID, adminID, "admin", ownerID, ownerID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", 
		uuid.New().String(), clubID, memberID, "member", ownerID, ownerID)
	// outsiderID is NOT added to the club

	// Test 1: Outsider tries to create news (should fail)
	t.Run("Non-member cannot create news", func(t *testing.T) {
		news := News{
			ClubID:  clubID,
			Title:   "Test News",
			Content: "Test Content",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, outsiderID)
		req := &http.Request{}

		err := news.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("Outsider should NOT be able to create news for the club")
		}
	})

	// Test 2: Regular member tries to create news (should be based on club policy)
	// According to the authorization model, only admins/owners should create news
	t.Run("Regular member cannot create news", func(t *testing.T) {
		news := News{
			ClubID:  clubID,
			Title:   "Test News",
			Content: "Test Content",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, memberID)
		req := &http.Request{}

		err := news.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("Regular member should NOT be able to create news")
		}
	})

	// Test 3: Admin can create news (should succeed)
	t.Run("Admin can create news", func(t *testing.T) {
		news := News{
			ClubID:  clubID,
			Title:   "Test News",
			Content: "Test Content",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req := &http.Request{}

		err := news.ODataBeforeCreate(ctx, req)
		if err != nil {
			t.Errorf("Admin should be able to create news, got error: %v", err)
		}
	})
}
