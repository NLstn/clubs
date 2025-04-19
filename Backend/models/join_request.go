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

func GetJoinRequests(clubId string) ([]JoinRequest, error) {
	var requests []JoinRequest
	err := database.Db.Where("club_id = ?", clubId).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}
