package models_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/database"
	"github.com/NLstn/civo/handlers"
	"github.com/NLstn/civo/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserPrivacySettingsAuthorization tests authorization for UserPrivacySettings
func TestUserPrivacySettingsAuthorization(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create test users
	user1, _ := handlers.CreateTestUser(t, "user1@example.com")
	user2, _ := handlers.CreateTestUser(t, "user2@example.com")

	// Create privacy settings for user1
	settings1 := models.UserPrivacySettings{
		UserID:         user1.ID,
		ShareBirthDate: true,
	}
	err := database.Db.Create(&settings1).Error
	require.NoError(t, err)

	// Create privacy settings for user2
	settings2 := models.UserPrivacySettings{
		UserID:         user2.ID,
		ShareBirthDate: false,
	}
	err = database.Db.Create(&settings2).Error
	require.NoError(t, err)

	t.Run("user_can_read_own_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("GET", "/api/v2/UserPrivacySettings", nil)

		ups := models.UserPrivacySettings{}
		scopes, err := ups.ODataBeforeReadCollection(ctx, req, nil)
		require.NoError(t, err)
		require.Len(t, scopes, 1)

		// Apply scope and fetch settings
		query := database.Db
		for _, scope := range scopes {
			query = scope(query)
		}

		var results []models.UserPrivacySettings
		err = query.Find(&results).Error
		require.NoError(t, err)

		// Should only see own settings
		assert.Len(t, results, 1)
		assert.Equal(t, user1.ID, results[0].UserID)
		assert.True(t, results[0].ShareBirthDate)
	})

	t.Run("user_cannot_read_other_user_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("GET", "/api/v2/UserPrivacySettings", nil)

		ups := models.UserPrivacySettings{}
		scopes, err := ups.ODataBeforeReadCollection(ctx, req, nil)
		require.NoError(t, err)

		// Apply scope and try to fetch user2's settings
		query := database.Db
		for _, scope := range scopes {
			query = scope(query)
		}

		var results []models.UserPrivacySettings
		err = query.Where("user_id = ?", user2.ID).Find(&results).Error
		require.NoError(t, err)

		// Should not see user2's settings
		assert.Len(t, results, 0)
	})

	t.Run("user_can_create_own_privacy_settings", func(t *testing.T) {
		user3, _ := handlers.CreateTestUser(t, "user3@example.com")
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user3.ID)
		req, _ := http.NewRequest("POST", "/api/v2/UserPrivacySettings", nil)

		newSettings := models.UserPrivacySettings{
			UserID:         user3.ID,
			ShareBirthDate: true,
		}

		err := newSettings.ODataBeforeCreate(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, user3.ID, newSettings.UserID)
	})

	t.Run("user_cannot_create_privacy_settings_for_another_user", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("POST", "/api/v2/UserPrivacySettings", nil)

		newSettings := models.UserPrivacySettings{
			UserID:         user2.ID,
			ShareBirthDate: true,
		}

		err := newSettings.ODataBeforeCreate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("user_can_update_own_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("PATCH", "/api/v2/UserPrivacySettings", nil)

		settings1.ShareBirthDate = false
		err := settings1.ODataBeforeUpdate(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("user_cannot_update_other_user_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("PATCH", "/api/v2/UserPrivacySettings", nil)

		settings2.ShareBirthDate = true
		err := settings2.ODataBeforeUpdate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("user_can_delete_own_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("DELETE", "/api/v2/UserPrivacySettings", nil)

		err := settings1.ODataBeforeDelete(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("user_cannot_delete_other_user_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("DELETE", "/api/v2/UserPrivacySettings", nil)

		err := settings2.ODataBeforeDelete(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

// TestMemberPrivacySettingsAuthorization tests authorization for MemberPrivacySettings
func TestMemberPrivacySettingsAuthorization(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create test users
	user1, _ := handlers.CreateTestUser(t, "user1@example.com")
	user2, _ := handlers.CreateTestUser(t, "user2@example.com")

	// Create test clubs and members
	club1 := handlers.CreateTestClub(t, user1, "Club 1")
	club2 := handlers.CreateTestClub(t, user2, "Club 2")

	// CreateTestClub already creates owner members, so just retrieve them
	var member1 models.Member
	err := database.Db.Where("user_id = ? AND club_id = ?", user1.ID, club1.ID).First(&member1).Error
	require.NoError(t, err)

	var member2 models.Member
	err = database.Db.Where("user_id = ? AND club_id = ?", user2.ID, club2.ID).First(&member2).Error
	require.NoError(t, err)

	// Create additional member for user1 in club2
	member1InClub2 := handlers.CreateTestMember(t, user1, club2, "member")

	// Create privacy settings for members
	memberSettings1 := models.MemberPrivacySettings{
		MemberID:       member1.ID,
		ShareBirthDate: true,
	}
	err = database.Db.Create(&memberSettings1).Error
	require.NoError(t, err)

	memberSettings2 := models.MemberPrivacySettings{
		MemberID:       member2.ID,
		ShareBirthDate: false,
	}
	err = database.Db.Create(&memberSettings2).Error
	require.NoError(t, err)

	t.Run("user_can_read_own_member_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("GET", "/api/v2/MemberPrivacySettings", nil)

		mps := models.MemberPrivacySettings{}
		scopes, err := mps.ODataBeforeReadCollection(ctx, req, nil)
		require.NoError(t, err)
		require.Len(t, scopes, 1)

		// Apply scope and fetch settings
		query := database.Db
		for _, scope := range scopes {
			query = scope(query)
		}

		var results []models.MemberPrivacySettings
		err = query.Find(&results).Error
		require.NoError(t, err)

		// Should only see own member privacy settings
		assert.Len(t, results, 1)
		assert.Equal(t, member1.ID, results[0].MemberID)
		assert.True(t, results[0].ShareBirthDate)
	})

	t.Run("user_cannot_read_other_user_member_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("GET", "/api/v2/MemberPrivacySettings", nil)

		mps := models.MemberPrivacySettings{}
		scopes, err := mps.ODataBeforeReadCollection(ctx, req, nil)
		require.NoError(t, err)

		// Apply scope and try to fetch user2's settings
		query := database.Db
		for _, scope := range scopes {
			query = scope(query)
		}

		var results []models.MemberPrivacySettings
		err = query.Where("member_id = ?", member2.ID).Find(&results).Error
		require.NoError(t, err)

		// Should not see user2's settings
		assert.Len(t, results, 0)
	})

	t.Run("user_can_create_privacy_settings_for_own_member", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("POST", "/api/v2/MemberPrivacySettings", nil)

		newSettings := models.MemberPrivacySettings{
			MemberID:       member1InClub2.ID,
			ShareBirthDate: true,
		}

		err := newSettings.ODataBeforeCreate(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("user_cannot_create_privacy_settings_for_another_user_member", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("POST", "/api/v2/MemberPrivacySettings", nil)

		newSettings := models.MemberPrivacySettings{
			MemberID:       member2.ID,
			ShareBirthDate: true,
		}

		err := newSettings.ODataBeforeCreate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("user_can_update_own_member_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("PATCH", "/api/v2/MemberPrivacySettings", nil)

		memberSettings1.ShareBirthDate = false
		err := memberSettings1.ODataBeforeUpdate(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("user_cannot_update_other_user_member_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("PATCH", "/api/v2/MemberPrivacySettings", nil)

		memberSettings2.ShareBirthDate = true
		err := memberSettings2.ODataBeforeUpdate(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("user_can_delete_own_member_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("DELETE", "/api/v2/MemberPrivacySettings", nil)

		err := memberSettings1.ODataBeforeDelete(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("user_cannot_delete_other_user_member_privacy_settings", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user1.ID)
		req, _ := http.NewRequest("DELETE", "/api/v2/MemberPrivacySettings", nil)

		err := memberSettings2.ODataBeforeDelete(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

// TestGetEffectivePrivacySettings tests the helper function for getting effective privacy settings
func TestGetEffectivePrivacySettings(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	user1, _ := handlers.CreateTestUser(t, "user1@example.com")
	club1 := handlers.CreateTestClub(t, user1, "Club 1")

	// CreateTestClub already creates a member for the owner, retrieve it
	var member1 models.Member
	err := database.Db.Where("user_id = ? AND club_id = ?", user1.ID, club1.ID).First(&member1).Error
	require.NoError(t, err)

	// Create global privacy settings
	globalSettings := models.UserPrivacySettings{
		UserID:         user1.ID,
		ShareBirthDate: true,
	}
	err = database.Db.Create(&globalSettings).Error
	require.NoError(t, err)

	t.Run("returns_global_setting_when_no_member_override", func(t *testing.T) {
		shareBirthDate, err := models.GetEffectivePrivacySettings(user1.ID, club1.ID)
		require.NoError(t, err)
		assert.True(t, shareBirthDate)
	})

	t.Run("returns_member_override_when_exists", func(t *testing.T) {
		// Create member-specific override
		memberSettings := models.MemberPrivacySettings{
			ID:             "test-member-privacy-settings-id",
			MemberID:       member1.ID,
			ShareBirthDate: false, // Different from global
		}
		err := database.Db.Create(&memberSettings).Error
		require.NoError(t, err)

		// Verify it was created
		var check models.MemberPrivacySettings
		err = database.Db.Where("member_id = ?", member1.ID).First(&check).Error
		require.NoError(t, err, "Should find the member privacy settings we just created")

		shareBirthDate, err := models.GetEffectivePrivacySettings(user1.ID, club1.ID)
		require.NoError(t, err)
		assert.False(t, shareBirthDate) // Should use member override
	})
}
