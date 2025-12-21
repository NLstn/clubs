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
	"gorm.io/gorm"
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

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
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

	tokenID, err := GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}

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

// hashAPIKey creates a SHA-256 hash for indexed lookup (not for security).
// This hash is used for fast database queries with an index, while bcrypt
// remains the security-focused verification mechanism.
// Returns a hex-encoded 64-character string (SHA-256 produces 32 bytes, hex encoding doubles to 64 chars).
func hashAPIKey(keyStr string) string {
	sum := sha256.Sum256([]byte(keyStr))
	return hex.EncodeToString(sum[:])
}

// GenerateAPIKey creates a new API key with the specified prefix
// Returns: plainKey (shown once), keyHash (for storage), keyPrefix (for identification), keyHashSHA256 (for lookup), error
func GenerateAPIKey(prefix string) (string, string, string, string, error) {
	// Generate 32 bytes of cryptographically secure random data
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate random bytes: %w", err)
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
		return "", "", "", "", fmt.Errorf("failed to hash API key: %w", err)
	}
	keyHash := string(hashBytes)

	keyHashSHA256 := hashAPIKey(plainKey)

	// The keyPrefix is just the prefix part before the underscore
	// This is used to quickly filter candidates in ValidateAPIKey
	keyPrefix := prefix

	return plainKey, keyHash, keyPrefix, keyHashSHA256, nil
}

// ValidateAPIKey validates an API key and returns the associated user ID and permissions
func ValidateAPIKey(keyStr string) (string, []string, error) {
	if keyStr == "" {
		return "", nil, errors.New("API key is empty")
	}

	keyHashSHA256 := hashAPIKey(keyStr)

	var key struct {
		ID            string
		UserID        string
		KeyHash       string
		KeyHashSHA256 *string
		ExpiresAt     *time.Time
		Permissions   string
	}

	err := database.Db.Table("api_keys").
		Select("id, user_id, key_hash, key_hash_sha256, expires_at, permissions").
		Where("key_hash_sha256 = ?", keyHashSHA256).
		First(&key).Error
	if err == nil {
		if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
			return "", nil, errors.New("API key has expired")
		}

		var permissions []string
		if key.Permissions != "" {
			if err := json.Unmarshal([]byte(key.Permissions), &permissions); err != nil {
				log.Printf("failed to unmarshal API key permissions (api_key_id=%s, user_id=%s): %v", key.ID, key.UserID, err)
			}
		}

		keyID := key.ID
		keyHash := keyHashSHA256
		db := database.Db
		go func() {
			now := time.Now()
			if keyID != "" {
				if err := db.Table("api_keys").
					Where("id = ?", keyID).
					Update("last_used_at", now).Error; err != nil {
					log.Printf("failed to update last_used_at for api key id %s: %v", keyID, err)
				}
				return
			}
			if err := db.Table("api_keys").
				Where("key_hash_sha256 = ?", keyHash).
				Update("last_used_at", now).Error; err != nil {
				log.Printf("failed to update last_used_at for api key hash %s: %v", keyHash, err)
			}
		}()

		return key.UserID, permissions, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, fmt.Errorf("database query failed: %w", err)
	}

	// Extract prefix from the key to find candidates
	// The key format is: prefix_randomString
	// We need to extract the prefix part, which is everything before the last underscore + random data
	// Since the prefix itself might contain underscores (e.g., "sk_live"),
	// we rely on knowing the prefix was passed to GenerateAPIKey
	// The simplest approach: try to match keys by trying bcrypt on all keys
	// But for efficiency, we first filter by checking if the key starts with a known prefix pattern

	// Query all (non-deleted) API keys
	// Expired keys will be cleaned up by the CleanupExpiredAPIKeys job
	// Note: Rate limiting for failed authentication attempts is handled by the RateLimitMiddleware
	// applied to OData endpoints via AuthMiddleware in handlers/middlewares.go
	var keys []struct {
		ID          string
		UserID      string
		KeyHash     string
		KeyPrefix   string
		ExpiresAt   *time.Time
		Permissions string
	}

	err = database.Db.Table("api_keys").
		Select("id, user_id, key_hash, key_prefix, expires_at, permissions").
		Where("key_hash_sha256 IS NULL").
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
		ExpiresAt   *time.Time
		Permissions string
		KeyHash     string
	}

	for _, key := range candidates {
		err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(keyStr))
		if err == nil {
			// Found matching key
			result = &struct {
				ID          string
				UserID      string
				ExpiresAt   *time.Time
				Permissions string
				KeyHash     string
			}{
				ID:          key.ID,
				UserID:      key.UserID,
				ExpiresAt:   key.ExpiresAt,
				Permissions: key.Permissions,
				KeyHash:     key.KeyHash,
			}
			break
		}
	}

	if result == nil {
		return "", nil, errors.New("invalid API key")
	}

	// Check if the key has expired
	if result.ExpiresAt != nil && time.Now().After(*result.ExpiresAt) {
		return "", nil, errors.New("API key has expired")
	}

	// Parse permissions from JSON
	var permissions []string
	if result.Permissions != "" {
		if err := json.Unmarshal([]byte(result.Permissions), &permissions); err != nil {
			log.Printf("failed to unmarshal API key permissions (user_id=%s): %v", result.UserID, err)
		}
	}

	// Update last used timestamp asynchronously
	// For SQLite tests where ID might be empty, we need to use a different identifier
	// Capture specific values before spawning goroutine to avoid pointer issues
	resultID := result.ID
	resultKeyHash := result.KeyHash
	db := database.Db
	go func() {
		now := time.Now()
		updates := map[string]interface{}{
			"last_used_at":    now,
			"key_hash_sha256": keyHashSHA256,
		}
		// Update by user_id and key_hash since ID might be empty in SQLite
		if resultID != "" {
			if err := db.Table("api_keys").
				Where("id = ?", resultID).
				Updates(updates).Error; err != nil {
				log.Printf("failed to update last_used_at for api key id %s: %v", resultID, err)
			}
			return
		}
		if err := db.Table("api_keys").
			Where("key_hash = ?", resultKeyHash).
			Updates(updates).Error; err != nil {
			log.Printf("failed to update last_used_at for api key hash %s: %v", resultKeyHash, err)
		}
	}()

	return result.UserID, permissions, nil
}
