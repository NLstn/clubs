package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

// Notification represents a single user notification
type Notification struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;index" json:"userId"`
	ClubID    string    `gorm:"type:uuid" json:"clubId,omitempty"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateNotification stores a new notification for a user
func CreateNotification(userID, clubID, message string) (Notification, error) {
	n := Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		ClubID:    clubID,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := database.Db.Create(&n).Error
	return n, err
}

// GetNotifications returns all notifications for the user ordered by date desc
func GetNotifications(userID string) ([]Notification, error) {
	var notis []Notification
	err := database.Db.Order("created_at desc").Where("user_id = ?", userID).Find(&notis).Error
	return notis, err
}

// UpdateNotificationRead sets the read status of a notification
func UpdateNotificationRead(id, userID string, read bool) error {
	return database.Db.Model(&Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("read", read).Error
}

// DeleteNotification removes a notification
func DeleteNotification(id, userID string) error {
	return database.Db.Where("id = ? AND user_id = ?", id, userID).Delete(&Notification{}).Error
}
