package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type Fine struct {
	ID        string  `json:"id" gorm:"type:uuid;primary_key"`
	ClubID    string  `json:"club_id" gorm:"type:uuid"`
	UserID    string  `json:"userId" gorm:"type:uuid"`
	Reason    string  `json:"reason"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Paid      bool    `json:"paid"`
}

func (c *Club) CreateFine(userID, reason string, amount float64) (Fine, error) {

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
	fine.CreatedAt = time.Now().Format(time.RFC3339)
	fine.UpdatedAt = time.Now().Format(time.RFC3339)

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
