package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Club struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by" gorm:"type:uuid"`
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

func GetClubsByIDs(clubIDs []string) ([]Club, error) {
	var clubs []Club
	if len(clubIDs) == 0 {
		return clubs, nil
	}
	err := database.Db.Where("id IN ?", clubIDs).Find(&clubs).Error
	return clubs, err
}

func CreateClub(name, description, ownerID string) (Club, error) {
	var club Club
	club.ID = uuid.New().String()
	club.Name = name
	club.Description = description
	club.CreatedBy = ownerID
	club.UpdatedBy = ownerID

	err := database.Db.Transaction(func(tx *gorm.DB) error {
		if dbErr := tx.Create(&club).Error; dbErr != nil {
			return dbErr
		}
		var member Member
		member.ID = uuid.New().String()
		member.ClubID = club.ID
		member.UserID = ownerID
		member.Role = "owner"
		member.CreatedBy = ownerID
		member.UpdatedBy = ownerID
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

func (c *Club) Update(name, description, updatedBy string) error {
	return database.Db.Model(c).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
		"updated_by":  updatedBy,
	}).Error
}
