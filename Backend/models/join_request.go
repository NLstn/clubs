package models

import (
	"fmt"
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
	return database.Db.Create(request).Error
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

	// Add user to club
	err = club.AddMember(joinRequest.UserID, "member")
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
