package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invite represents an admin invitation to a user to join a club
type Invite struct {
	ID        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid" odata:"required"`
	Email     string    `json:"Email" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// CreateInvite creates a new invitation from an admin to a user
func (c *Club) CreateInvite(email, createdBy string) error {
	invite := &Invite{
		ID:        uuid.New().String(),
		ClubID:    c.ID,
		Email:     email,
		CreatedBy: createdBy,
	}
	err := database.Db.Create(invite).Error
	if err != nil {
		return err
	}

	// Send invite notification to the user
	SendInviteReceivedNotifications(email, c.ID, c.Name, invite.ID)

	return nil
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

	// Add user to club via invite (this will skip the member_added notification)
	err = club.AddMemberViaInvite(user.ID, "member")
	if err != nil {
		return err
	}

	// Remove any invite notifications for this invite
	RemoveInviteNotifications(inviteId)

	// Delete the invite since it's now complete
	return database.Db.Delete(&Invite{}, "id = ?", inviteId).Error
}

// RejectInvite rejects an invitation
func RejectInvite(inviteId string) error {
	// Remove any invite notifications for this invite
	RemoveInviteNotifications(inviteId)

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

// ODataBeforeReadCollection filters invites - users see invites for their email or clubs they admin
func (i Invite) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get user email to check for personal invites
	var user User
	if err := database.Db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// User can see invites for their email OR invites for clubs they are admin/owner of
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ? OR club_id IN (SELECT club_id FROM members WHERE user_id = ? AND role IN ('admin', 'owner'))", user.Email, userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific invite
func (i Invite) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get user email to check for personal invites
	var user User
	if err := database.Db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// User can see invites for their email OR invites for clubs they are admin/owner of
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ? OR club_id IN (SELECT club_id FROM members WHERE user_id = ? AND role IN ('admin', 'owner'))", user.Email, userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates invite creation permissions
func (i *Invite) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", i.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create invites")
	}

	// Set CreatedBy
	now := time.Now()
	i.CreatedAt = now
	i.UpdatedAt = now
	i.CreatedBy = userID

	return nil
}

// ODataBeforeDelete validates invite deletion permissions
func (i *Invite) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get user email
	var user User
	if err := database.Db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// User can delete their own invites or invites for clubs they admin
	if i.Email == user.Email {
		return nil
	}

	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", i.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: can only delete your own invites or invites for clubs you admin")
	}

	return nil
}
