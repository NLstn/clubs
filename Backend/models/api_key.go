package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"gorm.io/gorm"
)

// APIKey represents a long-lived API key for programmatic access
type APIKey struct {
	ID          string     `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	UserID      string     `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
	User        User       `json:"User,omitempty" gorm:"foreignKey:UserID" odata:"navigationProperty"`
	Name        string     `json:"Name" gorm:"not null" odata:"required"`
	KeyHash     string     `json:"-" gorm:"uniqueIndex;not null"` // Never exposed via API
	KeyPrefix   string     `json:"KeyPrefix" gorm:"not null" odata:"immutable"`
	Permissions string     `json:"-" gorm:"type:text"` // Stored as JSON string
	LastUsedAt  *time.Time `json:"LastUsedAt,omitempty" gorm:"type:timestamp" odata:"nullable"`
	ExpiresAt   *time.Time `json:"ExpiresAt,omitempty" gorm:"type:timestamp" odata:"nullable"`
	IsActive    bool       `json:"IsActive" gorm:"default:true" odata:"required"`
	CreatedAt   time.Time  `json:"CreatedAt" odata:"immutable"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
}

// TableName specifies the table name for the APIKey model
func (APIKey) TableName() string {
	return "api_keys"
}

// BeforeCreate hook to ensure user ID is set
func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	if a.UserID == "" {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ODataBeforeUpdate validates API key update permissions
func (a *APIKey) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	// Get authenticated user from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user not authenticated")
	}

	// Users can only update their own API keys
	if a.UserID != userID {
		return fmt.Errorf("forbidden: cannot update API keys of other users")
	}

	// KeyPrefix, UserID, and CreatedAt are immutable
	// Check what fields are being updated - this is handled by gorm tags
	return nil
}

// ODataBeforeDelete validates API key deletion permissions
func (a *APIKey) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	// Get authenticated user from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user not authenticated")
	}

	// Users can only delete their own API keys
	if a.UserID != userID {
		return fmt.Errorf("forbidden: cannot delete API keys of other users")
	}

	return nil
}

// IsExpired checks if the API key has expired
func (a *APIKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// IsValid checks if the API key is valid (active and not expired)
func (a *APIKey) IsValid() bool {
	return a.IsActive && !a.IsExpired()
}

// GetPermissions returns the permissions as a string slice
func (a *APIKey) GetPermissions() []string {
	if a.Permissions == "" {
		return []string{}
	}
	var perms []string
	json.Unmarshal([]byte(a.Permissions), &perms)
	return perms
}

// SetPermissions sets the permissions from a string slice
func (a *APIKey) SetPermissions(perms []string) error {
	if perms == nil {
		a.Permissions = ""
		return nil
	}
	data, err := json.Marshal(perms)
	if err != nil {
		return err
	}
	a.Permissions = string(data)
	return nil
}
