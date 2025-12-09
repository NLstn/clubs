package models

import (
	"time"

	"github.com/NLstn/clubs/database"
)

type UserPrivacySettings struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	UserID         string    `gorm:"type:uuid;not null" odata:"required"`
	ClubID         *string   `gorm:"type:uuid" odata:"nullable"` // NULL means global setting
	ShareBirthDate bool      `gorm:"default:false"`
	CreatedAt      time.Time `odata:"immutable"`
	UpdatedAt      time.Time
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
