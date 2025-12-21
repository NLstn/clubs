package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invite represents an admin invitation to a user to join a club
type Invite struct {
	ID        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid;uniqueIndex:idx_club_email" odata:"required"`
	Email     string    `json:"Email" gorm:"uniqueIndex:idx_club_email" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// CreateInvite creates a new invitation from an admin to a user
func (c *Club) CreateInvite(email, createdBy string) error {
	// Check if invite already exists for this club and email
	var existingInvite Invite
	result := database.Db.Where("club_id = ? AND email = ?", c.ID, email).First(&existingInvite)
	if result.Error == nil {
		return fmt.Errorf("invite already exists for this email")
	}
	// Check for database errors other than "not found"
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing invite: %w", result.Error)
	}

	// Check if user with this email is already a member
	var existingMember Member
	result = database.Db.Joins("JOIN users ON users.id = members.user_id").
		Where("members.club_id = ? AND users.email = ?", c.ID, email).
		First(&existingMember)
	if result.Error == nil {
		return fmt.Errorf("user is already a member of this club")
	}
	// Check for database errors other than "not found"
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing member: %w", result.Error)
	}

	// Check rate limit for admin (10 invites per hour)
	if err := checkAdminInviteRateLimit(c.ID, createdBy); err != nil {
		return err
	}

	// Check rate limit for club (50 invites per hour)
	if err := checkClubInviteRateLimit(c.ID); err != nil {
		return err
	}

	invite := &Invite{
		ID:        uuid.New().String(),
		ClubID:    c.ID,
		Email:     email,
		CreatedBy: createdBy,
	}
	err := database.Db.Create(invite).Error
	if err != nil {
		// Handle unique constraint violation gracefully
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || 
		   strings.Contains(err.Error(), "duplicate key") ||
		   strings.Contains(err.Error(), "unique constraint") {
			return fmt.Errorf("invite already exists for this email")
		}
		return fmt.Errorf("failed to create invite: %w", err)
	}

	// Send invite notification to the user
	SendInviteReceivedNotifications(email, c.ID, c.Name, invite.ID)

	return nil
}

// checkAdminInviteRateLimit checks if an admin has exceeded the invite rate limit
// Limit: 10 invites per hour per admin
func checkAdminInviteRateLimit(clubID, adminID string) error {
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	var count int64
	err := database.Db.Model(&Invite{}).
		Where("club_id = ? AND created_by = ? AND created_at > ?",
			clubID, adminID, oneHourAgo).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check admin rate limit: %w", err)
	}

	if count >= 10 {
		return fmt.Errorf("rate limit exceeded: maximum 10 invites per hour per admin")
	}

	return nil
}

// checkClubInviteRateLimit checks if a club has exceeded the invite rate limit
// Limit: 50 invites per hour per club
func checkClubInviteRateLimit(clubID string) error {
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	var count int64
	err := database.Db.Model(&Invite{}).
		Where("club_id = ? AND created_at > ?", clubID, oneHourAgo).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check club rate limit: %w", err)
	}

	if count >= 50 {
		return fmt.Errorf("rate limit exceeded: maximum 50 invites per hour per club")
	}

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
