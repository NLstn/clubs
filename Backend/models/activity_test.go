package models_test

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
