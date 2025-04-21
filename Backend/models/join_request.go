package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type JoinRequest struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
	Email  string `json:"email"`
}

func CreateJoinRequest(clubId, email string) error {
	request := &JoinRequest{
		ID:     uuid.New().String(),
		ClubID: clubId,
		Email:  email,
	}
	return database.Db.Create(request).Error
}

func GetJoinRequestsForClub(clubId string) ([]JoinRequest, error) {
	var requests []JoinRequest
	err := database.Db.Where("club_id = ?", clubId).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func GetUserJoinRequests(userId string) ([]JoinRequest, error) {
	var requests []JoinRequest
	err := database.Db.Where("email = (SELECT email FROM users WHERE id = ?)", userId).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func GetUserCanEditJoinRequest(userId, requestId string) (bool, error) {
	var user User
	err := database.Db.Where("id = ?", userId).First(&user).Error
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

	err = club.AddMember(userId)
	if err != nil {
		return err
	}

	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}

func RejectJoinRequest(requestId string) error {
	// FIXME: inform admins about rejection?
	return database.Db.Delete(&JoinRequest{}, "id = ?", requestId).Error
}
