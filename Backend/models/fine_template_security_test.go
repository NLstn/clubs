package models

import (
	"context"
	"net/http"
	"testing"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/database"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupFineTemplateSecurityTestDB creates a test database for fine template security tests
func setupFineTemplateSecurityTestDB(t *testing.T) {
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
			created_at DATETIME,
			updated_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS clubs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT,
			deleted BOOLEAN DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS club_settings (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL UNIQUE,
			fines_enabled BOOLEAN DEFAULT 0,
			shifts_enabled BOOLEAN DEFAULT 0,
			teams_enabled BOOLEAN DEFAULT 0,
			news_enabled BOOLEAN DEFAULT 0,
			events_enabled BOOLEAN DEFAULT 0,
			members_list_visible BOOLEAN DEFAULT 0,
			discoverable_by_non_members BOOLEAN DEFAULT 0,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
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
		CREATE TABLE IF NOT EXISTS fine_templates (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			description TEXT NOT NULL,
			amount REAL NOT NULL,
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

// TestFineTemplateIDORVulnerability tests for Insecure Direct Object Reference in FineTemplate model
func TestFineTemplateIDORVulnerability(t *testing.T) {
	setupFineTemplateSecurityTestDB(t)

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

	// Create fine template in Club 1
	templateID := uuid.New().String()
	database.Db.Exec("INSERT INTO fine_templates (id, club_id, description, amount, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		templateID, club1ID, "Late to practice", 10.0, user1ID, user1ID)

	// Test 1: User2 tries to read fine template from Club 1 (should fail)
	t.Run("User cannot access fine template from club they don't belong to", func(t *testing.T) {
		template := FineTemplate{ID: templateID, ClubID: club1ID}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		scopes, err := template.ODataBeforeReadEntity(ctx, req, nil)
		if err != nil {
			t.Errorf("Expected no error from hook, got: %v", err)
		}

		// Apply scopes and try to fetch
		db := database.Db
		for _, scope := range scopes {
			db = scope(db)
		}

		var result FineTemplate
		err = db.Where("id = ?", templateID).First(&result).Error
		if err == nil {
			t.Error("User2 should NOT be able to access fine template from Club 1")
		}
	})

	// Test 2: User2 tries to update fine template from Club 1 (should fail)
	t.Run("User cannot update fine template from club they don't belong to", func(t *testing.T) {
		template := FineTemplate{
			ID:          templateID,
			ClubID:      club1ID,
			Description: "Modified Template",
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		err := template.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to update fine template from Club 1")
		}
	})

	// Test 3: User2 tries to delete fine template from Club 1 (should fail)
	t.Run("User cannot delete fine template from club they don't belong to", func(t *testing.T) {
		template := FineTemplate{
			ID:     templateID,
			ClubID: club1ID,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user2ID)
		req := &http.Request{}

		err := template.ODataBeforeDelete(ctx, req)
		if err == nil {
			t.Error("User2 should NOT be able to delete fine template from Club 1")
		}
	})
}

// TestFineTemplateCreationAuthorization tests that only admins/owners can create fine templates
func TestFineTemplateCreationAuthorization(t *testing.T) {
	setupFineTemplateSecurityTestDB(t)

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

	// Create club settings with fines enabled
	database.Db.Exec("INSERT INTO club_settings (id, club_id, fines_enabled, created_by, updated_by) VALUES (?, ?, ?, ?, ?)",
		uuid.New().String(), clubID, 1, ownerID, ownerID)

	// Add members with different roles
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), clubID, ownerID, "owner", ownerID, ownerID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), clubID, adminID, "admin", ownerID, ownerID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), clubID, memberID, "member", ownerID, ownerID)
	// outsiderID is NOT added to the club

	// Test 1: Outsider tries to create fine template (should fail)
	t.Run("Non-member cannot create fine template", func(t *testing.T) {
		template := FineTemplate{
			ClubID:      clubID,
			Description: "Test Template",
			Amount:      10.0,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, outsiderID)
		req := &http.Request{}

		err := template.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("Outsider should NOT be able to create fine template for the club")
		}
	})

	// Test 2: Regular member tries to create fine template (should fail)
	t.Run("Regular member cannot create fine template", func(t *testing.T) {
		template := FineTemplate{
			ClubID:      clubID,
			Description: "Test Template",
			Amount:      10.0,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, memberID)
		req := &http.Request{}

		err := template.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("Regular member should NOT be able to create fine template")
		}
	})

	// Test 3: Admin can create fine template (should succeed)
	t.Run("Admin can create fine template", func(t *testing.T) {
		template := FineTemplate{
			ClubID:      clubID,
			Description: "Test Template",
			Amount:      10.0,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req := &http.Request{}

		err := template.ODataBeforeCreate(ctx, req)
		if err != nil {
			t.Errorf("Admin should be able to create fine template, got error: %v", err)
		}
	})

	// Test 4: Owner can create fine template (should succeed)
	t.Run("Owner can create fine template", func(t *testing.T) {
		template := FineTemplate{
			ClubID:      clubID,
			Description: "Test Template",
			Amount:      10.0,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, ownerID)
		req := &http.Request{}

		err := template.ODataBeforeCreate(ctx, req)
		if err != nil {
			t.Errorf("Owner should be able to create fine template, got error: %v", err)
		}
	})
}

