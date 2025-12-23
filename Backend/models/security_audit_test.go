package models

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/database"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupSecurityTestDB creates a test database with all required models
func setupSecurityTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// For SQLite, we need to disable UUID generation in the default clause
	// Run migrations with minimal schema
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
		CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
		);
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			team_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			location TEXT,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT
		);
		CREATE TABLE IF NOT EXISTS fines (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			team_id TEXT,
			user_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			amount REAL NOT NULL,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT,
			paid BOOLEAN DEFAULT 0
		);
	`).Error
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Set the global database
	database.Db = db
}

// TestEventCreationTeamIDClubIDMismatch tests a critical security vulnerability
// where an attacker could create an event with a TeamID from a different club
func TestEventCreationTeamIDClubIDMismatch(t *testing.T) {
	setupSecurityTestDB(t)
	db := database.Db

	// Create two clubs
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	// Create users
	user1 := User{ID: user1ID, FirstName: "User", LastName: "One", Email: "user1@test.com"}
	user2 := User{ID: user2ID, FirstName: "User", LastName: "Two", Email: "user2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Club A (user1 is owner)
	clubA := Club{
		ID:        uuid.New().String(),
		Name:      "Club A",
		CreatedBy: user1ID,
		UpdatedBy: user1ID,
	}
	db.Create(&clubA)

	// Club B (user2 is owner)
	clubB := Club{
		ID:        uuid.New().String(),
		Name:      "Club B",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&clubB)

	// Add user1 as owner of Club A
	memberA := Member{
		ID:        uuid.New().String(),
		ClubID:    clubA.ID,
		UserID:    user1ID,
		Role:      "owner",
		CreatedBy: user1ID,
		UpdatedBy: user1ID,
	}
	db.Create(&memberA)

	// Add user2 as owner of Club B
	memberB := Member{
		ID:        uuid.New().String(),
		ClubID:    clubB.ID,
		UserID:    user2ID,
		Role:      "owner",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&memberB)

	// Create a team in Club B
	teamB := Team{
		ID:        uuid.New().String(),
		ClubID:    clubB.ID,
		Name:      "Team B",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&teamB)

	// Now try to create an event as user1 (admin of Club A)
	// with ClubID = A but TeamID = teamB (which belongs to Club B)
	// This should be REJECTED but currently may be ALLOWED
	event := Event{
		ID:          uuid.New().String(),
		ClubID:      clubA.ID,  // Club A
		TeamID:      &teamB.ID, // Team from Club B - SECURITY ISSUE!
		Name:        "Malicious Event",
		Description: strPtr("This event crosses club boundaries"),
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
	}

	ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
	req, _ := http.NewRequest("POST", "/api/v2/Events", nil)
	req = req.WithContext(ctx)

	err := event.ODataBeforeCreate(ctx, req)

	// The authorization should FAIL because TeamID belongs to a different club
	// If this passes, it's a CRITICAL security vulnerability
	if err == nil {
		t.Error("CRITICAL SECURITY VULNERABILITY: User from Club A can create events with Teams from Club B!")
		t.Error("This allows cross-club data manipulation and breaks club isolation")
	} else {
		t.Log("Security check passed: Cross-club team assignment prevented")
	}
}

// TestFineCreationTeamIDClubIDMismatch tests the same vulnerability for fines
func TestFineCreationTeamIDClubIDMismatch(t *testing.T) {
	setupSecurityTestDB(t)
	db := database.Db

	// Create two clubs
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	// Create users
	user1 := User{ID: user1ID, FirstName: "User", LastName: "One", Email: "user1@test.com"}
	user2 := User{ID: user2ID, FirstName: "User", LastName: "Two", Email: "user2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Club A (user1 is owner)
	clubA := Club{
		ID:        uuid.New().String(),
		Name:      "Club A",
		CreatedBy: user1ID,
		UpdatedBy: user1ID,
	}
	db.Create(&clubA)

	// Club B (user2 is owner)
	clubB := Club{
		ID:        uuid.New().String(),
		Name:      "Club B",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&clubB)

	// Add user1 as owner of Club A
	memberA := Member{
		ID:        uuid.New().String(),
		ClubID:    clubA.ID,
		UserID:    user1ID,
		Role:      "owner",
		CreatedBy: user1ID,
		UpdatedBy: user1ID,
	}
	db.Create(&memberA)

	// Add user2 as owner of Club B
	memberB := Member{
		ID:        uuid.New().String(),
		ClubID:    clubB.ID,
		UserID:    user2ID,
		Role:      "owner",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&memberB)

	// Create a team in Club B
	teamB := Team{
		ID:        uuid.New().String(),
		ClubID:    clubB.ID,
		Name:      "Team B",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&teamB)

	// Now try to create a fine as user1 (admin of Club A)
	// with ClubID = A but TeamID = teamB (which belongs to Club B)
	fine := Fine{
		ID:     uuid.New().String(),
		ClubID: clubA.ID,  // Club A
		TeamID: &teamB.ID, // Team from Club B - SECURITY ISSUE!
		UserID: user1ID,
		Reason: "Malicious Fine",
		Amount: 100.0,
	}

	ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
	req, _ := http.NewRequest("POST", "/api/v2/Fines", nil)
	req = req.WithContext(ctx)

	err := fine.ODataBeforeCreate(ctx, req)

	// The authorization should FAIL because TeamID belongs to a different club
	if err == nil {
		t.Error("CRITICAL SECURITY VULNERABILITY: User from Club A can create fines with Teams from Club B!")
		t.Error("This allows cross-club data manipulation and breaks club isolation")
	} else {
		t.Log("Security check passed: Cross-club team assignment prevented")
	}
}

// TestShiftCreationEventIDClubIDMismatch tests the vulnerability for shifts
func TestShiftCreationEventIDClubIDMismatch(t *testing.T) {
	setupSecurityTestDB(t)
	db := database.Db

	// Create two clubs
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()

	// Create users
	user1 := User{ID: user1ID, FirstName: "User", LastName: "One", Email: "user1@test.com"}
	user2 := User{ID: user2ID, FirstName: "User", LastName: "Two", Email: "user2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Club A (user1 is owner)
	clubA := Club{
		ID:        uuid.New().String(),
		Name:      "Club A",
		CreatedBy: user1ID,
		UpdatedBy: user1ID,
	}
	db.Create(&clubA)

	// Club B (user2 is owner)
	clubB := Club{
		ID:        uuid.New().String(),
		Name:      "Club B",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&clubB)

	// Add user1 as owner of Club A
	memberA := Member{
		ID:        uuid.New().String(),
		ClubID:    clubA.ID,
		UserID:    user1ID,
		Role:      "owner",
		CreatedBy: user1ID,
		UpdatedBy: user1ID,
	}
	db.Create(&memberA)

	// Add user2 as owner of Club B
	memberB := Member{
		ID:        uuid.New().String(),
		ClubID:    clubB.ID,
		UserID:    user2ID,
		Role:      "owner",
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&memberB)

	// Create an event in Club B
	eventB := Event{
		ID:        uuid.New().String(),
		ClubID:    clubB.ID,
		Name:      "Event B",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		CreatedBy: user2ID,
		UpdatedBy: user2ID,
	}
	db.Create(&eventB)

	// Now try to create a shift as user1 (admin of Club A)
	// with ClubID = A but EventID = eventB (which belongs to Club B)
	shift := Shift{
		ID:        uuid.New().String(),
		ClubID:    clubA.ID,  // Club A
		EventID:   eventB.ID, // Event from Club B - SECURITY ISSUE!
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	ctx := context.WithValue(context.Background(), auth.UserIDKey, user1ID)
	req, _ := http.NewRequest("POST", "/api/v2/Shifts", nil)
	req = req.WithContext(ctx)

	err := shift.ODataBeforeCreate(ctx, req)

	// The authorization should FAIL because EventID belongs to a different club
	if err == nil {
		t.Error("CRITICAL SECURITY VULNERABILITY: User from Club A can create shifts with Events from Club B!")
		t.Error("This allows cross-club data manipulation and breaks club isolation")
	} else {
		t.Log("Security check passed: Cross-club event assignment prevented")
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}

// TestPrivilegeEscalationViaRoleUpdate tests if admins can improperly escalate privileges
func TestPrivilegeEscalationViaRoleUpdate(t *testing.T) {
	setupSecurityTestDB(t)
	db := database.Db

	// Create club and users
	ownerID := uuid.New().String()
	adminID := uuid.New().String()
	memberID := uuid.New().String()

	// Create users
	owner := User{ID: ownerID, FirstName: "Owner", LastName: "User", Email: "owner@test.com"}
	admin := User{ID: adminID, FirstName: "Admin", LastName: "User", Email: "admin@test.com"}
	member := User{ID: memberID, FirstName: "Member", LastName: "User", Email: "member@test.com"}
	db.Create(&owner)
	db.Create(&admin)
	db.Create(&member)

	// Create club
	club := Club{
		ID:        uuid.New().String(),
		Name:      "Test Club",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&club)

	// Add owner
	ownerMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    ownerID,
		Role:      "owner",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&ownerMember)

	// Add admin
	adminMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    adminID,
		Role:      "admin",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&adminMember)

	// Add regular member
	regularMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    memberID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&regularMember)

	// TEST 1: Admin tries to promote themselves to owner (should FAIL)
	t.Run("Admin cannot promote themselves to owner", func(t *testing.T) {
		updatedMember := adminMember
		updatedMember.Role = "owner" // Try to escalate to owner

		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req, _ := http.NewRequest("PATCH", "/api/v2/Members", nil)
		req = req.WithContext(ctx)

		err := updatedMember.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("SECURITY VULNERABILITY: Admin was able to promote themselves to owner!")
		} else {
			t.Log("Security check passed: Admin cannot self-promote to owner")
		}
	})

	// TEST 2: Admin tries to promote regular member to owner (should FAIL)
	t.Run("Admin cannot promote member to owner", func(t *testing.T) {
		updatedMember := regularMember
		updatedMember.Role = "owner" // Try to make member an owner

		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req, _ := http.NewRequest("PATCH", "/api/v2/Members", nil)
		req = req.WithContext(ctx)

		err := updatedMember.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("SECURITY VULNERABILITY: Admin was able to promote member to owner!")
		} else {
			t.Log("Security check passed: Admin cannot promote members to owner")
		}
	})

	// TEST 3: Admin tries to demote owner (should FAIL)
	t.Run("Admin cannot demote owner", func(t *testing.T) {
		updatedMember := ownerMember
		updatedMember.Role = "member" // Try to demote owner

		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req, _ := http.NewRequest("PATCH", "/api/v2/Members", nil)
		req = req.WithContext(ctx)

		err := updatedMember.ODataBeforeUpdate(ctx, req)
		if err == nil {
			t.Error("SECURITY VULNERABILITY: Admin was able to demote owner!")
		} else {
			t.Log("Security check passed: Admin cannot demote owner")
		}
	})

	// TEST 4: Owner CAN promote admin to owner (should PASS)
	t.Run("Owner can promote admin to owner", func(t *testing.T) {
		updatedMember := adminMember
		updatedMember.Role = "owner"

		ctx := context.WithValue(context.Background(), auth.UserIDKey, ownerID)
		req, _ := http.NewRequest("PATCH", "/api/v2/Members", nil)
		req = req.WithContext(ctx)

		err := updatedMember.ODataBeforeUpdate(ctx, req)
		if err != nil {
			t.Errorf("Owner should be able to promote admin to owner, but got error: %v", err)
		} else {
			t.Log("Security check passed: Owner can promote to owner")
		}
	})
}

// TestPrivilegeEscalationViaCreate tests if users can create members with improper roles
func TestPrivilegeEscalationViaCreate(t *testing.T) {
	setupSecurityTestDB(t)
	db := database.Db

	// Create club and users
	ownerID := uuid.New().String()
	adminID := uuid.New().String()
	newUserID := uuid.New().String()

	// Create users
	owner := User{ID: ownerID, FirstName: "Owner", LastName: "User", Email: "owner@test.com"}
	admin := User{ID: adminID, FirstName: "Admin", LastName: "User", Email: "admin@test.com"}
	newUser := User{ID: newUserID, FirstName: "New", LastName: "User", Email: "new@test.com"}
	db.Create(&owner)
	db.Create(&admin)
	db.Create(&newUser)

	// Create club
	club := Club{
		ID:        uuid.New().String(),
		Name:      "Test Club",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&club)

	// Add owner
	ownerMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    ownerID,
		Role:      "owner",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&ownerMember)

	// Add admin
	adminMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    adminID,
		Role:      "admin",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&adminMember)

	// TEST 1: Admin tries to create a new owner (should FAIL)
	t.Run("Admin cannot create new owner", func(t *testing.T) {
		newMember := Member{
			ID:     uuid.New().String(),
			ClubID: club.ID,
			UserID: newUserID,
			Role:   "owner", // Admin trying to create an owner
		}

		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req, _ := http.NewRequest("POST", "/api/v2/Members", nil)
		req = req.WithContext(ctx)

		err := newMember.ODataBeforeCreate(ctx, req)
		if err == nil {
			t.Error("SECURITY VULNERABILITY: Admin was able to create a new owner!")
		} else {
			t.Log("Security check passed: Admin cannot create owners")
		}
	})

	// TEST 2: Owner CAN create a new owner (should PASS)
	t.Run("Owner can create new owner", func(t *testing.T) {
		newMember := Member{
			ID:     uuid.New().String(),
			ClubID: club.ID,
			UserID: newUserID,
			Role:   "owner",
		}

		ctx := context.WithValue(context.Background(), auth.UserIDKey, ownerID)
		req, _ := http.NewRequest("POST", "/api/v2/Members", nil)
		req = req.WithContext(ctx)

		err := newMember.ODataBeforeCreate(ctx, req)
		if err != nil {
			t.Errorf("Owner should be able to create new owner, but got error: %v", err)
		} else {
			t.Log("Security check passed: Owner can create owners")
		}
	})
}
