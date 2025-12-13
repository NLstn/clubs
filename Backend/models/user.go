package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"gorm.io/gorm"
)

type User struct {
	ID         string     `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	FirstName  string     `json:"FirstName" odata:"required"`
	LastName   string     `json:"LastName" odata:"required"`
	Email      string     `json:"Email" gorm:"uniqueIndex;not null" odata:"required"`
	KeycloakID *string    `json:"KeycloakID,omitempty" gorm:"uniqueIndex" odata:"nullable"`
	BirthDate  *time.Time `json:"BirthDate,omitempty" gorm:"type:date" odata:"nullable"`
	CreatedAt  time.Time  `json:"CreatedAt" odata:"immutable"`
	UpdatedAt  time.Time  `json:"UpdatedAt"`
}

type RefreshToken struct {
	ID        string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `json:"UserID" gorm:"type:uuid;not null"`
	Token     string    `json:"Token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"ExpiresAt"`
	UserAgent string    `json:"UserAgent"`
	IPAddress string    `json:"IPAddress"`
	CreatedAt time.Time `json:"CreatedAt"`
}

// UserSession represents an active user session exposed via OData API
// This maps to RefreshToken internally but provides a cleaner API
type UserSession struct {
	ID        string    `json:"ID" gorm:"column:id;type:uuid;primaryKey" odata:"key"`
	UserID    string    `json:"UserID" gorm:"column:user_id;type:uuid;not null" odata:"immutable"`
	UserAgent string    `json:"UserAgent" gorm:"column:user_agent" odata:"immutable"`
	IPAddress string    `json:"IPAddress" gorm:"column:ip_address" odata:"immutable"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"column:created_at" odata:"immutable"`
	ExpiresAt time.Time `json:"ExpiresAt" gorm:"column:expires_at" odata:"immutable"`
	IsCurrent bool      `json:"IsCurrent" gorm:"-" odata:"-"` // Computed field
}

// TableName specifies that UserSession uses the refresh_tokens table
func (UserSession) TableName() string {
	return "refresh_tokens"
}

// ODataBeforeReadCollection filters sessions to only show current user's active sessions
// This implements the ReadHook interface from go-odata
func (UserSession) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return []func(*gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("1 = 0") // Return no results if no user context
			},
		}, nil
	}

	return []func(*gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ? AND expires_at > NOW()", userID)
		},
	}, nil
}

// ODataBeforeReadEntity filters to only allow current user to read their own sessions
// This implements the ReadHook interface from go-odata
func (UserSession) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return []func(*gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("1 = 0") // Return no results if no user context
			},
		}, nil
	}

	// For DELETE operations, allow reading expired sessions to enable deletion
	// The OData framework loads the entity before calling BeforeDelete
	if r.Method == "DELETE" {
		return []func(*gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("user_id = ?", userID)
			},
		}, nil
	}

	return []func(*gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ? AND expires_at > NOW()", userID)
		},
	}, nil
}

// ODataAfterReadCollection adds the IsCurrent computed field to each session
// This implements the ReadHook interface from go-odata
func (UserSession) ODataAfterReadCollection(ctx context.Context, r *http.Request, opts interface{}, results interface{}) (interface{}, error) {
	// Try both pointer to slice and slice directly
	var sessions []UserSession
	if sessionsPtr, ok := results.(*[]UserSession); ok {
		sessions = *sessionsPtr
	} else if sessionsSlice, ok := results.([]UserSession); ok {
		sessions = sessionsSlice
	} else {
		return results, nil // Unknown type, return as-is
	}

	if len(sessions) == 0 {
		return results, nil
	}

	// Get current refresh token to identify current session
	currentRefreshToken := r.Header.Get("X-Refresh-Token")
	if currentRefreshToken == "" {
		return results, nil // No current token, all sessions have IsCurrent = false
	}

	hashedToken := HashToken(currentRefreshToken)

	// Get tokens for all sessions to compare
	ids := make([]string, len(sessions))
	for i, session := range sessions {
		ids[i] = session.ID
	}

	var tokens []struct {
		ID    string
		Token string
	}
	if err := database.Db.Raw("SELECT id, token FROM refresh_tokens WHERE id IN ?", ids).Scan(&tokens).Error; err == nil {
		tokenMap := make(map[string]string)
		for _, t := range tokens {
			tokenMap[t.ID] = t.Token
		}

		for i := range sessions {
			if token, ok := tokenMap[sessions[i].ID]; ok {
				sessions[i].IsCurrent = token == hashedToken
			}
		}
	}

	// Return the modified sessions in the same format we received
	if _, ok := results.(*[]UserSession); ok {
		return &sessions, nil
	}
	return sessions, nil
}

