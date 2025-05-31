package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type Fine struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	ClubID    string    `json:"club_id" gorm:"type:uuid"`
	UserID    string    `json:"userId" gorm:"type:uuid"`
	Reason    string    `json:"reason"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
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
