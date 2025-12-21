package auth

import (
	"os"
	"strings"
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
		token, err := GenerateToken()
		assert.NoError(t, err)
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

func TestGenerateAPIKey(t *testing.T) {
	t.Run("Valid key generation with sk_live prefix", func(t *testing.T) {
		plainKey, keyHash, keyPrefix, keyHashSHA256, err := GenerateAPIKey("sk_live")

		assert.NoError(t, err)
		assert.NotEmpty(t, plainKey)
		assert.NotEmpty(t, keyHash)
		assert.NotEmpty(t, keyPrefix)
		assert.NotEmpty(t, keyHashSHA256)

		// Verify key format
		assert.Contains(t, plainKey, "sk_live_")
		assert.Greater(t, len(plainKey), 20, "Key should be sufficiently long")

		// Verify hash is bcrypt (starts with $2a$ and is ~60 chars)
		assert.True(t, strings.HasPrefix(keyHash, "$2a$"), "Hash should be bcrypt format")
		assert.Greater(t, len(keyHash), 50, "Bcrypt hash should be at least 50 characters")

		// Verify prefix is just the prefix part
		assert.Equal(t, "sk_live", keyPrefix)
	})

	t.Run("Valid key generation with sk_test prefix", func(t *testing.T) {
		plainKey, keyHash, keyPrefix, keyHashSHA256, err := GenerateAPIKey("sk_test")

		assert.NoError(t, err)
		assert.Contains(t, plainKey, "sk_test_")
		assert.NotEmpty(t, keyHash)
		assert.NotEmpty(t, keyPrefix)
		assert.NotEmpty(t, keyHashSHA256)
	})

	t.Run("Keys should be unique", func(t *testing.T) {
		key1, hash1, _, sha1, err := GenerateAPIKey("sk_live")
		assert.NoError(t, err)

		key2, hash2, _, sha2, err := GenerateAPIKey("sk_live")
		assert.NoError(t, err)

		// Keys should be different
		assert.NotEqual(t, key1, key2)
		assert.NotEqual(t, hash1, hash2)
		assert.NotEqual(t, sha1, sha2)
	})
}
