package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/database"
	"gorm.io/gorm"
)

// UserPrivacySettings stores global privacy settings for a user
// These settings apply to all clubs unless overridden by MemberPrivacySettings
type UserPrivacySettings struct {
	ID             string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	UserID         string    `json:"UserID" gorm:"type:uuid;not null;uniqueIndex" odata:"required"`
	ShareBirthDate bool      `json:"ShareBirthDate" gorm:"default:false"`
	CreatedAt      time.Time `json:"CreatedAt" odata:"immutable"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}

// EntitySetName returns the custom entity set name for the UserPrivacySettings entity.
// By default, "UserPrivacySettings" would be pluralized to "UserPrivacySettingses", so we override it.
func (UserPrivacySettings) EntitySetName() string {
	return "UserPrivacySettings"
}

// MemberPrivacySettings stores club-specific privacy settings (overrides for global settings)
// When present, these settings take precedence over UserPrivacySettings for the specific club
type MemberPrivacySettings struct {
	ID             string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	MemberID       string    `json:"MemberID" gorm:"type:uuid;not null;uniqueIndex" odata:"required"`
	ShareBirthDate bool      `json:"ShareBirthDate" gorm:"default:false"`
	CreatedAt      time.Time `json:"CreatedAt" odata:"immutable"`
	UpdatedAt      time.Time `json:"UpdatedAt"`

	// Navigation property for OData
	Member *Member `gorm:"foreignKey:MemberID" json:"Member,omitempty" odata:"nav"`
}

// EntitySetName returns the custom entity set name for the MemberPrivacySettings entity.
// By default, "MemberPrivacySettings" would be pluralized to "MemberPrivacySettingses", so we override it.
func (MemberPrivacySettings) EntitySetName() string {
	return "MemberPrivacySettings"
}

// GetUserGlobalPrivacySettings returns global privacy settings for a user
func GetUserGlobalPrivacySettings(userID string) (*UserPrivacySettings, error) {
	var settings UserPrivacySettings
	err := database.Db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return default settings if none exist
			return &UserPrivacySettings{
				UserID:         userID,
				ShareBirthDate: false,
			}, nil
		}
		// Return other database errors
		return nil, err
	}
	return &settings, nil
}

// GetMemberPrivacySettings returns privacy settings for a specific member
// If no member-specific setting exists, returns nil
func GetMemberPrivacySettings(memberID string) (*MemberPrivacySettings, error) {
	var settings MemberPrivacySettings
	err := database.Db.Where("member_id = ?", memberID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

// GetEffectivePrivacySettings returns the effective privacy settings for a user in a club
// If member-specific settings exist, they take precedence over global settings
func GetEffectivePrivacySettings(userID, clubID string) (bool, error) {
	// First get the member record
	var member Member
	err := database.Db.Where("user_id = ? AND club_id = ?", userID, clubID).First(&member).Error
	if err != nil {
		return false, err
	}

	// Try to get member-specific settings
	memberSettings, err := GetMemberPrivacySettings(member.ID)
	if err != nil {
		return false, err
	}
	if memberSettings != nil {
		return memberSettings.ShareBirthDate, nil
	}

	// Fall back to global settings
	globalSettings, err := GetUserGlobalPrivacySettings(userID)
	if err != nil {
		return false, err
	}
	return globalSettings.ShareBirthDate, nil
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

// MemberPrivacySettings OData Hooks

// ODataBeforeReadCollection filters member privacy settings to those belonging to user's members
func (mps MemberPrivacySettings) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see privacy settings for their own member records
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("member_id IN (SELECT id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to specific member privacy settings
func (mps MemberPrivacySettings) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see privacy settings for their own member records
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("member_id IN (SELECT id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates member privacy settings creation
func (mps *MemberPrivacySettings) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Verify the member belongs to the current user
	var member Member
	err := database.Db.Where("id = ? AND user_id = ?", mps.MemberID, userID).First(&member).Error
	if err != nil {
		return fmt.Errorf("unauthorized: cannot create privacy settings for another user's member record")
	}

	// Set CreatedAt and UpdatedAt
	now := time.Now()
	mps.CreatedAt = now
	mps.UpdatedAt = now

	return nil
}

// ODataBeforeUpdate validates member privacy settings update permissions
func (mps *MemberPrivacySettings) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Verify the member belongs to the current user
	var member Member
	err := database.Db.Where("id = ? AND user_id = ?", mps.MemberID, userID).First(&member).Error
	if err != nil {
		return fmt.Errorf("unauthorized: can only update privacy settings for your own member records")
	}

	// Set UpdatedAt
	mps.UpdatedAt = time.Now()

	return nil
}

// ODataBeforeDelete validates member privacy settings deletion permissions
func (mps *MemberPrivacySettings) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Verify the member belongs to the current user
	var member Member
	err := database.Db.Where("id = ? AND user_id = ?", mps.MemberID, userID).First(&member).Error
	if err != nil {
		return fmt.Errorf("unauthorized: can only delete privacy settings for your own member records")
	}

	return nil
}
