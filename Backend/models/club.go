package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type Club struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id" gorm:"type:uuid"`
}

func CreateClub(club *Club, ownerID string) error {
	club.ID = uuid.New().String()
	club.OwnerID = ownerID
	return database.Db.Create(club).Error
}

func GetAllClubs() ([]Club, error) {
	var clubs []Club
	err := database.Db.Find(&clubs).Error
	return clubs, err
}

func GetClubByID(id string) (Club, error) {
	var club Club
	result := database.Db.First(&club, "id = ?", id)
	return club, result.Error
}

func (c *Club) IsOwner(userID string) bool {
	return c.OwnerID == userID
}
