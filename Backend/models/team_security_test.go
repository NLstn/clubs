package models

import (
	"context"
	"net/http"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTeamSecurityTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// For SQLite, we need to create tables manually without UUID generation
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
			updated_by TEXT,
			FOREIGN KEY (club_id) REFERENCES clubs(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT,
			FOREIGN KEY (club_id) REFERENCES clubs(id)
		);
		CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT DEFAULT 'member',
			created_at DATETIME,
			created_by TEXT,
			updated_at DATETIME,
			updated_by TEXT,
			FOREIGN KEY (team_id) REFERENCES teams(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`).Error
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	database.Db = db
}

// TestTeamMemberPrivilegeEscalationViaRoleUpdate tests if team members can improperly escalate privileges via direct OData PATCH
func TestTeamMemberPrivilegeEscalationViaRoleUpdate(t *testing.T) {
	setupTeamSecurityTestDB(t)
	db := database.Db

	// Create club and users
	ownerID := uuid.New().String()
	adminID := uuid.New().String()
	memberID := uuid.New().String()
	teamAdminID := uuid.New().String()
	teamMemberID := uuid.New().String()

	// Create users
	owner := User{ID: ownerID, FirstName: "Owner", LastName: "User", Email: "owner@test.com"}
	admin := User{ID: adminID, FirstName: "Admin", LastName: "User", Email: "admin@test.com"}
	member := User{ID: memberID, FirstName: "Member", LastName: "User", Email: "member@test.com"}
	teamAdmin := User{ID: teamAdminID, FirstName: "TeamAdmin", LastName: "User", Email: "teamadmin@test.com"}
	teamMember := User{ID: teamMemberID, FirstName: "TeamMember", LastName: "User", Email: "teammember@test.com"}
	db.Create(&owner)
	db.Create(&admin)
	db.Create(&member)
	db.Create(&teamAdmin)
	db.Create(&teamMember)

	// Create club
	club := Club{
		ID:        uuid.New().String(),
		Name:      "Test Club",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&club)

	// Add club members
	ownerMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    ownerID,
		Role:      "owner",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&ownerMember)

	clubAdmin := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    adminID,
		Role:      "admin",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&clubAdmin)

	regularMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    memberID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&regularMember)

	// Create team
	team := Team{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		Name:      "Test Team",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&team)

	// Add team members
	teamAdminMember := TeamMember{
		ID:        uuid.New().String(),
		TeamID:    team.ID,
		UserID:    teamAdminID,
		Role:      "admin",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&teamAdminMember)

	regularTeamMember := TeamMember{
		ID:        uuid.New().String(),
		TeamID:    team.ID,
		UserID:    teamMemberID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&regularTeamMember)

	// Also add teamAdmin and teamMember to club as regular members
	db.Create(&Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    teamAdminID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	})
	db.Create(&Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    teamMemberID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	})

	// Test Case 1: Regular team member trying to promote themselves to team admin via PATCH
	// This should FAIL
	t.Run("TeamMember cannot self-promote to admin", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, teamMemberID)
		req, _ := http.NewRequestWithContext(ctx, "PATCH", "/api/v2/TeamMembers('"+regularTeamMember.ID+"')", nil)

		// Attempt to promote self to admin
		updatedTeamMember := regularTeamMember
		updatedTeamMember.Role = "admin"

		err := updatedTeamMember.ODataBeforeUpdate(ctx, req)
		assert.Error(t, err, "Regular team member should not be able to promote themselves to admin")
		assert.Contains(t, err.Error(), "unauthorized", "Error should indicate authorization failure")
	})

	// Test Case 2: Club admin trying to update role via PATCH (should work with proper validation)
	t.Run("Club admin can change team member roles", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, adminID)
		req, _ := http.NewRequestWithContext(ctx, "PATCH", "/api/v2/TeamMembers('"+regularTeamMember.ID+"')", nil)

		// Club admin attempts to promote team member to admin
		updatedTeamMember := regularTeamMember
		updatedTeamMember.Role = "admin"

		err := updatedTeamMember.ODataBeforeUpdate(ctx, req)
		// This should succeed as club admins can change team roles
		assert.NoError(t, err, "Club admin should be able to change team member roles")
	})

	// Test Case 3: Team admin trying to promote another member via PATCH (should work with proper validation)
	t.Run("Team admin can change team member roles", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, teamAdminID)
		req, _ := http.NewRequestWithContext(ctx, "PATCH", "/api/v2/TeamMembers('"+regularTeamMember.ID+"')", nil)

		// Team admin attempts to promote team member to admin
		updatedTeamMember := regularTeamMember
		updatedTeamMember.Role = "admin"

		err := updatedTeamMember.ODataBeforeUpdate(ctx, req)
		// This should succeed as team admins can change team roles
		assert.NoError(t, err, "Team admin should be able to change team member roles")
	})

	// Test Case 4: Prevent demoting the last team admin
	t.Run("Cannot demote last team admin", func(t *testing.T) {
		// Remove the regular team member who might have been promoted
		db.Delete(&TeamMember{}, "id = ?", regularTeamMember.ID)

		ctx := context.WithValue(context.Background(), auth.UserIDKey, ownerID)
		req, _ := http.NewRequestWithContext(ctx, "PATCH", "/api/v2/TeamMembers('"+teamAdminMember.ID+"')", nil)

		// Attempt to demote the only team admin
		updatedTeamMember := teamAdminMember
		updatedTeamMember.Role = "member"

		err := updatedTeamMember.ODataBeforeUpdate(ctx, req)
		assert.Error(t, err, "Should not be able to demote the last team admin")
		assert.Contains(t, err.Error(), "last", "Error should mention last admin")
	})
}

// TestTeamMemberPrivilegeEscalationViaCreate tests if users can create team members with improper roles
func TestTeamMemberPrivilegeEscalationViaCreate(t *testing.T) {
	setupTeamSecurityTestDB(t)
	db := database.Db

	// Create club and users
	ownerID := uuid.New().String()
	memberID := uuid.New().String()
	newUserID := uuid.New().String()

	// Create users
	owner := User{ID: ownerID, FirstName: "Owner", LastName: "User", Email: "owner@test.com"}
	member := User{ID: memberID, FirstName: "Member", LastName: "User", Email: "member@test.com"}
	newUser := User{ID: newUserID, FirstName: "New", LastName: "User", Email: "new@test.com"}
	db.Create(&owner)
	db.Create(&member)
	db.Create(&newUser)

	// Create club
	club := Club{
		ID:        uuid.New().String(),
		Name:      "Test Club",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&club)

	// Add club members
	ownerMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    ownerID,
		Role:      "owner",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&ownerMember)

	regularMember := Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    memberID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&regularMember)

	db.Create(&Member{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		UserID:    newUserID,
		Role:      "member",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	})

	// Create team
	team := Team{
		ID:        uuid.New().String(),
		ClubID:    club.ID,
		Name:      "Test Team",
		CreatedBy: ownerID,
		UpdatedBy: ownerID,
	}
	db.Create(&team)

	// Test Case 1: Regular club member trying to add someone as team admin (should fail)
	t.Run("Regular member cannot add team admins", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, memberID)
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v2/TeamMembers", nil)

		newTeamMember := TeamMember{
			ID:     uuid.New().String(),
			TeamID: team.ID,
			UserID: newUserID,
			Role:   "admin",
		}

		err := newTeamMember.ODataBeforeCreate(ctx, req)
		assert.Error(t, err, "Regular member should not be able to add team members")
	})

	// Test Case 2: Club owner can add team members with any role
	t.Run("Club owner can add team admins", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, ownerID)
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v2/TeamMembers", nil)

		newTeamMember := TeamMember{
			ID:     uuid.New().String(),
			TeamID: team.ID,
			UserID: newUserID,
			Role:   "admin",
		}

		err := newTeamMember.ODataBeforeCreate(ctx, req)
		assert.NoError(t, err, "Club owner should be able to add team admins")
	})
}
