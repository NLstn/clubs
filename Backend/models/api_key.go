package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"gorm.io/gorm"
)

// APIKey represents a long-lived API key for programmatic access
type APIKey struct {
	ID            string     `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	UserID        string     `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
	User          User       `json:"User,omitempty" gorm:"foreignKey:UserID" odata:"navigationProperty"`
	Name          string     `json:"Name" gorm:"not null" odata:"required"`
	KeyHash       string     `json:"-" gorm:"uniqueIndex;not null"`      // Never exposed via API
	KeyHashSHA256 *string    `json:"-" gorm:"uniqueIndex;type:char(64)"` // Indexed lookup hash for API key (SHA-256: 32 bytes -> 64 hex chars, hence char(64))
	KeyPrefix     string     `json:"KeyPrefix" gorm:"not null" odata:"immutable"`
	Permissions   string     `json:"-" gorm:"type:text"` // Stored as JSON string
	LastUsedAt    *time.Time `json:"LastUsedAt,omitempty" gorm:"type:timestamp" odata:"nullable"`
	ExpiresAt     *time.Time `json:"ExpiresAt,omitempty" gorm:"type:timestamp" odata:"nullable"`
	CreatedAt     time.Time  `json:"CreatedAt" odata:"immutable"`
	UpdatedAt     time.Time  `json:"UpdatedAt"`
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

// ODataBeforeReadCollection filters API keys to only those belonging to the user
// This prevents users from enumerating all API keys in the system
func (a APIKey) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user not authenticated")
	}

	// User can only see their own API keys
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific API key record
// Users can only read their own API keys
func (a APIKey) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user not authenticated")
	}

	// User can only see their own API keys
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates API key creation permissions
func (a *APIKey) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	// Get authenticated user from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user not authenticated")
	}

	// Users can only create API keys for themselves
	if a.UserID == "" {
		a.UserID = userID
	} else if a.UserID != userID {
		return fmt.Errorf("forbidden: cannot create API keys for other users")
	}

	// Set CreatedAt and UpdatedAt
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

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

// IsValid checks if the API key is valid (not expired)
func (a *APIKey) IsValid() bool {
	return !a.IsExpired()
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

// CleanupExpiredAPIKeys removes expired API keys from the database
// This should be called periodically to prevent the table from growing indefinitely
// Expired keys are hard-deleted rather than being marked inactive
func CleanupExpiredAPIKeys() error {
	var db gorm.DB = *database.Db
	result := db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&APIKey{})
	return result.Error
}
