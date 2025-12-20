package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OAuthState represents a state parameter used in OAuth flows for CSRF protection
type OAuthState struct {
	ID           string    `json:"ID" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	State        string    `json:"State" gorm:"type:varchar(255);unique;not null"`
	CodeVerifier string    `json:"CodeVerifier" gorm:"type:varchar(255);not null"` // PKCE code verifier
	ExpiresAt    time.Time `json:"ExpiresAt" gorm:"not null;index"`
	CreatedAt    time.Time `json:"CreatedAt" gorm:"autoCreateTime"`
}

// TableName overrides the table name used by OAuthState to `oauth_states`
func (OAuthState) TableName() string {
	return "oauth_states"
}

// BeforeCreate hook for OAuthState
func (o *OAuthState) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	// Default expiration: 10 minutes from creation
	if o.ExpiresAt.IsZero() {
		o.ExpiresAt = time.Now().Add(10 * time.Minute)
	}
	return nil
}

// CleanupExpiredOAuthStates removes expired OAuth state records from the database
// This should be called periodically to prevent the table from growing indefinitely
func CleanupExpiredOAuthStates() error {
	result := database.Db.Where("expires_at < ?", time.Now()).Delete(&OAuthState{})
	return result.Error
}

// CreateOAuthState creates a new OAuth state record
func CreateOAuthState(state, codeVerifier string) error {
	oauthState := &OAuthState{
		State:        state,
		CodeVerifier: codeVerifier,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}
	return database.Db.Create(oauthState).Error
}

// GetOAuthStateByState retrieves an OAuth state by its state parameter
func GetOAuthStateByState(state string) (*OAuthState, error) {
	var oauthState OAuthState
	err := database.Db.Where("state = ? AND expires_at > ?", state, time.Now()).First(&oauthState).Error
	if err != nil {
		return nil, err
	}
	return &oauthState, nil
}

// DeleteOAuthState deletes an OAuth state record (used after successful validation)
func DeleteOAuthState(state string) error {
	return database.Db.Where("state = ?", state).Delete(&OAuthState{}).Error
}
