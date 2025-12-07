package models

import (
	"errors"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/notifications"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrLastOwnerDemotion = errors.New("cannot demote the last owner of the club")

type Member struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID    string    `json:"club_id" gorm:"type:uuid"`
	UserID    string    `json:"user_id" gorm:"type:uuid"`
	Role      string    `json:"role" gorm:"default:member"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
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
	return c.addMember(userId, role, true)
}

func (c *Club) AddMemberViaInvite(userId, role string) error {
	return c.addMember(userId, role, false)
}

func (c *Club) addMember(userId, role string, sendNotification bool) error {
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
		return err
	}
	if sendNotification {
		member.notifyAdded()
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

func (m *Member) notifyAdded() {
	var club Club
	if err := database.Db.Where("id = ?", m.ClubID).First(&club).Error; err != nil {
		return
	}

	var user User
	if err := database.Db.Where("id = ?", m.UserID).First(&user).Error; err != nil {
		return
	}

	// Create activity entry for timeline so admins can see new member joins
	err := CreateMemberJoinedActivity(club.ID, user.ID, club.Name)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
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
		// TODO: Add proper logging
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

	if changingUserRole == "owner" {
		return true, nil
	}
	if changingUserRole == "admin" && (targetMember.Role == "member" || newRole == "admin") {
		return true, nil
	}
	return false, nil
}
