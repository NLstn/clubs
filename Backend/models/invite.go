package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

// Invite represents an admin invitation to a user to join a club
type Invite struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID    string    `json:"club_id" gorm:"type:uuid"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateInvite creates a new invitation from an admin to a user
func (c *Club) CreateInvite(email, createdBy string) error {
	invite := &Invite{
		ID:        uuid.New().String(),
		ClubID:    c.ID,
		Email:     email,
		CreatedBy: createdBy,
	}
	return database.Db.Create(invite).Error
}

// GetInvites returns all pending invites for a club (admin view)
func (c *Club) GetInvites() ([]Invite, error) {
	var invites []Invite
	err := database.Db.Where("club_id = ?", c.ID).Find(&invites).Error
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// GetUserInvites returns all pending invites for a user
func (u *User) GetUserInvites() ([]Invite, error) {
	var invites []Invite
	err := database.Db.Where("email = (SELECT email FROM users WHERE id = ?)", u.ID).Find(&invites).Error
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// AcceptInvite accepts an invitation and adds the user to the club
func AcceptInvite(inviteId, userId string) error {
	var invite Invite
	err := database.Db.Where("id = ?", inviteId).First(&invite).Error
	if err != nil {
		return err
	}

	var club Club
	err = database.Db.Where("id = ?", invite.ClubID).First(&club).Error
	if err != nil {
		return err
	}

	// Get the user accepting the invite
	var user User
	err = database.Db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return err
	}

	// Verify the user's email matches the invite
	if user.Email != invite.Email {
		return fmt.Errorf("user email does not match invite")
	}

	// Add user to club
	err = club.AddMember(user.ID, "member")
	if err != nil {
		return err
	}

	// Delete the invite since it's now complete
	return database.Db.Delete(&Invite{}, "id = ?", inviteId).Error
}

// RejectInvite rejects an invitation
func RejectInvite(inviteId string) error {
	return database.Db.Delete(&Invite{}, "id = ?", inviteId).Error
}

// CanUserEditInvite checks if a user can accept/reject an invite
func (u *User) CanUserEditInvite(inviteId string) (bool, error) {
	var user User
	err := database.Db.Where("id = ?", u.ID).First(&user).Error
	if err != nil {
		return false, err
	}

	var invite Invite
	err = database.Db.Where("id = ?", inviteId).First(&invite).Error
	if err != nil {
		return false, err
	}

	// User can accept if it's their own invite (their email matches)
	return user.Email == invite.Email, nil
}
