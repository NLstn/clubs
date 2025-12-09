package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type Fine struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID    string    `json:"club_id" gorm:"type:uuid" odata:"required"`
	TeamID    *string   `json:"team_id,omitempty" gorm:"type:uuid" odata:"nullable"` // Optional team association
	UserID    string    `json:"userId" gorm:"type:uuid" odata:"required"`
	Reason    string    `json:"reason" odata:"required"`
	Amount    float64   `json:"amount" odata:"required"`
	CreatedAt time.Time `json:"created_at" odata:"immutable"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid" odata:"required"`
	Paid      bool      `json:"paid"`
}

func (c *Club) CreateFine(userID, reason, createdBy string, amount float64) (Fine, error) {

	user, err := GetUserByID(userID)
	if err != nil {
		return Fine{}, err
	}
	if !c.IsMember(user) {
		err = fmt.Errorf("user is not a member of the club")
		return Fine{}, err
	}

	var fine Fine
	fine.ID = uuid.New().String()
	fine.ClubID = c.ID
	fine.UserID = userID
	fine.Reason = reason
	fine.Amount = amount
	fine.CreatedBy = createdBy
	fine.UpdatedBy = createdBy

	err = database.Db.Create(&fine).Error
	if err != nil {
		return Fine{}, err
	}

	return fine, nil
}

func (c *Club) GetFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Where("club_id = ?", c.ID).Find(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}

func (c *Club) DeleteFine(fineID string) error {
	return database.Db.Where("id = ? AND club_id = ?", fineID, c.ID).Delete(&Fine{}).Error
}

func (t *Team) CreateFine(userID, reason, createdBy string, amount float64) (Fine, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return Fine{}, err
	}

	// Check if user is a member of the team
	if !t.IsMember(user) {
		err = fmt.Errorf("user is not a member of the team")
		return Fine{}, err
	}

	var fine Fine
	fine.ID = uuid.New().String()
	fine.ClubID = t.ClubID
	fine.TeamID = &t.ID
	fine.UserID = userID
	fine.Reason = reason
	fine.Amount = amount
	fine.CreatedBy = createdBy
	fine.UpdatedBy = createdBy

	err = database.Db.Create(&fine).Error
	if err != nil {
		return Fine{}, err
	}

	return fine, nil
}

func (t *Team) GetFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Where("team_id = ?", t.ID).Find(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}

func (t *Team) DeleteFine(fineID string) error {
	return database.Db.Where("id = ? AND team_id = ?", fineID, t.ID).Delete(&Fine{}).Error
}
