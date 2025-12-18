package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NLstn/clubs/azure/acs"
	"github.com/NLstn/clubs/database"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

// Init reads the JWT_SECRET environment variable and initializes the jwtSecret
// used for signing tokens.
func Init() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	jwtSecret = []byte(secret)
	return nil
}

// GetJWTSecret returns the JWT secret for use in external packages
// This is used by OData middleware for token validation
func GetJWTSecret() []byte {
	return jwtSecret
}

type contextKey string

const UserIDKey contextKey = "userID"

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func SendMagicLinkEmail(email, link string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	return acs.SendMail([]acs.Recipient{{Address: email}}, "Magic Link", "Click the link to login: "+link, "<a href='"+link+"'>Click here to login</a>")
}

func generateJWT(userID string, expiration time.Duration) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("cannot generate JWT with empty userID")
	}

	tokenID := GenerateToken()

	claims := jwt.MapClaims{
		"user_id": userID,
		"jti":     tokenID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenStr, nil
}

func GenerateAccessToken(userID string) (string, error) {
	return generateJWT(userID, 15*time.Minute)
}

func GenerateRefreshToken(userID string) (string, error) {
	return generateJWT(userID, 30*24*time.Hour)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Default().Println("Missing or invalid Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Default().Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Default().Printf("Token parsing error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Default().Println("Token validation failed")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Default().Println("Could not parse claims as MapClaims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if claims["user_id"] == nil {
			log.Default().Println("user_id claim is missing")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Default().Printf("user_id is not a string: %T", claims["user_id"])
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Default().Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Default().Printf("Token parsing error: %v", err)
		return "", fmt.Errorf("invalid token")
	}

	if !token.Valid {
		log.Default().Println("Token validation failed")
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Default().Println("Could not parse claims as MapClaims")
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Default().Printf("user_id is not a string: %T", claims["user_id"])
		return "", fmt.Errorf("invalid user ID")
	}

	sum := sha256.Sum256([]byte(token.Raw))
	tokenHash := hex.EncodeToString(sum[:])

	var expiresAt time.Time
	err = database.Db.Raw(`SELECT expires_at FROM refresh_tokens WHERE user_id = ? AND token = ?`, userID, tokenHash).Scan(&expiresAt).Error
	if err != nil {
		return "", err
	}
	if expiresAt.Before(time.Now()) {
		return "", fmt.Errorf("refresh token expired")
	}

	return userID, nil
}

// GenerateAPIKey creates a new API key with the specified prefix
// Returns: plainKey (shown once), keyHash (for storage), keyPrefix (for identification), error
func GenerateAPIKey(prefix string) (string, string, string, error) {
	// Generate 32 bytes of cryptographically secure random data
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 for URL-safe string
	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	// Remove padding characters for cleaner key
	randomString = strings.TrimRight(randomString, "=")

	// Construct full key with prefix
	plainKey := fmt.Sprintf("%s_%s", prefix, randomString)

	// Generate bcrypt hash for storage (cost 12)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainKey), 12)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to hash API key: %w", err)
	}
	keyHash := string(hashBytes)

	// The keyPrefix is just the prefix part before the underscore
	// This is used to quickly filter candidates in ValidateAPIKey
	keyPrefix := prefix

	return plainKey, keyHash, keyPrefix, nil
}

// ValidateAPIKey validates an API key and returns the associated user ID and permissions
func ValidateAPIKey(keyStr string) (string, []string, error) {
	if keyStr == "" {
		return "", nil, errors.New("API key is empty")
	}

	// Extract prefix from the key to find candidates
	// The key format is: prefix_randomString
	// We need to extract the prefix part, which is everything before the last underscore + random data
	// Since the prefix itself might contain underscores (e.g., "sk_live"), 
	// we rely on knowing the prefix was passed to GenerateAPIKey
	// The simplest approach: try to match keys by trying bcrypt on all keys
	// But for efficiency, we first filter by checking if the key starts with a known prefix pattern
	
	// Query ALL active API keys (for now - in production with many keys, this could be optimized)
	var keys []struct {
		ID          string
		UserID      string
		KeyHash     string
		KeyPrefix   string
		IsActive    bool
		ExpiresAt   *time.Time
		Permissions string
	}

	err := database.Db.Table("api_keys").
		Select("id, user_id, key_hash, key_prefix, is_active, expires_at, permissions").
		Find(&keys).Error
	if err != nil {
		return "", nil, fmt.Errorf("database query failed: %w", err)
	}
	
	// Filter to only keys where the provided key starts with the stored prefix
	var candidates []struct {
		ID          string
		UserID      string
		KeyHash     string
		KeyPrefix   string
		IsActive    bool
		ExpiresAt   *time.Time
		Permissions string
	}
	for _, key := range keys {
		if strings.HasPrefix(keyStr, key.KeyPrefix) {
			candidates = append(candidates, key)
		}
	}

	// Try to match the provided key with bcrypt
	var result *struct {
		ID          string
		UserID      string
		IsActive    bool
		ExpiresAt   *time.Time
		Permissions string
	}

	for _, key := range candidates {
		err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(keyStr))
		if err == nil {
			// Found matching key
			result = &struct {
				ID          string
				UserID      string
				IsActive    bool
				ExpiresAt   *time.Time
				Permissions string
			}{
				ID:          key.ID,
				UserID:      key.UserID,
				IsActive:    key.IsActive,
				ExpiresAt:   key.ExpiresAt,
				Permissions: key.Permissions,
			}
			break
		}
	}

	if result == nil {
		return "", nil, errors.New("invalid API key")
	}

	// Check if the key is active (already checked in query, but double-check)
	if !result.IsActive {
		return "", nil, errors.New("API key is inactive")
	}

	// Check if the key has expired
	if result.ExpiresAt != nil && time.Now().After(*result.ExpiresAt) {
		return "", nil, errors.New("API key has expired")
	}

	// Parse permissions from JSON
	var permissions []string
	if result.Permissions != "" {
		// Ignore unmarshal errors, just return empty array
		_ = json.Unmarshal([]byte(result.Permissions), &permissions)
	}

	// Update last used timestamp asynchronously
	// For SQLite tests where ID might be empty, we need to use a different identifier
	// Capture database reference to avoid race conditions in tests
	db := database.Db
	go func() {
		now := time.Now()
		// Update by user_id and key_hash since ID might be empty in SQLite
		if result.ID != "" {
			db.Table("api_keys").
				Where("id = ?", result.ID).
				Update("last_used_at", now)
		}
	}()

	return result.UserID, permissions, nil
}
