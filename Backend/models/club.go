package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Club struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func CreateClub(name, description, ownerID string) (Club, error) {
	var club Club
	club.ID = uuid.New().String()
	club.Name = name
	club.Description = description

	err := database.Db.Transaction(func(tx *gorm.DB) error {
		if dbErr := tx.Create(&club).Error; dbErr != nil {
			return dbErr
		}
		var member Member
		member.ID = uuid.New().String()
		member.ClubID = club.ID
		member.UserID = ownerID
		member.Role = "owner"
		if dbErr := tx.Create(&member).Error; dbErr != nil {
			return dbErr
		}
		return nil
	})
	if err != nil {
		return Club{}, err
	}

	return club, nil
}

func (c *Club) Update(name, description string) error {
	return database.Db.Model(c).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
	}).Error
}
