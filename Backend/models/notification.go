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
	Type      string    `json:"type" gorm:"not null"` // e.g., "member_added", "event_created", "fine_assigned", "invite_received"
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Read      bool      `json:"read" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
	// Optional data for linking to specific resources
	ClubID   *string `json:"clubId,omitempty" gorm:"type:uuid"`
	EventID  *string `json:"eventId,omitempty" gorm:"type:uuid"`
	FineID   *string `json:"fineId,omitempty" gorm:"type:uuid"`
	InviteID *string `json:"inviteId,omitempty" gorm:"type:uuid"`
}

// UserNotificationPreferences represents user's notification settings
type UserNotificationPreferences struct {
	ID                  string    `json:"id" gorm:"type:uuid;primary_key"`
	UserID              string    `json:"userId" gorm:"type:uuid;not null;unique"`
	MemberAddedInApp    bool      `json:"memberAddedInApp" gorm:"default:true"`
	MemberAddedEmail    bool      `json:"memberAddedEmail" gorm:"default:true"`
	InviteReceivedInApp bool      `json:"inviteReceivedInApp" gorm:"default:true"`
	InviteReceivedEmail bool      `json:"inviteReceivedEmail" gorm:"default:true"`
	EventCreatedInApp   bool      `json:"eventCreatedInApp" gorm:"default:true"`
	EventCreatedEmail   bool      `json:"eventCreatedEmail" gorm:"default:false"`
	FineAssignedInApp   bool      `json:"fineAssignedInApp" gorm:"default:true"`
	FineAssignedEmail   bool      `json:"fineAssignedEmail" gorm:"default:true"`
	NewsCreatedInApp    bool      `json:"newsCreatedInApp" gorm:"default:true"`
	NewsCreatedEmail    bool      `json:"newsCreatedEmail" gorm:"default:false"`
	RoleChangedInApp    bool      `json:"roleChangedInApp" gorm:"default:true"`
	RoleChangedEmail    bool      `json:"roleChangedEmail" gorm:"default:true"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
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

// DeleteNotification deletes a notification for a user
func DeleteNotification(notificationID, userID string) error {
	return database.Db.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&Notification{}).Error
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

// CreateNotificationWithInvite creates a new notification with invite reference
func CreateNotificationWithInvite(userID, notificationType, title, message string, clubID, eventID, fineID, inviteID *string) error {
	notification := Notification{
		UserID:   userID,
		Type:     notificationType,
		Title:    title,
		Message:  message,
		ClubID:   clubID,
		EventID:  eventID,
		FineID:   fineID,
		InviteID: inviteID,
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
		UserID:              userID,
		MemberAddedInApp:    true,
		MemberAddedEmail:    true,
		InviteReceivedInApp: true,
		InviteReceivedEmail: true,
		EventCreatedInApp:   true,
		EventCreatedEmail:   false,
		FineAssignedInApp:   true,
		FineAssignedEmail:   true,
		NewsCreatedInApp:    true,
		NewsCreatedEmail:    false,
		RoleChangedInApp:    true,
		RoleChangedEmail:    true,
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

// SendInviteReceivedNotifications handles both in-app and email notifications for invite received
func SendInviteReceivedNotifications(userEmail, clubID, clubName, inviteID string) error {
	// Find user by email to get their preferences
	var user User
	err := database.Db.Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		// User might not be registered yet, so we'll skip notification preferences
		// but still create a basic notification if they register later
		return nil
	}

	// Get user notification preferences
	preferences, err := GetUserNotificationPreferences(user.ID)
	if err != nil {
		// If preferences don't exist, create default ones and continue
		preferences, err = CreateDefaultUserNotificationPreferences(user.ID)
		if err != nil {
			return fmt.Errorf("failed to create notification preferences: %v", err)
		}
	}

	// Send in-app notification if enabled
	if preferences.InviteReceivedInApp {
		title := "Invitation to " + clubName
		message := fmt.Sprintf("You have been invited to join the club %s.", clubName)
		err := CreateNotificationWithInvite(user.ID, "invite_received", title, message, &clubID, nil, nil, &inviteID)
		if err != nil {
			return fmt.Errorf("failed to create in-app notification: %v", err)
		}
	}

	// Email sending will be handled by the caller if needed

	return nil
}

// RemoveInviteNotifications removes invite notifications when an invite is accepted or rejected
func RemoveInviteNotifications(inviteID string) error {
	return database.Db.Where("invite_id = ? AND type = ?", inviteID, "invite_received").Delete(&Notification{}).Error
}

// SendRoleChangedNotifications handles both in-app and email notifications for role changes
func SendRoleChangedNotifications(userID, clubID, clubName, oldRole, newRole string) error {
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
	if preferences.RoleChangedInApp {
		title := "Role Updated in " + clubName
		message := fmt.Sprintf("Your role in %s has been changed from %s to %s.", clubName, oldRole, newRole)
		err := CreateNotification(userID, "role_changed", title, message, &clubID, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to create in-app notification: %v", err)
		}
	}

	// Send email notification if enabled
	if preferences.RoleChangedEmail {
		var user User
		if err := database.Db.Where("id = ?", userID).First(&user).Error; err != nil {
			return fmt.Errorf("failed to find user for email notification: %v", err)
		}

		// Import the notifications package for email sending
		// Using a simplified approach to avoid circular imports
		// The email sending will be handled by the caller using the notification package functions
	}

	return nil
}
