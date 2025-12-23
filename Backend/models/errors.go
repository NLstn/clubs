package models

import (
	"fmt"

	"github.com/NLstn/civo/database"
	"gorm.io/gorm"
)

// FeatureDisabledError represents an error when a club feature is disabled
// This error type is used to signal that an HTTP 400 should be returned
type FeatureDisabledError struct {
	Feature string
	Message string
}

// Error implements the error interface
func (e *FeatureDisabledError) Error() string {
	return e.Message
}

// NewFeatureDisabledError creates a new FeatureDisabledError
func NewFeatureDisabledError(featureName string) *FeatureDisabledError {
	return &FeatureDisabledError{
		Feature: featureName,
		Message: fmt.Sprintf("bad request: %s feature is disabled for this club", featureName),
	}
}

// CheckFeatureEnabled checks if a feature is enabled for a given club
// Returns a FeatureDisabledError if the feature is disabled
// If settings don't exist, returns nil (feature enabled by default for backward compatibility with tests)
func CheckFeatureEnabled(clubID, featureName string) error {
	var settings ClubSettings
	err := database.Db.Where("club_id = ?", clubID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// If settings don't exist, allow the operation (backward compatibility)
		// In production, settings should always exist for a club
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check club settings: %w", err)
	}

	switch featureName {
	case "fines":
		if !settings.FinesEnabled {
			return NewFeatureDisabledError(featureName)
		}
	case "shifts":
		if !settings.ShiftsEnabled {
			return NewFeatureDisabledError(featureName)
		}
	case "teams":
		if !settings.TeamsEnabled {
			return NewFeatureDisabledError(featureName)
		}
	case "news":
		if !settings.NewsEnabled {
			return NewFeatureDisabledError(featureName)
		}
	case "events":
		if !settings.EventsEnabled {
			return NewFeatureDisabledError(featureName)
		}
	default:
		return fmt.Errorf("unknown feature: %s", featureName)
	}

	return nil
}

// IsFeatureEnabled checks if a feature is enabled for a club without returning an error
func IsFeatureEnabled(clubID, featureName string) bool {
	return CheckFeatureEnabled(clubID, featureName) == nil
}
