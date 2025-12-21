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

func DeleteMagicLink(token string) error {
	tx := database.Db.Exec(`DELETE FROM magic_links WHERE token = ?`, token)
	return tx.Error
}
