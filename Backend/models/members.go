package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/database"
	"github.com/NLstn/civo/notifications"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrLastOwnerDemotion = errors.New("cannot demote the last owner of the club")
var ErrMemberAlreadyExists = errors.New("member already exists for club and user")

type Member struct {
	ID        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid;uniqueIndex:idx_members_club_user" odata:"required"`
	UserID    string    `json:"UserID" gorm:"type:uuid;uniqueIndex:idx_members_club_user" odata:"required"`
	Role      string    `json:"Role" gorm:"default:member" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	UpdatedBy string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`

	// Navigation properties for OData
	User            *User                  `gorm:"foreignKey:UserID" json:"User,omitempty" odata:"nav"`
	Club            *Club                  `gorm:"foreignKey:ClubID" json:"Club,omitempty" odata:"nav"`
	PrivacySettings *MemberPrivacySettings `gorm:"foreignKey:MemberID" json:"PrivacySettings,omitempty" odata:"nav"`
}

func (c *Club) IsOwner(user User) bool {
	role, err := c.GetMemberRole(user)
	if err != nil {
		return false
	}
	if role == "owner" {
		return true
	}
	return false
}

func (c *Club) IsAdmin(user User) bool {
	role, err := c.GetMemberRole(user)
	if err != nil {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	return false
}

// CountOwners returns the number of owners in the club
func (c *Club) CountOwners() (int64, error) {
	var count int64
	err := database.Db.Model(&Member{}).Where("club_id = ? AND role = ?", c.ID, "owner").Count(&count).Error
	return count, err
}

// IsMember reports whether the provided user belongs to the club. If the
// user has an empty ID the function simply returns false.
func (c *Club) IsMember(user User) bool {
	if user.ID == "" {
		return false
	}

	result := database.Db.Where("club_id = ? AND user_id = ?", c.ID, user.ID).Limit(1).Find(&Member{})
	if result.Error != nil {
		return false
	}
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func (c *Club) GetClubMembers() ([]Member, error) {
	var members []Member
	err := database.Db.Where("club_id = ?", c.ID).Find(&members).Error
	return members, err
}

func (c *Club) AddMember(userId, role string) error {
	return c.addMemberWithActor(userId, role, true, nil)
}

func (c *Club) AddMemberWithActor(userId, role, actorID string) error {
	return c.addMemberWithActor(userId, role, true, &actorID)
}

func (c *Club) AddMemberViaInvite(userId, role string) error {
	return c.addMemberWithActor(userId, role, false, nil)
}

func (c *Club) addMemberWithActor(userId, role string, sendNotification bool, actorID *string) error {
	var member Member
	member.ID = uuid.New().String()
	member.ClubID = c.ID
	member.UserID = userId
	member.Role = role
	// For now, set created_by to the user being added since we don't have the adding user's ID
	member.CreatedBy = userId
	member.UpdatedBy = userId
	err := database.Db.Create(&member).Error
	if err != nil {
		// Check if this is a unique constraint violation
		errMsg := err.Error()
		if strings.Contains(errMsg, "UNIQUE constraint failed") ||
			strings.Contains(errMsg, "duplicate key") ||
			strings.Contains(errMsg, "unique constraint") {
			return ErrMemberAlreadyExists
		}
		return err
	}
	if sendNotification {
		member.notifyAdded(actorID)
	}
	return nil
}

func (c *Club) DeleteMember(memberID string) (int64, error) {
	result := database.Db.Where("id = ? AND club_id = ?", memberID, c.ID).Delete(&Member{})
	return result.RowsAffected, result.Error
}

func (c *Club) DeleteMemberByUserID(userID string) error {
	result := database.Db.Where("user_id = ? AND club_id = ?", userID, c.ID).Delete(&Member{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (c *Club) GetMemberRole(user User) (string, error) {
	var member Member
	result := database.Db.Where("club_id = ? AND user_id = ?", c.ID, user.ID).First(&member)
	if result.Error != nil {
		return "", result.Error
	}
	return member.Role, nil
}

func (c *Club) UpdateMemberRole(changingUser User, memberID, role string) error {
	var member Member
	result := database.Db.Where("id = ? AND club_id = ?", memberID, c.ID).First(&member)
	if result.Error != nil {
		return result.Error
	}
	if role != "owner" && role != "admin" && role != "member" {
		return gorm.ErrInvalidData
	}

	if canChange, err := c.canChangeRole(changingUser, member, role); err != nil {
		if err == ErrLastOwnerDemotion {
			return err
		}
		return gorm.ErrInvalidData
	} else if !canChange {
		return gorm.ErrInvalidData
	}

	// Store the old role before updating
	oldRole := member.Role
	member.Role = role
	member.UpdatedBy = changingUser.ID

	err := database.Db.Save(&member).Error
	if err != nil {
		return err
	}

	// Send notification if role actually changed
	if oldRole != role {
		member.notifyRoleChanged(oldRole, role, c.Name, changingUser.ID)
	}

	return nil
}

func (m *Member) notifyAdded(actorID *string) {
	var club Club
	if err := database.Db.Where("id = ?", m.ClubID).First(&club).Error; err != nil {
		return
	}

	var user User
	if err := database.Db.Where("id = ?", m.UserID).First(&user).Error; err != nil {
		return
	}

	// Create activity entry for timeline so admins can see new member joins
	err := CreateMemberJoinedActivity(club.ID, user.ID, club.Name, actorID)
	if err != nil {
		// Log error but don't fail the operation
		log.Printf("Failed to create member joined activity for user %s in club %s: %v", user.ID, club.ID, err)
	}

	// Send in-app notification based on preferences
	SendMemberAddedNotifications(user.ID, user.Email, club.ID, club.Name)

	// Also send the traditional email notification for backward compatibility
	notifications.SendMemberAddedNotification(user.Email, club.ID, club.Name)
}

func (m *Member) notifyRoleChanged(oldRole, newRole, clubName, actorID string) {
	// Create activity entry for the role change
	err := CreateRoleChangeActivity(m.ClubID, m.UserID, actorID, clubName, oldRole, newRole)
	if err != nil {
		// Log error but don't fail the operation
		log.Printf("Failed to create role change activity for user %s in club %s: %v", m.UserID, m.ClubID, err)
	}

	// Get user notification preferences
	preferences, err := GetUserNotificationPreferences(m.UserID)
	if err != nil {
		// If preferences don't exist, create default ones and continue
		preferences, err = CreateDefaultUserNotificationPreferences(m.UserID)
		if err != nil {
			return
		}
	}

	// Send in-app notification based on preferences
	SendRoleChangedNotifications(m.UserID, m.ClubID, clubName, oldRole, newRole)

	// Send email notification if enabled in preferences
	if preferences.RoleChangedEmail {
		var user User
		if err := database.Db.Where("id = ?", m.UserID).First(&user).Error; err != nil {
			return
		}
		notifications.SendRoleChangedNotification(user.Email, m.ClubID, clubName, oldRole, newRole)
	}
}

func (c *Club) canChangeRole(changingUser User, targetMember Member, newRole string) (bool, error) {
	changingUserRole, err := c.GetMemberRole(changingUser)
	if err != nil {
		return false, err
	}

	// Check if the changing user is trying to demote themselves as the last owner
	if changingUser.ID == targetMember.UserID && targetMember.Role == "owner" && newRole != "owner" {
		ownerCount, err := c.CountOwners()
		if err != nil {
			return false, err
		}
		if ownerCount <= 1 {
			return false, ErrLastOwnerDemotion
		}
	}

	// Only owners can change any role
	if changingUserRole == "owner" {
		return true, nil
	}

	// SECURITY FIX: Admins can only promote members to admin or demote admins to member
	// They CANNOT create/promote to owner or demote owners
	if changingUserRole == "admin" {
		// Admins cannot touch owner roles
		if targetMember.Role == "owner" || newRole == "owner" {
			return false, nil
		}
		// Admins can change between member and admin roles
		if (targetMember.Role == "member" || targetMember.Role == "admin") &&
			(newRole == "member" || newRole == "admin") {
			return true, nil
		}
	}

	// Members cannot change roles
	// This explicit check makes the authorization logic clearer
	if changingUserRole == "member" {
		return false, nil
	}

	return false, nil
}

// BeforeCreate GORM hook - sets UUID and timestamps if not provided
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	// Set timestamps if not already set
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}

	return nil
}

// ODataBeforeReadCollection filters members to only those in clubs the user belongs to
func (m Member) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see members of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific member record
func (m Member) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see members of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates member creation permissions
func (m *Member) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get the creating user's role
	var creatingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ?", m.ClubID, userID).First(&creatingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only club members can add other members")
	}

	// SECURITY: Only owners can add new owners, only owners and admins can add admins
	if m.Role == "owner" && creatingMember.Role != "owner" {
		return fmt.Errorf("unauthorized: only owners can add new owners")
	}
	if m.Role == "admin" && creatingMember.Role != "owner" && creatingMember.Role != "admin" {
		return fmt.Errorf("unauthorized: only admins and owners can add new admins")
	}

	// Check if user is an admin/owner of the club
	if creatingMember.Role != "admin" && creatingMember.Role != "owner" {
		return fmt.Errorf("unauthorized: only admins and owners can add members")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.CreatedBy = userID
	m.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates member update permissions
func (m *Member) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get the current member state from database before update
	var currentMember Member
	if err := database.Db.Where("id = ?", m.ID).First(&currentMember).Error; err != nil {
		return fmt.Errorf("member not found")
	}

	// SECURITY: Prevent changing the club of an existing member (ClubID is immutable)
	if m.ClubID != currentMember.ClubID {
		return fmt.Errorf("forbidden: club cannot be changed for an existing member")
	}

	// SECURITY: Check if role is being changed
	if m.Role != currentMember.Role {
		// Role changes require special authorization
		club := Club{ID: currentMember.ClubID}
		changingUser := User{ID: userID}

		// Use the same authorization logic as UpdateMemberRole
		// Note: We use currentMember (from DB) to avoid TOCTOU race conditions
		canChange, err := club.canChangeRole(changingUser, currentMember, m.Role)
		if err != nil {
			if err == ErrLastOwnerDemotion {
				return fmt.Errorf("cannot demote the last owner of the club")
			}
			return fmt.Errorf("unauthorized: cannot change member role")
		}
		if !canChange {
			return fmt.Errorf("unauthorized: insufficient permissions to change member role")
		}
	}

	// Check if user is an admin/owner of the club for other updates
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", currentMember.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update members")
	}

	// Set UpdatedBy
	now := time.Now()
	m.UpdatedAt = now
	m.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates member deletion permissions
func (m *Member) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club, or is deleting themselves
	if m.UserID == userID {
		// Users can leave clubs (delete their own membership)
		return nil
	}

	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", m.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can remove members")
	}

	return nil
}
