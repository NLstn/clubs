package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClubSettings struct {
	ID            string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID        string    `json:"clubId" gorm:"type:uuid;not null;unique"`
	FinesEnabled  bool      `json:"finesEnabled" gorm:"default:true"`
	ShiftsEnabled bool      `json:"shiftsEnabled" gorm:"default:true"`
	TeamsEnabled  bool      `json:"teamsEnabled" gorm:"default:true"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy" gorm:"type:uuid"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UpdatedBy     string    `json:"updatedBy" gorm:"type:uuid"`
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
		ID:            uuid.New().String(),
		ClubID:        clubID,
		FinesEnabled:  true,
		ShiftsEnabled: true,
		TeamsEnabled:  true,
		CreatedBy:     clubID, // Using clubID as default since we don't have user context here
		UpdatedBy:     clubID,
	}
	err := database.Db.Create(&settings).Error
	return settings, err
}

func (s *ClubSettings) Update(finesEnabled, shiftsEnabled, teamsEnabled bool, updatedBy string) error {
	return database.Db.Model(s).Updates(map[string]interface{}{
		"fines_enabled":  finesEnabled,
		"shifts_enabled": shiftsEnabled,
		"teams_enabled":  teamsEnabled,
		"updated_by":     updatedBy,
	}).Error
}
