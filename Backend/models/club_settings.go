package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClubSettings struct {
	ID                        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID                    string    `json:"ClubID" gorm:"type:uuid;not null;unique" odata:"required"`
	FinesEnabled              bool      `json:"FinesEnabled" gorm:"default:true"`
	ShiftsEnabled             bool      `json:"ShiftsEnabled" gorm:"default:true"`
	TeamsEnabled              bool      `json:"TeamsEnabled" gorm:"default:true"`
	MembersListVisible        bool      `json:"MembersListVisible" gorm:"default:true"`
	DiscoverableByNonMembers  bool      `json:"DiscoverableByNonMembers" gorm:"default:false"`
	CreatedAt                 time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy                 string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt                 time.Time `json:"UpdatedAt"`
	UpdatedBy                 string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
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
		FinesEnabled:             true,
		ShiftsEnabled:            true,
		TeamsEnabled:             true,
		MembersListVisible:       true,
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