// TestFineTemplateClubIsolation tests that users can only see templates from their clubs
func TestFineTemplateClubIsolation(t *testing.T) {
	setupFineTemplateSecurityTestDB(t)

	// Create users
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

	// Create club settings with fines enabled
	database.Db.Exec("INSERT INTO club_settings (id, club_id, fines_enabled, created_by, updated_by) VALUES (?, ?, ?, ?, ?)",
		uuid.New().String(), club1ID, 1, user1ID, user1ID)
	database.Db.Exec("INSERT INTO club_settings (id, club_id, fines_enabled, created_by, updated_by) VALUES (?, ?, ?, ?, ?)",
		uuid.New().String(), club2ID, 1, user2ID, user2ID)

	// user1 is member of club1
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), club1ID, user1ID, "owner", user1ID, user1ID)

	// user2 is member of club2
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), club2ID, user2ID, "owner", user2ID, user2ID)

	// Create templates in both clubs
	template1ID := uuid.New().String()
	template2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO fine_templates (id, club_id, description, amount, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		template1ID, club1ID, "Club 1 Template", 10.0, user1ID, user1ID)
	database.Db.Exec("INSERT INTO fine_templates (id, club_id, description, amount, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		template2ID, club2ID, "Club 2 Template", 15.0, user2ID, user2ID)

	// Test: user1 tries to read fine template list (should only see club1 templates)
	t.Run("User can only see fine templates from their own clubs", func(t *testing.T) {
		template := FineTemplate{}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
		req := &http.Request{}

		scopes, err := template.ODataBeforeReadCollection(ctx, req, nil)
		if err != nil {
			t.Fatalf("Expected no error from hook, got: %v", err)
		}

		// Apply scopes and try to fetch all templates
		db := database.Db
		for _, scope := range scopes {
			db = scope(db)
		}

		var results []FineTemplate
		err = db.Find(&results).Error
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}

		// user1 should only see templates from club1
		if len(results) != 1 {
			t.Errorf("Expected 1 template visible to user1, got %d", len(results))
		}

		if len(results) > 0 && results[0].ClubID != club1ID {
			t.Error("user1 can see templates from other clubs!")
		}
	})
}

// TestFineTemplateClubIDImmutable tests that ClubID cannot be changed after creation
func TestFineTemplateClubIDImmutable(t *testing.T) {
	setupFineTemplateSecurityTestDB(t)

	// Create user and two clubs
	userID := uuid.New().String()
	database.Db.Exec("INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?)",
		userID, "User", "Test", "user@test.com")

	club1ID := uuid.New().String()
	club2ID := uuid.New().String()
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)",
		club1ID, "Club 1", userID, userID)
	database.Db.Exec("INSERT INTO clubs (id, name, created_by, updated_by) VALUES (?, ?, ?, ?)",
		club2ID, "Club 2", userID, userID)

	// Create club settings with fines enabled
	database.Db.Exec("INSERT INTO club_settings (id, club_id, fines_enabled, created_by, updated_by) VALUES (?, ?, ?, ?, ?)",
		uuid.New().String(), club1ID, 1, userID, userID)
	database.Db.Exec("INSERT INTO club_settings (id, club_id, fines_enabled, created_by, updated_by) VALUES (?, ?, ?, ?, ?)",
		uuid.New().String(), club2ID, 1, userID, userID)

	// User is owner of both clubs
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), club1ID, userID, "owner", userID, userID)
	database.Db.Exec("INSERT INTO members (id, club_id, user_id, role, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), club2ID, userID, "owner", userID, userID)

	// Create template in Club 1
	templateID := uuid.New().String()
	database.Db.Exec("INSERT INTO fine_templates (id, club_id, description, amount, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
		templateID, club1ID, "Template", 10.0, userID, userID)

	// Test: Try to change ClubID to Club 2 (should fail)
	t.Run("ClubID cannot be changed after creation", func(t *testing.T) {
		template := FineTemplate{
			ID:          templateID,
			ClubID:      club2ID, // Trying to change to club2
			Description: "Updated Template",
			Amount:      15.0,
		}
		ctx := context.WithValue(context.Background(), auth.UserIDKey, userID)
		req := &http.Request{}

		err := template.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("Should NOT be able to change ClubID of an existing template")
		}
		if err != nil && err.Error() != "forbidden: club cannot be changed for an existing fine template" {
			t.Errorf("Expected immutability error, got: %v", err)
		}
	})
}