// ODataAfterReadEntity adds the IsCurrent computed field to the session
// This implements the ReadHook interface from go-odata
func (UserSession) ODataAfterReadEntity(ctx context.Context, r *http.Request, opts interface{}, entity interface{}) (interface{}, error) {
	session, ok := entity.(*UserSession)
	if !ok {
		return entity, nil
	}

	// Get current refresh token to identify current session
	currentRefreshToken := r.Header.Get("X-Refresh-Token")
	if currentRefreshToken == "" {
		return entity, nil // No current token, IsCurrent = false
	}

	hashedToken := HashToken(currentRefreshToken)

	// Check if this session matches the current token
	var tokenCheck struct{ Token string }
	if err := database.Db.Raw("SELECT token FROM refresh_tokens WHERE id = ?", session.ID).Scan(&tokenCheck).Error; err == nil {
		session.IsCurrent = tokenCheck.Token == hashedToken
	}

	return entity, nil
}

// ODataBeforeDelete ensures users can only delete their own sessions
// This implements the EntityHook interface from go-odata
func (s *UserSession) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		fmt.Printf("DEBUG: BeforeDelete - userID not found in context, ok=%v, userID=%v\n", ok, userID)
		return fmt.Errorf("unauthorized")
	}

	fmt.Printf("DEBUG: BeforeDelete - userID from context: %s, session.ID: %s, session.UserID: %s\n", userID, s.ID, s.UserID)

	// The ID field is already populated by OData framework
	// We just need to verify it belongs to the current user by checking UserID
	// Note: The OData framework loads the entity before calling this hook,
	// so s.UserID should be populated if BeforeRead hooks filtered correctly
	if s.UserID == "" {
		// If UserID is not populated, query it from database
		var session UserSession
		if err := database.Db.Where("id = ?", s.ID).First(&session).Error; err != nil {
			fmt.Printf("DEBUG: BeforeDelete - session not found in DB for ID: %s, error: %v\n", s.ID, err)
			return fmt.Errorf("session not found")
		}
		s.UserID = session.UserID
		fmt.Printf("DEBUG: BeforeDelete - loaded UserID from DB: %s\n", s.UserID)
	}

	// Verify the session belongs to the current user
	if s.UserID != userID {
		fmt.Printf("DEBUG: BeforeDelete - user mismatch: session.UserID=%s, contextUserID=%s\n", s.UserID, userID)
		return fmt.Errorf("cannot delete another user's session")
	}

	fmt.Printf("DEBUG: BeforeDelete - authorization successful for session %s\n", s.ID)
	return nil
}

// HashToken returns a sha256 hash of the provided token encoded as hex.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func FindOrCreateUser(email string) (User, error) {
	var user User
	err := database.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err.Error() == "record not found" {
			user = User{Email: email}
			err = database.Db.Create(&user).Error
			if err != nil {
				return User{}, err
			}
			return user, nil
		}
		return User{}, err
	}
	return user, nil
}

func FindOrCreateUserWithKeycloakID(keycloakID, email, fullName string) (User, error) {
	var user User

	// First try to find user by Keycloak ID
	err := database.Db.Where("keycloak_id = ?", keycloakID).First(&user).Error
	if err == nil {
		// User found, update email if different
		if user.Email != email {
			user.Email = email
			database.Db.Save(&user)
		}
		return user, nil
	}

	// If not found by Keycloak ID, try to find by email
	err = database.Db.Where("email = ?", email).First(&user).Error
	if err == nil {
		// User exists with this email, update with Keycloak ID
		user.KeycloakID = &keycloakID
		database.Db.Save(&user)
		return user, nil
	}

	// User doesn't exist, create new one
	if err.Error() == "record not found" {
		// Parse full name into first and last name
		firstName, lastName := parseFullName(fullName)

		user = User{
			Email:      email,
			KeycloakID: &keycloakID,
			FirstName:  firstName,
			LastName:   lastName,
		}
		err = database.Db.Create(&user).Error
		if err != nil {
			return User{}, err
		}
		return user, nil
	}

	return User{}, err
}

