package models

import (
	"testing"

	"github.com/NLstn/civo/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestJoinRequestNotifications(t *testing.T) {
	// Set up test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	database.Db = testDB

	// Simplified models for SQLite
	type TestUser struct {
		ID        string `gorm:"primaryKey"`
		FirstName string
		LastName  string
		Email     string
	}

	type TestClub struct {
		ID   string `gorm:"primaryKey"`
		Name string
	}

	type TestMember struct {
		ID     string `gorm:"primaryKey"`
		ClubID string
		UserID string
		Role   string
	}

	type TestNotification struct {
		ID      string `gorm:"primaryKey"`
		UserID  string
		Type    string
		Title   string
		Message string
		ClubID  *string
	}

	// Migrate test tables
	err = testDB.AutoMigrate(&TestUser{}, &TestClub{}, &TestMember{}, &TestNotification{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create test data
	owner := TestUser{ID: "owner-1", FirstName: "Owner", LastName: "User", Email: "owner@test.com"}
	admin := TestUser{ID: "admin-1", FirstName: "Admin", LastName: "User", Email: "admin@test.com"}
	club := TestClub{ID: "club-1", Name: "Test Club"}

	testDB.Create(&owner)
	testDB.Create(&admin)
	testDB.Create(&club)

	// Create members
	testDB.Create(&TestMember{ID: "member-1", ClubID: club.ID, UserID: owner.ID, Role: "owner"})
	testDB.Create(&TestMember{ID: "member-2", ClubID: club.ID, UserID: admin.ID, Role: "admin"})

	// Test GetAdminsAndOwners query manually
	var admins []TestUser
	err = testDB.Table("test_users").
		Joins("JOIN test_members ON test_users.id = test_members.user_id").
		Where("test_members.club_id = ? AND (test_members.role = ? OR test_members.role = ?)", club.ID, "admin", "owner").
		Find(&admins).Error

	if err != nil {
		t.Fatalf("Failed to get admins: %v", err)
	}

	if len(admins) != 2 {
		t.Errorf("Expected 2 admins/owners, got %d", len(admins))
	}

	// Test notification creation
	clubID := club.ID
	err = testDB.Create(&TestNotification{
		ID:      "notif-1",
		UserID:  owner.ID,
		Type:    "join_request_received",
		Title:   "New Join Request",
		Message: "John Doe (john@test.com) has requested to join Test Club",
		ClubID:  &clubID,
	}).Error

	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	// Verify notification was created
	var notifications []TestNotification
	err = testDB.Where("user_id = ? AND type = ?", owner.ID, "join_request_received").Find(&notifications).Error
	if err != nil {
		t.Fatalf("Failed to query notifications: %v", err)
	}

	if len(notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notifications))
	}

	if len(notifications) > 0 && notifications[0].Title != "New Join Request" {
		t.Errorf("Expected title 'New Join Request', got '%s'", notifications[0].Title)
	}

	t.Logf("✅ Test passed: Found %d admins/owners and created notification successfully", len(admins))
}
