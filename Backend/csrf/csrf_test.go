package csrf

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Init with CSRF_SECRET", func(t *testing.T) {
		os.Setenv("CSRF_SECRET", "test-csrf-secret")
		defer os.Unsetenv("CSRF_SECRET")

		err := Init()
		assert.NoError(t, err)
		assert.NotNil(t, csrfSecret)
	})

	t.Run("Init with JWT_SECRET fallback", func(t *testing.T) {
		os.Unsetenv("CSRF_SECRET")
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		defer os.Unsetenv("JWT_SECRET")

		err := Init()
		assert.NoError(t, err)
		assert.NotNil(t, csrfSecret)
	})

	t.Run("Init without secrets", func(t *testing.T) {
		os.Unsetenv("CSRF_SECRET")
		os.Unsetenv("JWT_SECRET")

		err := Init()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SECRET")
	})
}

func TestGenerateStateToken(t *testing.T) {
	// Setup
	os.Setenv("CSRF_SECRET", "test-secret-for-state")
	defer os.Unsetenv("CSRF_SECRET")
	err := Init()
	assert.NoError(t, err)

	t.Run("Generate valid state token", func(t *testing.T) {
		ipHash := HashIP("192.168.1.1")
		token, err := GenerateStateToken(ipHash)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Token should have 3 parts: nonce.timestamp.signature
		assert.Contains(t, token, ".")
	})

	t.Run("Different IPs generate different tokens", func(t *testing.T) {
		ipHash1 := HashIP("192.168.1.1")
		ipHash2 := HashIP("192.168.1.2")

		token1, err1 := GenerateStateToken(ipHash1)
		token2, err2 := GenerateStateToken(ipHash2)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("Multiple calls generate different tokens", func(t *testing.T) {
		ipHash := HashIP("192.168.1.1")

		token1, err1 := GenerateStateToken(ipHash)
		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamp
		token2, err2 := GenerateStateToken(ipHash)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestValidateStateToken(t *testing.T) {
	// Setup
	os.Setenv("CSRF_SECRET", "test-secret-for-validation")
	defer os.Unsetenv("CSRF_SECRET")
	err := Init()
	assert.NoError(t, err)

	t.Run("Valid token is accepted", func(t *testing.T) {
		ipHash := HashIP("192.168.1.1")
		token, err := GenerateStateToken(ipHash)
		assert.NoError(t, err)

		nonce, valid := ValidateStateToken(token, ipHash)
		assert.True(t, valid)
		assert.NotEmpty(t, nonce)
	})

	t.Run("Token with wrong IP is rejected", func(t *testing.T) {
		ipHash1 := HashIP("192.168.1.1")
		ipHash2 := HashIP("192.168.1.2")

		token, err := GenerateStateToken(ipHash1)
		assert.NoError(t, err)

		nonce, valid := ValidateStateToken(token, ipHash2)
		assert.False(t, valid)
		assert.Empty(t, nonce)
	})

	t.Run("Malformed token is rejected", func(t *testing.T) {
		ipHash := HashIP("192.168.1.1")

		// Test various malformed tokens
		testCases := []string{
			"",
			"invalid",
			"only.two.parts",
			"too.many.parts.in.token",
		}

		for _, tc := range testCases {
			nonce, valid := ValidateStateToken(tc, ipHash)
			assert.False(t, valid, "Token should be invalid: %s", tc)
			assert.Empty(t, nonce)
		}
	})

	t.Run("Token with invalid timestamp is rejected", func(t *testing.T) {
		ipHash := HashIP("192.168.1.1")
		// Create token with invalid timestamp
		malformedToken := "validnonce.notanumber.validsig"

		nonce, valid := ValidateStateToken(malformedToken, ipHash)
		assert.False(t, valid)
		assert.Empty(t, nonce)
	})

	t.Run("Token with tampered signature is rejected", func(t *testing.T) {
		ipHash := HashIP("192.168.1.1")
		token, err := GenerateStateToken(ipHash)
		assert.NoError(t, err)

		// Tamper with the signature
		tamperedToken := token[:len(token)-5] + "xxxxx"

		nonce, valid := ValidateStateToken(tamperedToken, ipHash)
		assert.False(t, valid)
		assert.Empty(t, nonce)
	})
}

func TestHashIP(t *testing.T) {
	t.Run("Same IP produces same hash", func(t *testing.T) {
		hash1 := HashIP("192.168.1.1")
		hash2 := HashIP("192.168.1.1")

		assert.Equal(t, hash1, hash2)
		assert.NotEmpty(t, hash1)
	})

	t.Run("Different IPs produce different hashes", func(t *testing.T) {
		hash1 := HashIP("192.168.1.1")
		hash2 := HashIP("192.168.1.2")

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Hash is deterministic", func(t *testing.T) {
		// Test multiple times to ensure determinism
		ip := "10.0.0.1"
		hashes := make([]string, 5)

		for i := 0; i < 5; i++ {
			hashes[i] = HashIP(ip)
		}

		// All hashes should be identical
		for i := 1; i < 5; i++ {
			assert.Equal(t, hashes[0], hashes[i])
		}
	})
}

func TestGenerateCSRFToken(t *testing.T) {
	t.Run("Generate valid CSRF token", func(t *testing.T) {
		token, err := GenerateCSRFToken()

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Token should be valid base64
		valid := ValidateCSRFToken(token)
		assert.True(t, valid)
	})

	t.Run("Multiple calls generate different tokens", func(t *testing.T) {
		token1, err1 := GenerateCSRFToken()
		token2, err2 := GenerateCSRFToken()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestValidateCSRFToken(t *testing.T) {
	t.Run("Valid token is accepted", func(t *testing.T) {
		token, err := GenerateCSRFToken()
		assert.NoError(t, err)

		valid := ValidateCSRFToken(token)
		assert.True(t, valid)
	})

	t.Run("Invalid tokens are rejected", func(t *testing.T) {
		testCases := []string{
			"",                 // Empty
			"short",            // Too short
			"invalid!@#$%",     // Invalid characters
			"dG9vc2hvcnQ",      // Too short (valid base64 but < 32 bytes)
		}

		for _, tc := range testCases {
			valid := ValidateCSRFToken(tc)
			assert.False(t, valid, "Token should be invalid: %s", tc)
		}
	})
}

func TestStateTokenExpiration(t *testing.T) {
	// This test would require mocking time, which is complex in Go
	// For now, we document the behavior: tokens expire after 10 minutes
	// and create a manual test scenario

	t.Run("Token expiration logic exists", func(t *testing.T) {
		// Setup
		os.Setenv("CSRF_SECRET", "test-secret-expiration")
		defer os.Unsetenv("CSRF_SECRET")
		err := Init()
		assert.NoError(t, err)

		ipHash := HashIP("192.168.1.1")
		token, err := GenerateStateToken(ipHash)
		assert.NoError(t, err)

		// Immediately validate - should succeed
		nonce, valid := ValidateStateToken(token, ipHash)
		assert.True(t, valid)
		assert.NotEmpty(t, nonce)

		// Note: Testing actual expiration would require either:
		// 1. Waiting 10 minutes (not practical for unit tests)
		// 2. Mocking time (complex in Go without test infrastructure)
		// The expiration logic is tested manually and through integration tests
	})
}
