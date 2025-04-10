package models

import (
	"time"

	"github.com/NLstn/clubs/database"
)

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindOrCreateUser(email string) (string, error) {
	var userID string
	err := database.Db.Raw(`SELECT id FROM users WHERE email = ?`, email).Scan(&userID).Error
	if userID == "" {
		err = database.Db.Raw(`INSERT INTO users (email) VALUES (?) RETURNING id`, email).Scan(&userID).Error
	}
	return userID, err
}
