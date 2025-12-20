package core_test

import (
	"encoding/json"
	"testing"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
)

func TestCreateRoleChangeActivity(t *testing.T) {
	testCases := []struct {
		name         string
		oldRole      string
		newRole      string
		expectedType string
	}{
		{"promotion", "member", "admin", "member_promoted"},
		{"demotion", "admin", "member", "member_demoted"},
		{"unchanged", "admin", "admin", "role_changed"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handlers.SetupTestDB(t)
			defer handlers.TeardownTestDB(t)
			db := database.Db

			clubID := "club-1"
			userID := "user-1"
			actorID := "actor-1"
			clubName := "Test Club"

			err := models.CreateRoleChangeActivity(clubID, userID, actorID, clubName, tc.oldRole, tc.newRole)
			if err != nil {
				t.Fatalf("CreateRoleChangeActivity returned error: %v", err)
			}

			var activity models.Activity
			if err := db.First(&activity).Error; err != nil {
				t.Fatalf("failed to fetch activity: %v", err)
			}

			if activity.Type != tc.expectedType {
				t.Errorf("expected type %s, got %s", tc.expectedType, activity.Type)
			}

			var meta map[string]interface{}
			if err := json.Unmarshal([]byte(activity.Metadata), &meta); err != nil {
				t.Fatalf("failed to unmarshal metadata: %v", err)
			}

			if meta["old_role"] != tc.oldRole {
				t.Errorf("metadata old_role expected %s, got %v", tc.oldRole, meta["old_role"])
			}
			if meta["new_role"] != tc.newRole {
				t.Errorf("metadata new_role expected %s, got %v", tc.newRole, meta["new_role"])
			}
			if meta["club_name"] != clubName {
				t.Errorf("metadata club_name expected %s, got %v", clubName, meta["club_name"])
			}
		})
	}
}

func TestCreateMemberJoinedActivity(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)
	db := database.Db

	clubID := "club-1"
	userID := "user-1"
	clubName := "Test Club"

	err := models.CreateMemberJoinedActivity(clubID, userID, clubName, nil)
	if err != nil {
		t.Fatalf("CreateMemberJoinedActivity returned error: %v", err)
	}

	var activity models.Activity
	if err := db.First(&activity).Error; err != nil {
		t.Fatalf("failed to fetch activity: %v", err)
	}

	if activity.Type != "member_joined" {
		t.Errorf("expected type 'member_joined', got %s", activity.Type)
	}

	if activity.ClubID != clubID {
		t.Errorf("expected club_id %s, got %s", clubID, activity.ClubID)
	}

	if activity.UserID != userID {
		t.Errorf("expected user_id %s, got %s", userID, activity.UserID)
	}

	if activity.Title != "" {
		t.Errorf("expected empty title, got %s", activity.Title)
	}

	if activity.Content != "" {
		t.Errorf("expected empty content, got %s", activity.Content)
	}

	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(activity.Metadata), &meta); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}

	if meta["club_name"] != clubName {
		t.Errorf("metadata club_name expected %s, got %v", clubName, meta["club_name"])
	}
}

func TestAddMemberCreatesActivity(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)
	db := database.Db

	// Create test club and users
	owner, _ := handlers.CreateTestUser(t, "owner@example.com")
	newMember, _ := handlers.CreateTestUser(t, "newmember@example.com")
	club := handlers.CreateTestClub(t, owner, "Test Club")

	// Add new member
	err := club.AddMember(newMember.ID, "member")
	if err != nil {
		t.Fatalf("AddMember returned error: %v", err)
	}

	// Verify activity was created
	var activities []models.Activity
	err = db.Where("club_id = ? AND user_id = ? AND type = ?", club.ID, newMember.ID, "member_joined").Find(&activities).Error
	if err != nil {
		t.Fatalf("failed to fetch activities: %v", err)
	}

	if len(activities) == 0 {
		t.Fatalf("expected activity to be created, but none found")
	}

	// Verify activity details
	activity := activities[0]
	if activity.Type != "member_joined" {
		t.Errorf("expected type 'member_joined', got %s", activity.Type)
	}

	if activity.Title != "" {
		t.Errorf("expected empty title, got %s", activity.Title)
	}

	if activity.Content != "" {
		t.Errorf("expected empty content, got %s", activity.Content)
	}

	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(activity.Metadata), &meta); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}

	if meta["club_name"] != club.Name {
		t.Errorf("metadata club_name expected %s, got %v", club.Name, meta["club_name"])
	}
}
