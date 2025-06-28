package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents an in-app notification for a user
type Notification struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	UserID    string    `json:"userId" gorm:"type:uuid;not null"`
	Type      string    `json:"type" gorm:"not null"` // e.g., "member_added", "event_created", "fine_assigned"
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Read      bool      `json:"read" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
	// Optional data for linking to specific resources
	ClubID  *string `json:"clubId,omitempty" gorm:"type:uuid"`
	EventID *string `json:"eventId,omitempty" gorm:"type:uuid"`
	FineID  *string `json:"fineId,omitempty" gorm:"type:uuid"`
}

// UserNotificationPreferences represents user's notification settings
type UserNotificationPreferences struct {
	ID                    string    `json:"id" gorm:"type:uuid;primary_key"`
	UserID                string    `json:"userId" gorm:"type:uuid;not null;unique"`
	MemberAddedInApp      bool      `json:"memberAddedInApp" gorm:"default:true"`
	MemberAddedEmail      bool      `json:"memberAddedEmail" gorm:"default:true"`
	EventCreatedInApp     bool      `json:"eventCreatedInApp" gorm:"default:true"`
	EventCreatedEmail     bool      `json:"eventCreatedEmail" gorm:"default:false"`
	FineAssignedInApp     bool      `json:"fineAssignedInApp" gorm:"default:true"`
	FineAssignedEmail     bool      `json:"fineAssignedEmail" gorm:"default:true"`
	NewsCreatedInApp      bool      `json:"newsCreatedInApp" gorm:"default:true"`
	NewsCreatedEmail      bool      `json:"newsCreatedEmail" gorm:"default:false"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// BeforeCreate sets the ID for new notifications
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate sets the ID for new user notification preferences
func (p *UserNotificationPreferences) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// GetUserNotifications retrieves notifications for a user
func GetUserNotifications(userID string, limit int) ([]Notification, error) {
	var notifications []Notification
	query := database.Db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&notifications).Error
	return notifications, err
}

// GetUnreadNotificationCount returns the count of unread notifications for a user
func GetUnreadNotificationCount(userID string) (int64, error) {
	var count int64
	err := database.Db.Model(&Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&count).Error
	return count, err
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(notificationID, userID string) error {
	return database.Db.Model(&Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("read", true).Error
}

// MarkAllNotificationsAsRead marks all notifications for a user as read
func MarkAllNotificationsAsRead(userID string) error {
	return database.Db.Model(&Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error
}

// CreateNotification creates a new notification
func CreateNotification(userID, notificationType, title, message string, clubID, eventID, fineID *string) error {
	notification := Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		ClubID:  clubID,
		EventID: eventID,
		FineID:  fineID,
	}
	return database.Db.Create(&notification).Error
}

// GetUserNotificationPreferences retrieves user notification preferences
func GetUserNotificationPreferences(userID string) (UserNotificationPreferences, error) {
	var preferences UserNotificationPreferences
	err := database.Db.Where("user_id = ?", userID).First(&preferences).Error
	if err != nil {
		if err.Error() == "record not found" {
			// Create default preferences if none exist
			return CreateDefaultUserNotificationPreferences(userID)
		}
		return UserNotificationPreferences{}, err
	}
	return preferences, nil
}

// CreateDefaultUserNotificationPreferences creates default notification preferences for a user
func CreateDefaultUserNotificationPreferences(userID string) (UserNotificationPreferences, error) {
	preferences := UserNotificationPreferences{
		UserID:                userID,
		MemberAddedInApp:      true,
		MemberAddedEmail:      true,
		EventCreatedInApp:     true,
		EventCreatedEmail:     false,
		FineAssignedInApp:     true,
		FineAssignedEmail:     true,
		NewsCreatedInApp:      true,
		NewsCreatedEmail:      false,
	}
	err := database.Db.Create(&preferences).Error
	return preferences, err
}

// UpdateUserNotificationPreferences updates user notification preferences
func (p *UserNotificationPreferences) Update() error {
	return database.Db.Save(p).Error
}

// SendMemberAddedNotifications handles both in-app and email notifications for member addition
func SendMemberAddedNotifications(userID, userEmail, clubID, clubName string) error {
	// Get user notification preferences
	preferences, err := GetUserNotificationPreferences(userID)
	if err != nil {
		// If preferences don't exist, create default ones and continue
		preferences, err = CreateDefaultUserNotificationPreferences(userID)
		if err != nil {
			return fmt.Errorf("failed to create notification preferences: %v", err)
		}
	}

	// Send in-app notification if enabled
	if preferences.MemberAddedInApp {
		title := "Welcome to " + clubName
		message := fmt.Sprintf("You have been added to the club %s as a member.", clubName)
		err := CreateNotification(userID, "member_added", title, message, &clubID, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to create in-app notification: %v", err)
		}
	}

	// Import the notifications package for email sending
	// We'll use a simplified approach here to avoid circular imports
	// The email sending will be handled by the caller
	
	return nil
}