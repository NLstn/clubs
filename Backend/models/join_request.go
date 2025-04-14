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
