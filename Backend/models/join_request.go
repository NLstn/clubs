package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type JoinRequest struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID    string    `json:"club_id" gorm:"type:uuid"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
}

func (c *Club) CreateJoinRequest(email, createdBy string) error {
	request := &JoinRequest{
		ID:        uuid.New().String(),
		ClubID:    c.ID,
		Email:     email,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
	return database.Db.Create(request).Error
}

func (c *Club) GetJoinRequests() ([]JoinRequest, error) {
	var requests []JoinRequest
	err := database.Db.Where("club_id = ?", c.ID).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (u *User) GetJoinRequests() ([]JoinRequest, error) {
	var requests []JoinRequest
	err := database.Db.Where("email = (SELECT email FROM users WHERE id = ?)", u.ID).Find(&requests).Error
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

	if user.Email != request.Email {
		return false, nil
	}

	return true, nil
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

	err = club.AddMember(userId, "member")
	if err != nil {
		return err
	}

	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}

func RejectJoinRequest(requestId string) error {
	// FIXME: inform admins about rejection?
	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}