func parseFullName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func GetUserByID(userID string) (User, error) {
	var user User
	err := database.Db.Raw(`SELECT * FROM users WHERE id = ?`, userID).Scan(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUsersByIDs(userIDs []string) ([]User, error) {
	var users []User
	if len(userIDs) == 0 {
		return users, nil
	}
	err := database.Db.Where("id IN ?", userIDs).Find(&users).Error
	return users, err
}

func (u *User) UpdateUserName(firstName, lastName string) error {
	return database.Db.Exec(`UPDATE users SET first_name = ?, last_name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, firstName, lastName, u.ID).Error
}

func (u *User) UpdateBirthDate(birthDate *time.Time) error {
	return database.Db.Exec(`UPDATE users SET birth_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, birthDate, u.ID).Error
}

// Helper method to get full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// Helper method to check if user has completed profile setup
func (u *User) IsProfileComplete() bool {
	return u.FirstName != "" && u.LastName != ""
}

func (u *User) StoreRefreshToken(token, userAgent, ipAddress string) error {
	// Delete any existing refresh tokens for the same user and IP address
	err := database.Db.Exec(`DELETE FROM refresh_tokens WHERE user_id = ? AND ip_address = ?`, u.ID, ipAddress).Error
	if err != nil {
		return fmt.Errorf("failed to delete existing sessions for IP %s: %v", ipAddress, err)
	}

	refreshToken := RefreshToken{
		UserID:    u.ID,
		Token:     HashToken(token),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}
	return database.Db.Exec(`INSERT INTO refresh_tokens (user_id, token, expires_at, user_agent, ip_address, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, refreshToken.UserAgent, refreshToken.IPAddress, refreshToken.CreatedAt).Error
}

// Helper function to extract device information from HTTP request
func GetDeviceInfo(r *http.Request) (userAgent, ipAddress string) {
	userAgent = r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}

	// Try to get real IP from various headers (common proxy headers)
	ipAddress = r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Real-IP")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// Clean up the IP address (remove port if present)
	if colon := strings.LastIndex(ipAddress, ":"); colon != -1 {
		if bracket := strings.LastIndex(ipAddress, "]"); bracket == -1 || bracket < colon {
			ipAddress = ipAddress[:colon]
		}
	}

	if ipAddress == "" {
		ipAddress = "Unknown"
	}

	return userAgent, ipAddress
}

func (u *User) ValidateRefreshToken(token string) error {
	var refreshToken RefreshToken
	err := database.Db.Raw(`SELECT * FROM refresh_tokens WHERE user_id = ? AND token = ?`, u.ID, HashToken(token)).Scan(&refreshToken).Error
	if err != nil {
		return err
	}
	if refreshToken.ID == "" {
		return fmt.Errorf("invalid refresh token")
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("refresh token expired")
	}

	return nil
}

func (u *User) DeleteRefreshToken(token string) error {
	return database.Db.Exec(`DELETE FROM refresh_tokens WHERE user_id = ? AND token = ?`, u.ID, HashToken(token)).Error
}

func (u *User) DeleteAllRefreshTokens() error {
	return database.Db.Exec(`DELETE FROM refresh_tokens WHERE user_id = ?`, u.ID).Error
}

func (u *User) GetFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Raw(`SELECT * FROM fines WHERE user_id = ?`, u.ID).Scan(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}

func (u *User) GetUnpaidFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Raw(`SELECT * FROM fines WHERE user_id = ? AND paid = FALSE`, u.ID).Scan(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}

func (u *User) GetActiveSessions() ([]RefreshToken, error) {
	var sessions []RefreshToken
	err := database.Db.Raw(`SELECT * FROM refresh_tokens WHERE user_id = ? AND expires_at > ?`, u.ID, time.Now()).Scan(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (u *User) DeleteSession(sessionID string) error {
	return database.Db.Exec(`DELETE FROM refresh_tokens WHERE user_id = ? AND id = ?`, u.ID, sessionID).Error
}
