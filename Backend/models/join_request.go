package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JoinRequest represents a user request to join a club (via invitation link)
type JoinRequest struct {
	ID        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid" odata:"required"`
	UserID    string    `json:"UserID" gorm:"type:uuid" odata:"required"`
	Email     string    `json:"Email" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	UpdatedAt time.Time `json:"UpdatedAt"`
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
	err = c.notifyAdminsAboutJoinRequest(userID, email, request.ID)
	if err != nil {
		// Log error but don't fail the operation
		log.Printf("ERROR: Failed to notify admins about join request %s for club %s: %v", request.ID, c.ID, err)
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

	// Remove notifications for this join request
	err = RemoveJoinRequestNotifications(joinRequest.ID)
	if err != nil {
		// Log error but don't fail the operation
		log.Printf("ERROR: Failed to remove join request notifications for request %s: %v", joinRequest.ID, err)
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

	// Remove notifications for this join request
	err = RemoveJoinRequestNotifications(joinRequest.ID)
	if err != nil {
		// Log error but don't fail the operation
		log.Printf("ERROR: Failed to remove join request notifications for request %s: %v", joinRequest.ID, err)
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
func (c *Club) notifyAdminsAboutJoinRequest(userID, email, joinRequestID string) error {
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
			err := CreateNotificationWithJoinRequest(admin.ID, "join_request_received", title, message, &c.ID, nil, nil, &joinRequestID)
			if err != nil {
				// Log error but continue with other notifications
				log.Printf("ERROR: Failed to create join request notification for admin %s in club %s: %v", admin.ID, c.ID, err)
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

// ODataBeforeReadCollection filters join requests - admins see requests for their clubs, users see their own
func (jr JoinRequest) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can see their own requests OR requests for clubs they are admin/owner of
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? OR club_id IN (SELECT club_id FROM members WHERE user_id = ? AND role IN ('admin', 'owner'))", userID, userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific join request
func (jr JoinRequest) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can see their own requests OR requests for clubs they are admin/owner of
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? OR club_id IN (SELECT club_id FROM members WHERE user_id = ? AND role IN ('admin', 'owner'))", userID, userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates join request creation permissions
func (jr *JoinRequest) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Users can only create their own join requests
	if jr.UserID == "" {
		jr.UserID = userID
	} else if jr.UserID != userID {
		return fmt.Errorf("unauthorized: cannot create join request for another user")
	}

	// Set CreatedAt and UpdatedAt
	now := time.Now()
	jr.CreatedAt = now
	jr.UpdatedAt = now

	return nil
}

// ODataBeforeDelete validates join request deletion permissions
func (jr *JoinRequest) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can delete their own requests or requests for clubs they admin
	if jr.UserID == userID {
		return nil
	}

	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", jr.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: can only delete your own join requests or requests for clubs you admin")
	}

	return nil
}
