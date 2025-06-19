package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type JoinRequest struct {
	ID            string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID        string    `json:"club_id" gorm:"type:uuid"`
	Email         string    `json:"email"`
	AdminApproved bool      `json:"admin_approved" gorm:"default:false"`
	UserApproved  bool      `json:"user_approved" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by" gorm:"type:uuid"`
}

func (c *Club) CreateJoinRequest(email, createdBy string, adminApproved, userApproved bool) error {
	request := &JoinRequest{
		ID:            uuid.New().String(),
		ClubID:        c.ID,
		Email:         email,
		AdminApproved: adminApproved,
		UserApproved:  userApproved,
		CreatedBy:     createdBy,
		UpdatedBy:     createdBy,
	}
	return database.Db.Create(request).Error
}

func (c *Club) GetJoinRequests() ([]JoinRequest, error) {
	var requests []JoinRequest
	// Admins should see requests where admin approval is needed (AdminApproved=false)
	err := database.Db.Where("club_id = ? AND admin_approved = ?", c.ID, false).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (u *User) GetJoinRequests() ([]JoinRequest, error) {
	var requests []JoinRequest
	// Users should see invites where user approval is needed (UserApproved=false)
	err := database.Db.Where("email = (SELECT email FROM users WHERE id = ?) AND user_approved = ?", u.ID, false).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (u *User) GetUserCanEditJoinRequest(requestId string) (bool, error) {
	var user User
	err := database.Db.Where("id = ?", u.ID).First(&user).Error
	if err != nil {
		return false, err
	}

	var request JoinRequest
	err = database.Db.Where("id = ?", requestId).First(&request).Error
	if err != nil {
		return false, err
	}

	// User can accept if it's their own invite (their email matches)
	if user.Email == request.Email {
		return true, nil
	}

	// Admin can accept if it's a user request to join their club
	var club Club
	err = database.Db.Where("id = ?", request.ClubID).First(&club).Error
	if err != nil {
		return false, err
	}

	if club.IsOwner(user) || club.IsAdmin(user) {
		return true, nil
	}

	return false, nil
}

func AcceptJoinRequest(requestId, userId string) error {
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

	// Get the user who is accepting the request
	var acceptingUser User
	err = database.Db.Where("id = ?", userId).First(&acceptingUser).Error
	if err != nil {
		return err
	}

	// Determine who should be added to the club based on the request email
	var targetUser User
	err = database.Db.Where("email = ?", joinRequest.Email).First(&targetUser).Error
	if err != nil {
		return err
	}

	// Update approval status based on who is accepting
	if acceptingUser.Email == joinRequest.Email {
		// User is accepting their own invite (admin invited them)
		joinRequest.UserApproved = true
	} else if club.IsOwner(acceptingUser) || club.IsAdmin(acceptingUser) {
		// Admin is accepting a user's request to join
		joinRequest.AdminApproved = true
	} else {
		return fmt.Errorf("user not authorized to accept this request")
	}

	// If both approvals are now true, add the target user to the club
	if joinRequest.AdminApproved && joinRequest.UserApproved {
		err = club.AddMember(targetUser.ID, "member")
		if err != nil {
			return err
		}
		// Delete the join request since it's now complete
		return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
	} else {
		// Update the request with the new approval status
		return database.Db.Save(&joinRequest).Error
	}
}

func RejectJoinRequest(requestId string) error {
	// FIXME: inform admins about rejection?
	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}
