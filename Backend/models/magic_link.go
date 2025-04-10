package models

import (
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
)

type MagicLink struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}

func CreateMagicLink(email string) (string, error) {
	token := auth.GenerateToken()
	expiresAt := time.Now().Add(15 * time.Minute)

	tx := database.Db.Exec(`INSERT INTO magic_links (email, token, expires_at) VALUES (?, ?, ?)`,
		email, token, expiresAt)

	return token, tx.Error
}

func VerifyMagicLink(token string) (string, bool, error) {
	var result struct {
		Email     string
		ExpiresAt time.Time
	}

	err := database.Db.Raw(`SELECT email, expires_at FROM magic_links WHERE token = ?`, token).
		Scan(&result).Error

	if err != nil || time.Now().After(result.ExpiresAt) {
		return "", false, err
	}

	return result.Email, true, nil
}

func DeleteMagicLink(token string) error {
	tx := database.Db.Exec(`DELETE FROM magic_links WHERE token = ?`, token)
	return tx.Error
}
