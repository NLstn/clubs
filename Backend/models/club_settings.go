package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClubSettings struct {
	ID                       string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID                   string    `json:"ClubID" gorm:"type:uuid;not null;unique" odata:"required"`
	FinesEnabled             bool      `json:"FinesEnabled" gorm:"default:false"`
	ShiftsEnabled            bool      `json:"ShiftsEnabled" gorm:"default:false"`
	TeamsEnabled             bool      `json:"TeamsEnabled" gorm:"default:false"`
	MembersListVisible       bool      `json:"MembersListVisible" gorm:"default:false"`
	DiscoverableByNonMembers bool      `json:"DiscoverableByNonMembers" gorm:"default:false"`
	CreatedAt                time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy                string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt                time.Time `json:"UpdatedAt"`
	UpdatedBy                string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`

	// Navigation properties for OData
	Club *Club `gorm:"foreignKey:ClubID" json:"Club,omitempty" odata:"nav"`
}

// EntitySetName returns the custom entity set name to prevent double pluralization
// Without this, the OData library would pluralize "ClubSettings" to "ClubSettingses"
func (ClubSettings) EntitySetName() string {
	return "ClubSettings"
}

func GetClubSettings(clubID string) (ClubSettings, error) {
	var settings ClubSettings
	result := database.Db.First(&settings, "club_id = ?", clubID)
	if result.Error == gorm.ErrRecordNotFound {
		// Create default settings if none exist
		return CreateDefaultClubSettings(clubID)
	}
	return settings, result.Error
}

func CreateDefaultClubSettings(clubID string) (ClubSettings, error) {
	settings := ClubSettings{
		ID:                       uuid.New().String(),
		ClubID:                   clubID,
		FinesEnabled:             false,
		ShiftsEnabled:            false,
		TeamsEnabled:             false,
		MembersListVisible:       false,
		DiscoverableByNonMembers: false,
		CreatedBy:                clubID, // Using clubID as default since we don't have user context here
		UpdatedBy:                clubID,
	}
	err := database.Db.Create(&settings).Error
	return settings, err
}

func (s *ClubSettings) Update(finesEnabled, shiftsEnabled, teamsEnabled, membersListVisible, discoverableByNonMembers bool, updatedBy string) error {
	return database.Db.Model(s).Updates(map[string]interface{}{
		"fines_enabled":               finesEnabled,
		"shifts_enabled":              shiftsEnabled,
		"teams_enabled":               teamsEnabled,
		"members_list_visible":        membersListVisible,
		"discoverable_by_non_members": discoverableByNonMembers,
		"updated_by":                  updatedBy,
	}).Error
}

// ODataBeforeReadCollection OData hook - restricts ClubSettings access to club members only
func (s ClubSettings) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	scope := func(db *gorm.DB) *gorm.DB {
		// Only show settings for clubs where user is a member
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity OData hook - restricts ClubSettings access to club members only
func (s ClubSettings) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	scope := func(db *gorm.DB) *gorm.DB {
		// Only allow access to settings for clubs where user is a member
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeUpdate OData hook - restricts ClubSettings updates to club admins only
func (s *ClubSettings) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get the club to check if user is admin
	var club Club
	if err := database.Db.First(&club, "id = ?", s.ClubID).Error; err != nil {
		return fmt.Errorf("failed to find club: %w", err)
	}

	// Get user
	var user User
	if err := database.Db.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Only admins and owners can update settings
	if !club.IsAdmin(user) {
		return fmt.Errorf("forbidden: only club admins can update settings")
	}

	// Set UpdatedBy and UpdatedAt
	now := time.Now()
	s.UpdatedAt = now
	s.UpdatedBy = userID

	return nil
}
