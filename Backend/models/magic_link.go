package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
)

type MagicLink struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"not null;uniqueIndex"`
	OTPCode   *string   `gorm:"uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}

func CreateMagicLink(email string) (string, error) {
	token, err := auth.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate magic link token: %w", err)
	}
	expiresAt := time.Now().Add(15 * time.Minute)

	tx := database.Db.Exec(`INSERT INTO magic_links (email, token, expires_at) VALUES (?, ?, ?)`,
		email, token, expiresAt)

	return token, tx.Error
}

// CreateMagicLinkWithCode generates a magic link token and a 6-digit OTP code.
// Both are stored server-side for verification. Returns token and code.
func CreateMagicLinkWithCode(email string) (string, string, error) {
	token, err := auth.GenerateToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate magic link token: %w", err)
	}
	// 6-digit numeric OTP (leading zeros allowed)
	// Use crypto/rand seeded math/rand for simplicity via auth.GenerateToken? Here derive digits safely.
	// Generate 6 digits by taking random bytes and mapping to 0-9.
	digits := make([]byte, 6)
	randTok, err := auth.GenerateToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate otp seed: %w", err)
	}
	for i := 0; i < 6; i++ {
		// Map char to 0-9 deterministically
		b := randTok[i]
		digits[i] = byte('0' + (b % 10))
	}
	code := string(digits)

	expiresAt := time.Now().Add(15 * time.Minute)

	tx := database.Db.Exec(`INSERT INTO magic_links (email, token, otp_code, expires_at) VALUES (?, ?, ?, ?)`,
		email, token, code, expiresAt)

	return token, code, tx.Error
}

func VerifyMagicLink(token string) (string, bool, error) {
	var result struct {
		Email string
	}

	tx := database.Db.Raw(
		`DELETE FROM magic_links WHERE token = ? AND expires_at > ? RETURNING email`,
		token,
		time.Now(),
	).Scan(&result)
	if tx.Error != nil {
		return "", false, tx.Error
	}
	if tx.RowsAffected == 0 {
		return "", false, nil
	}

	return result.Email, true, nil
}

// VerifyMagicCode verifies a 6-digit OTP code and consumes it on success.
func VerifyMagicCode(code string) (string, bool, error) {
	if code == "" {
		return "", false, nil
	}
	var result struct {
		Email string
	}
	tx := database.Db.Raw(
		`DELETE FROM magic_links WHERE otp_code = ? AND expires_at > ? RETURNING email`,
		code,
		time.Now(),
	).Scan(&result)
	if tx.Error != nil {
		return "", false, tx.Error
	}
	if tx.RowsAffected == 0 {
		return "", false, nil
	}
	return result.Email, true, nil
}

func DeleteMagicLink(token string) error {
	tx := database.Db.Exec(`DELETE FROM magic_links WHERE token = ?`, token)
	return tx.Error
}
