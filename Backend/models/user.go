package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
)

type User struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	CreatedBy string    `gorm:"type:uuid"`
	UpdatedAt time.Time
	UpdatedBy string    `gorm:"type:uuid"`
}

type RefreshToken struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	Token     string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
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
			
			// Update the created_by and updated_by fields with the user's own ID (self-reference)
			if user.ID != "" {
				user.CreatedBy = user.ID
				user.UpdatedBy = user.ID
				err = database.Db.Save(&user).Error
				if err != nil {
					return User{}, err
				}
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
	return database.Db.Exec(`UPDATE users SET name = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, name, u.ID, u.ID).Error
}

func (u *User) StoreRefreshToken(token string) error {
	// Delete all existing refresh tokens for this user first
	if err := u.DeleteAllRefreshTokens(); err != nil {
		fmt.Printf("Error deleting all refresh tokens for user %s: %v\n", u.ID, err)
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

func (u *User) GetUnpaidFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Raw(`SELECT * FROM fines WHERE user_id = ? AND paid = FALSE`, u.ID).Scan(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}
