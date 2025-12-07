package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

// JoinRequest represents a user request to join a club (via invitation link)
type JoinRequest struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID    string    `json:"club_id" gorm:"type:uuid"`
	UserID    string    `json:"user_id" gorm:"type:uuid"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateJoinRequest creates a new join request from a user
func (c *Club) CreateJoinRequest(userID, email string) error {
	request := &JoinRequest{
		ID:     uuid.New().String(),
		ClubID: c.ID,
		UserID: userID,
		Email:  email,
	}
	err := database.Db.Create(request).Error
	if err != nil {
		return err
	}

	// Notify all admins and owners about the join request
	err = c.notifyAdminsAboutJoinRequest(userID, email)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

// GetJoinRequests returns all pending join requests for a club (admin view)
func (c *Club) GetJoinRequests() ([]JoinRequest, error) {
	var requests []JoinRequest
	err := database.Db.Where("club_id = ?", c.ID).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// AcceptJoinRequest accepts a join request and adds the user to the club
func AcceptJoinRequest(requestId, adminUserId string) error {
	var joinRequest JoinRequest
	err := database.Db.Where("id = ?", requestId).First(&joinRequest).Error
	if err != nil {
		return err
	}

	var club Club
	err = database.Db.Where("id = ?", joinRequest.ClubID).First(&club).Error
	if err != nil {
		return err
	}

	// Verify the admin has permission
	var admin User
	err = database.Db.Where("id = ?", adminUserId).First(&admin).Error
	if err != nil {
		return err
	}

	if !club.IsOwner(admin) && !club.IsAdmin(admin) {
		return fmt.Errorf("user not authorized to accept this request")
	}

	// Add user to club with the admin as the actor
	err = club.AddMemberWithActor(joinRequest.UserID, "member", adminUserId)
	if err != nil {
		return err
	}

	// Delete the join request since it's now complete
	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}

// RejectJoinRequest rejects a join request
func RejectJoinRequest(requestId, adminUserId string) error {
	var joinRequest JoinRequest
	err := database.Db.Where("id = ?", requestId).First(&joinRequest).Error
	if err != nil {
		return err
	}

	var club Club
	err = database.Db.Where("id = ?", joinRequest.ClubID).First(&club).Error
	if err != nil {
		return err
	}

	// Verify the admin has permission
	var admin User
	err = database.Db.Where("id = ?", adminUserId).First(&admin).Error
	if err != nil {
		return err
	}

	if !club.IsOwner(admin) && !club.IsAdmin(admin) {
		return fmt.Errorf("user not authorized to reject this request")
	}

	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}

// HasPendingJoinRequest checks if a user already has a pending join request for this club
func (c *Club) HasPendingJoinRequest(userID string) (bool, error) {
	var count int64
	err := database.Db.Model(&JoinRequest{}).Where("club_id = ? AND user_id = ?", c.ID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasPendingInvite checks if a user already has a pending invite for this club
func (c *Club) HasPendingInvite(email string) (bool, error) {
	var count int64
	err := database.Db.Model(&Invite{}).Where("club_id = ? AND email = ?", c.ID, email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// notifyAdminsAboutJoinRequest sends notifications to all admins and owners about a new join request
func (c *Club) notifyAdminsAboutJoinRequest(userID, email string) error {
	// Get all admins and owners
	admins, err := c.GetAdminsAndOwners()
	if err != nil {
		return fmt.Errorf("failed to get club admins: %v", err)
	}

	// Get the user who made the request
	var requestingUser User
	err = database.Db.Where("id = ?", userID).First(&requestingUser).Error
	if err != nil {
		return fmt.Errorf("failed to get requesting user: %v", err)
	}

	// Send notification to each admin/owner
	for _, admin := range admins {
		// Get user notification preferences
		preferences, err := GetUserNotificationPreferences(admin.ID)
		if err != nil {
			// If preferences don't exist, create default ones and continue
			preferences, err = CreateDefaultUserNotificationPreferences(admin.ID)
			if err != nil {
				continue // Skip this admin if we can't get/create preferences
			}
		}

		// Send in-app notification if enabled
		if preferences.JoinRequestInApp {
			title := "New Join Request"
			userName := strings.TrimSpace(requestingUser.FirstName + " " + requestingUser.LastName)
			if userName == "" {
				userName = "A user"
			}
			message := fmt.Sprintf("%s (%s) has requested to join %s", userName, email, c.Name)
			err := CreateNotification(admin.ID, "join_request_received", title, message, &c.ID, nil, nil)
			if err != nil {
				// Log error but continue with other notifications
				// TODO: Add proper logging
			}
		}

		// Send email notification if enabled and notifications package is available
		if preferences.JoinRequestEmail {
			// TODO: Implement email notification for join requests if needed
			// This would require adding a function to the notifications package
		}
	}

	return nil
}

// GetJoinRequestCount returns the number of pending join requests for a club
func (c *Club) GetJoinRequestCount() (int64, error) {
	var count int64
	err := database.Db.Model(&JoinRequest{}).Where("club_id = ?", c.ID).Count(&count).Error
	return count, err
}
