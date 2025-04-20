package models

import (
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

func FindOrCreateUser(email string) (string, error) {
	var userID string
	err := database.Db.Raw(`SELECT id FROM users WHERE email = ?`, email).Scan(&userID).Error
	if userID == "" {
		err = database.Db.Raw(`INSERT INTO users (email) VALUES (?) RETURNING id`, email).Scan(&userID).Error
	}
	return userID, err
}

func GetUserByID(userID string) (User, error) {
	var user User
	err := database.Db.Raw(`SELECT * FROM users WHERE id = ?`, userID).Scan(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUserName(userID, name string) error {
	return database.Db.Exec(`UPDATE users SET name = ? WHERE id = ?`, name, userID).Error
}
