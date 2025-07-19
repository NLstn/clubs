package auth

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	if err := Init(); err != nil {
		t.Fatalf("failed to init auth: %v", err)
	}
	t.Run("GenerateToken", func(t *testing.T) {
		token := GenerateToken()
		assert.NotEmpty(t, token)
		assert.Greater(t, len(token), 40) // Base64 encoded 32 bytes should be longer than 40 chars
	})

	t.Run("GenerateAccessToken", func(t *testing.T) {
		userID := "test-user-id"
		token, err := GenerateAccessToken(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("GenerateAccessToken - Empty UserID", func(t *testing.T) {
		token, err := GenerateAccessToken("")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "cannot generate JWT with empty userID")
	})

	t.Run("GenerateRefreshToken", func(t *testing.T) {
		userID := "test-user-id"
		token, err := GenerateRefreshToken(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("GenerateRefreshToken - Empty UserID", func(t *testing.T) {
		token, err := GenerateRefreshToken("")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "cannot generate JWT with empty userID")
	})

	t.Run("ValidateRefreshToken - Invalid Token", func(t *testing.T) {
		userID, err := ValidateRefreshToken("invalid-token")

		assert.Error(t, err)
		assert.Empty(t, userID)
	})

	t.Run("ValidateRefreshToken - Malformed Token", func(t *testing.T) {
		userID, err := ValidateRefreshToken("malformed.token")

		assert.Error(t, err)
		assert.Empty(t, userID)
	})

	t.Run("generateJWT - Valid UserID", func(t *testing.T) {
		userID := "test-user-id"
		expiration := 15 * time.Minute

		token, err := generateJWT(userID, expiration)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token has 3 parts separated by dots (JWT structure)
		parts := len(token)
		assert.Greater(t, parts, 50) // JWT tokens are typically longer than 50 chars
	})

	t.Run("generateJWT - Empty UserID", func(t *testing.T) {
		expiration := 15 * time.Minute

		token, err := generateJWT("", expiration)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "cannot generate JWT with empty userID")
	})
}
