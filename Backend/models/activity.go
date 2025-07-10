package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID        string    `json:"id" gorm:"primaryKey;type:char(36)"`
	ClubID    string    `json:"club_id" gorm:"type:char(36);not null;index"`
	UserID    string    `json:"user_id" gorm:"type:char(36);not null;index"` // User who performed the action or was affected
	ActorID   *string   `json:"actor_id" gorm:"type:char(36);index"`         // User who initiated the action (e.g., admin who promoted someone)
	Type      string    `json:"type" gorm:"type:varchar(50);not null;index"` // "role_changed", "member_promoted", etc.
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Metadata  string    `json:"metadata" gorm:"type:json"` // JSON field for additional data
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Club  Club  `json:"club,omitempty" gorm:"foreignKey:ClubID"`
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Actor *User `json:"actor,omitempty" gorm:"foreignKey:ActorID"`
}

func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// CreateActivity creates a new activity record
func CreateActivity(clubID, userID string, actorID *string, activityType, title, content string, metadata map[string]interface{}) error {
	activity := Activity{
		ClubID:  clubID,
		UserID:  userID,
		ActorID: actorID,
		Type:    activityType,
		Title:   title,
		Content: content,
	}

	// Convert metadata to JSON string if provided
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			activity.Metadata = string(metadataBytes)
		}
	}

	return database.Db.Create(&activity).Error
}

// GetClubActivities retrieves activities for a specific club
func GetClubActivities(clubID string, limit int) ([]Activity, error) {
	var activities []Activity
	query := database.Db.Where("club_id = ?", clubID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&activities).Error
	return activities, err
}

// GetRecentActivities retrieves recent activities for multiple clubs
func GetRecentActivities(clubIDs []string, daysBack int, limit int) ([]Activity, error) {
	var activities []Activity
	query := database.Db.Where("club_id IN ? AND created_at > ?", clubIDs, time.Now().AddDate(0, 0, -daysBack)).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&activities).Error
	return activities, err
}

// CreateRoleChangeActivity creates an activity when a user's role changes
func CreateRoleChangeActivity(clubID, userID, actorID, clubName, oldRole, newRole string) error {
	title := "Role updated"
	content := ""
	activityType := "role_changed"

	// Determine if this is a promotion or demotion
	roleHierarchy := map[string]int{"member": 1, "admin": 2, "owner": 3}
	oldLevel := roleHierarchy[oldRole]
	newLevel := roleHierarchy[newRole]

	if newLevel > oldLevel {
		title = "Member promoted"
		content = fmt.Sprintf("Role changed from %s to %s", oldRole, newRole)
		activityType = "member_promoted"
	} else if newLevel < oldLevel {
		title = "Member demoted"
		content = fmt.Sprintf("Role changed from %s to %s", oldRole, newRole)
		activityType = "member_demoted"
	} else {
		content = fmt.Sprintf("Role changed from %s to %s", oldRole, newRole)
	}

	metadata := map[string]interface{}{
		"old_role":  oldRole,
		"new_role":  newRole,
		"club_name": clubName,
	}

	var actor *string
	if actorID != "" {
		actor = &actorID
	}

	return CreateActivity(clubID, userID, actor, activityType, title, content, metadata)
}
