package models_test

import (
	"testing"
	"time"

	"github.com/NLstn/civo/handlers"
	"github.com/NLstn/civo/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateOAuthState(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	state := "test-state-123"
	codeVerifier := "test-verifier-456"

	err := models.CreateOAuthState(state, codeVerifier)
	assert.NoError(t, err)

	// Verify the state was created
	oauthState, err := models.GetOAuthStateByState(state)
	assert.NoError(t, err)
	assert.Equal(t, state, oauthState.State)
	assert.Equal(t, codeVerifier, oauthState.CodeVerifier)
	assert.False(t, oauthState.ExpiresAt.IsZero())
	assert.True(t, oauthState.ExpiresAt.After(time.Now()))
}

func TestGetOAuthStateByState(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("existing state", func(t *testing.T) {
		state := "valid-state"
		codeVerifier := "valid-verifier"

		err := models.CreateOAuthState(state, codeVerifier)
		assert.NoError(t, err)

		oauthState, err := models.GetOAuthStateByState(state)
		assert.NoError(t, err)
		assert.Equal(t, state, oauthState.State)
		assert.Equal(t, codeVerifier, oauthState.CodeVerifier)
	})

	t.Run("non-existing state", func(t *testing.T) {
		_, err := models.GetOAuthStateByState("non-existing-state")
		assert.Error(t, err)
	})

	t.Run("expired state", func(t *testing.T) {
		// Create an expired state directly in the database
		expiredState := &models.OAuthState{
			State:        "expired-state",
			CodeVerifier: "expired-verifier",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		err := handlers.GetDB().Create(expiredState).Error
		assert.NoError(t, err)

		// Should not be able to retrieve expired state
		_, err = models.GetOAuthStateByState("expired-state")
		assert.Error(t, err)
	})
}

func TestDeleteOAuthState(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	state := "delete-state"
	codeVerifier := "delete-verifier"

	err := models.CreateOAuthState(state, codeVerifier)
	assert.NoError(t, err)

	// Verify it exists
	oauthState, err := models.GetOAuthStateByState(state)
	assert.NoError(t, err)
	assert.NotNil(t, oauthState)

	// Delete it
	err = models.DeleteOAuthState(state)
	assert.NoError(t, err)

	// Verify it no longer exists
	_, err = models.GetOAuthStateByState(state)
	assert.Error(t, err)
}

func TestCleanupExpiredOAuthStates(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create some OAuth states
	// Valid state 1
	err := models.CreateOAuthState("valid-1", "verifier-1")
	assert.NoError(t, err)

	// Valid state 2
	err = models.CreateOAuthState("valid-2", "verifier-2")
	assert.NoError(t, err)

	// Expired state 1
	expiredState1 := &models.OAuthState{
		State:        "expired-1",
		CodeVerifier: "expired-verifier-1",
		ExpiresAt:    time.Now().Add(-1 * time.Hour),
	}
	err = handlers.GetDB().Create(expiredState1).Error
	assert.NoError(t, err)

	// Expired state 2
	expiredState2 := &models.OAuthState{
		State:        "expired-2",
		CodeVerifier: "expired-verifier-2",
		ExpiresAt:    time.Now().Add(-2 * time.Hour),
	}
	err = handlers.GetDB().Create(expiredState2).Error
	assert.NoError(t, err)

	// Verify we have 4 states total
	var countBefore int64
	handlers.GetDB().Model(&models.OAuthState{}).Count(&countBefore)
	assert.Equal(t, int64(4), countBefore)

	// Cleanup expired states
	err = models.CleanupExpiredOAuthStates()
	assert.NoError(t, err)

	// Verify we only have 2 valid states remaining
	var countAfter int64
	handlers.GetDB().Model(&models.OAuthState{}).Count(&countAfter)
	assert.Equal(t, int64(2), countAfter)

	// Verify the valid states still exist
	_, err = models.GetOAuthStateByState("valid-1")
	assert.NoError(t, err)

	_, err = models.GetOAuthStateByState("valid-2")
	assert.NoError(t, err)

	// Verify the expired states are gone
	var expired1 models.OAuthState
	err = handlers.GetDB().Where("state = ?", "expired-1").First(&expired1).Error
	assert.Error(t, err) // Should be record not found

	var expired2 models.OAuthState
	err = handlers.GetDB().Where("state = ?", "expired-2").First(&expired2).Error
	assert.Error(t, err) // Should be record not found
}

func TestOAuthStateBeforeCreate(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("auto-generates ID", func(t *testing.T) {
		state := &models.OAuthState{
			State:        "test-state-id",
			CodeVerifier: "test-verifier",
			ExpiresAt:    time.Now().Add(10 * time.Minute),
		}
		err := handlers.GetDB().Create(state).Error
		assert.NoError(t, err)
		assert.NotEmpty(t, state.ID)
	})

	t.Run("sets default expiration", func(t *testing.T) {
		state := &models.OAuthState{
			State:        "test-state-expiry",
			CodeVerifier: "test-verifier",
		}
		err := handlers.GetDB().Create(state).Error
		assert.NoError(t, err)
		assert.False(t, state.ExpiresAt.IsZero())
		assert.True(t, state.ExpiresAt.After(time.Now()))
		
		// Should be approximately 10 minutes from now (allowing 1 minute margin)
		expectedExpiry := time.Now().Add(10 * time.Minute)
		timeDiff := state.ExpiresAt.Sub(expectedExpiry).Abs()
		assert.Less(t, timeDiff, 1*time.Minute)
	})
}
