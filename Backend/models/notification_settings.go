package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

// NotificationSetting stores user preference for notification channels
// If ClubID is empty, the setting is global
// Channels supported: Email and InApp

type NotificationSetting struct {
	ID     string `gorm:"type:uuid;primaryKey" json:"id"`
	UserID string `gorm:"type:uuid;index" json:"userId"`
	ClubID string `gorm:"type:uuid" json:"clubId,omitempty"`
	Email  bool   `json:"email"`
	InApp  bool   `json:"inApp"`
}

// GetUserNotificationSettings returns all settings for a user
func GetUserNotificationSettings(userID string) ([]NotificationSetting, error) {
	var settings []NotificationSetting
	err := database.Db.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

// UpsertNotificationSetting inserts or updates a setting
func UpsertNotificationSetting(s *NotificationSetting) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
		return database.Db.Create(s).Error
	}
	return database.Db.Save(s).Error
}

// DeleteNotificationSetting removes a setting
func DeleteNotificationSetting(id, userID string) error {
	return database.Db.Where("id = ? AND user_id = ?", id, userID).Delete(&NotificationSetting{}).Error
}
