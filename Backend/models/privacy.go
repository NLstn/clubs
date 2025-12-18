package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"gorm.io/gorm"
)

type UserPrivacySettings struct {
	ID             string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	UserID         string    `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
	ClubID         *string   `json:"ClubID,omitempty" gorm:"type:uuid" odata:"nullable"` // NULL means global setting
	ShareBirthDate bool      `json:"ShareBirthDate" gorm:"default:false"`
	CreatedAt      time.Time `json:"CreatedAt" odata:"immutable"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}

// EntitySetName returns the custom entity set name to prevent double pluralization
// Without this, the OData library would pluralize "UserPrivacySettings" to "UserPrivacySettingses"
func (UserPrivacySettings) EntitySetName() string {
	return "UserPrivacySettings"
}

// GetUserPrivacySettings returns privacy settings for a user and specific club
// If no club-specific setting exists, returns the global setting
// If no settings exist at all, returns default (private) settings
func GetUserPrivacySettings(userID, clubID string) (*UserPrivacySettings, error) {
	var settings UserPrivacySettings

	// First try to get club-specific settings
	err := database.Db.Where("user_id = ? AND club_id = ?", userID, clubID).First(&settings).Error
	if err == nil {
		return &settings, nil
	}

	// If no club-specific settings, try global settings
	err = database.Db.Where("user_id = ? AND club_id IS NULL", userID).First(&settings).Error
	if err == nil {
		return &settings, nil
	}

	// If no settings exist, return default private settings
	return &UserPrivacySettings{
		UserID:         userID,
		ClubID:         nil,
		ShareBirthDate: false,
	}, nil
}

// GetUserGlobalPrivacySettings returns global privacy settings for a user
func GetUserGlobalPrivacySettings(userID string) (*UserPrivacySettings, error) {
	var settings UserPrivacySettings
	err := database.Db.Where("user_id = ? AND club_id IS NULL", userID).First(&settings).Error
	if err != nil {
		// Return default settings if none exist
		return &UserPrivacySettings{
			UserID:         userID,
			ClubID:         nil,
			ShareBirthDate: false,
		}, nil
	}
	return &settings, nil
}

// GetUserClubSpecificPrivacySettings returns all club-specific privacy settings for a user
func GetUserClubSpecificPrivacySettings(userID string) ([]UserPrivacySettings, error) {
	var settings []UserPrivacySettings
	err := database.Db.Where("user_id = ? AND club_id IS NOT NULL", userID).Find(&settings).Error
	return settings, err
}

// UpdateOrCreatePrivacySettings updates or creates privacy settings
func UpdateOrCreatePrivacySettings(userID, clubID string, shareBirthDate bool) error {
	var settings UserPrivacySettings

	var err error
	if clubID == "" {
		// Global setting
		err = database.Db.Where("user_id = ? AND club_id IS NULL", userID).First(&settings).Error
	} else {
		// Club-specific setting
		err = database.Db.Where("user_id = ? AND club_id = ?", userID, clubID).First(&settings).Error
	}

	if err != nil {
		// Create new settings
		settings = UserPrivacySettings{
			UserID:         userID,
			ShareBirthDate: shareBirthDate,
		}
		if clubID != "" {
			settings.ClubID = &clubID
		}
		// For global settings (clubID == ""), ClubID should remain empty/null
		return database.Db.Create(&settings).Error
	} else {
		// Update existing settings
		settings.ShareBirthDate = shareBirthDate
		return database.Db.Save(&settings).Error
	}
}

// ODataBeforeReadCollection filters privacy settings to only those belonging to the user
func (ups UserPrivacySettings) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see their own privacy settings
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to specific privacy settings
func (ups UserPrivacySettings) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see their own privacy settings
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates privacy settings creation
func (ups *UserPrivacySettings) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Users can only create their own privacy settings
	if ups.UserID == "" {
		ups.UserID = userID
	} else if ups.UserID != userID {
		return fmt.Errorf("unauthorized: cannot create privacy settings for another user")
	}

	// Set CreatedAt and UpdatedAt
	now := time.Now()
	ups.CreatedAt = now
	ups.UpdatedAt = now

	return nil
}

// ODataBeforeUpdate validates privacy settings update permissions
func (ups *UserPrivacySettings) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Users can only update their own privacy settings
	if ups.UserID != userID {
		return fmt.Errorf("unauthorized: can only update your own privacy settings")
	}

	// Set UpdatedAt
	ups.UpdatedAt = time.Now()

	return nil
}

// ODataBeforeDelete validates privacy settings deletion permissions
func (ups *UserPrivacySettings) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Users can only delete their own privacy settings
	if ups.UserID != userID {
		return fmt.Errorf("unauthorized: can only delete your own privacy settings")
	}

	return nil
}
