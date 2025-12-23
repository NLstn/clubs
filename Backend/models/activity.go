package models

import (
	"encoding/json"
	"time"

	"github.com/NLstn/civo/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID        string    `json:"ID" gorm:"primaryKey;type:char(36)"`
	ClubID    string    `json:"ClubID" gorm:"type:char(36);not null;index"`
	UserID    string    `json:"UserID" gorm:"type:char(36);not null;index"`   // User who performed the action or was affected
	ActorID   *string   `json:"ActorID,omitempty" gorm:"type:char(36);index"` // User who initiated the action (e.g., admin who promoted someone)
	Type      string    `json:"Type" gorm:"type:varchar(50);not null;index"`  // "role_changed", "member_promoted", etc.
	Title     string    `json:"Title" gorm:"type:varchar(255);not null"`
	Content   string    `json:"Content" gorm:"type:text"`
	Metadata  string    `json:"Metadata" gorm:"type:json"` // JSON field for additional data
	CreatedAt time.Time `json:"CreatedAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"autoUpdateTime"`

	// Relations
	Club  Club  `json:"Club,omitempty" gorm:"foreignKey:ClubID"`
	User  User  `json:"User,omitempty" gorm:"foreignKey:UserID"`
	Actor *User `json:"Actor,omitempty" gorm:"foreignKey:ActorID"`
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
	activityType := "role_changed"

	// Determine if this is a promotion or demotion
	roleHierarchy := map[string]int{"member": 1, "admin": 2, "owner": 3}
	oldLevel := roleHierarchy[oldRole]
	newLevel := roleHierarchy[newRole]

	if newLevel > oldLevel {
		activityType = "member_promoted"
	} else if newLevel < oldLevel {
		activityType = "member_demoted"
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

	return CreateActivity(clubID, userID, actor, activityType, "", "", metadata)
}

// CreateMemberJoinedActivity creates an activity when a new member joins the club
func CreateMemberJoinedActivity(clubID, userID, clubName string, actorID *string) error {
	metadata := map[string]interface{}{
		"club_name": clubName,
	}

	return CreateActivity(clubID, userID, actorID, "member_joined", "", "", metadata)
}
