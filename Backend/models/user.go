package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NLstn/clubs/database"
)

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string
	Email     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshToken struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	Token     string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
	UserAgent string    // Browser/device information
	IPAddress string    // IP address for session tracking
	CreatedAt time.Time
}

func FindOrCreateUser(email string) (User, error) {
	var user User
	err := database.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err.Error() == "record not found" {
			// For new user creation, use raw SQL to insert with NULL for created_by/updated_by initially
			// Then update after we have the user ID
			err = database.Db.Raw(`INSERT INTO users (email, created_by, updated_by) VALUES (?, NULL, NULL) RETURNING *`, email).Scan(&user).Error
			if err != nil {
				// If the above fails (e.g., in SQLite), try with GORM and handle the self-reference after
				user = User{Email: email}
				err = database.Db.Create(&user).Error
				if err != nil {
					return User{}, err
				}
			}

			err = database.Db.Save(&user).Error
			if err != nil {
				return User{}, err
			}
			return user, nil
		}
		return User{}, err
	}
	return user, nil
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

func (u *User) UpdateUserName(name string) error {
	return database.Db.Exec(`UPDATE users SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, name, u.ID).Error
}

func (u *User) StoreRefreshToken(token, userAgent, ipAddress string) error {
	refreshToken := RefreshToken{
		UserID:    u.ID,
		Token:     token,
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
	err := database.Db.Raw(`SELECT * FROM refresh_tokens WHERE user_id = ? AND token = ?`, u.ID, token).Scan(&refreshToken).Error
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
	return database.Db.Exec(`DELETE FROM refresh_tokens WHERE user_id = ? AND token = ?`, u.ID, token).Error
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
