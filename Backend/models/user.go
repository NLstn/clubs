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

type RefreshToken struct {
	ID     string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID string `gorm:"type:uuid;not null"`
	Token  string `gorm:"uniqueIndex;not null"`
}

func FindOrCreateUser(email string) (User, error) {
	var user User
	err := database.Db.Raw(`SELECT * FROM users WHERE email = ?`, email).Scan(&user).Error
	if err != nil {
		return User{}, err
	}
	if user.ID == "" {
		user = User{Email: email}
		err = database.Db.Raw(`INSERT INTO users (email) VALUES (?) RETURNING *`, email).Scan(&user).Error
		if err != nil {
			return User{}, err
		}
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

func (u *User) UpdateUserName(name string) error {
	return database.Db.Exec(`UPDATE users SET name = ? WHERE id = ?`, name, u.ID).Error
}

func (u *User) StoreRefreshToken(token string) error {
	refreshToken := RefreshToken{UserID: u.ID, Token: token}
	return database.Db.Exec(`INSERT INTO refresh_tokens (user_id, token) VALUES (?, ?)`, refreshToken.UserID, refreshToken.Token).Error
}
