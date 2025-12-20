package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	csrfSecret []byte
)

// Init initializes the CSRF package with secrets from environment
func Init() error {
	secret := os.Getenv("CSRF_SECRET")
	if secret == "" {
		// Use JWT secret as fallback for CSRF
		secret = os.Getenv("JWT_SECRET")
		if secret == "" {
			return fmt.Errorf("CSRF_SECRET or JWT_SECRET environment variable is required")
		}
	}
	csrfSecret = []byte(secret)
	return nil
}

// GenerateStateToken creates a signed OAuth state token
// Format: nonce.timestamp.signature
// The signature covers: nonce + timestamp + ipHash
func GenerateStateToken(ipHash string) (string, error) {
	// Generate random nonce
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	nonce := base64.RawURLEncoding.EncodeToString(nonceBytes)

	// Current timestamp
	timestamp := time.Now().Unix()

	// Create signature: HMAC-SHA256(nonce + timestamp + ipHash)
	message := fmt.Sprintf("%s.%d.%s", nonce, timestamp, ipHash)
	signature := generateHMAC(message)

	// Format: nonce.timestamp.signature
	stateToken := fmt.Sprintf("%s.%d.%s", nonce, timestamp, signature)

	return stateToken, nil
}

// ValidateStateToken verifies a signed OAuth state token
// Returns the nonce if valid, empty string otherwise
func ValidateStateToken(stateToken string, ipHash string) (string, bool) {
	// Parse token: nonce.timestamp.signature
	parts := strings.Split(stateToken, ".")
	if len(parts) != 3 {
		return "", false
	}

	nonce := parts[0]
	timestampStr := parts[1]
	signature := parts[2]

	// Parse timestamp
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", false
	}

	// Check if token has expired (10 minutes)
	tokenTime := time.Unix(timestamp, 0)
	if time.Since(tokenTime) > 10*time.Minute {
		return "", false
	}

	// Verify signature
	message := fmt.Sprintf("%s.%d.%s", nonce, timestamp, ipHash)
	expectedSignature := generateHMAC(message)

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", false
	}

	return nonce, true
}

// HashIP creates a SHA-256 hash of an IP address for privacy
func HashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}

// generateHMAC creates an HMAC-SHA256 signature
func generateHMAC(message string) string {
	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// GenerateCSRFToken generates a random CSRF token for general API protection
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ValidateCSRFToken validates a CSRF token (for now, just checks format)
// In a stateless architecture, CSRF tokens could be signed similar to state tokens
func ValidateCSRFToken(token string) bool {
	// Basic validation - token should be a valid base64 string of sufficient length
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	// Require at least 32 bytes (256 bits) of entropy
	return len(decoded) >= 32
}
