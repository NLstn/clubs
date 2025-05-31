package models

import (
	"fmt"
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
	// Delete all existing refresh tokens for this user first
	if err := u.DeleteAllRefreshTokens(); err != nil {
		return err
	}
	
	refreshToken := RefreshToken{UserID: u.ID, Token: token, ExpiresAt: time.Now().Add(30 * 24 * time.Hour)}
	return database.Db.Exec(`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt).Error
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
