package models

import (
	"time"

	"github.com/NLstn/clubs/database"
)

type Shift struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClubID    string    `gorm:"type:uuid;not null"`
	EventID   string    `gorm:"type:uuid;not null" json:"eventId"`
	StartTime time.Time `gorm:"not null" json:"startTime"`
	EndTime   time.Time `gorm:"not null" json:"endTime"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
}

type ShiftMember struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShiftID   string    `gorm:"type:uuid;not null"`
	UserID    string    `gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
}

func (c *Club) CreateShift(startTime, endTime time.Time, createdBy string, eventID string) (string, error) {
	shift := Shift{
		ClubID:    c.ID,
		EventID:   eventID,
		StartTime: startTime,
		EndTime:   endTime,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}

	tx := database.Db.Create(&shift)
	if tx.Error != nil {
		return "", tx.Error
	}

	return shift.ID, nil
}

func AddMemberToShift(shiftID, userID, createdBy string) error {
	shiftMember := ShiftMember{
		ShiftID:   shiftID,
		UserID:    userID,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}

	tx := database.Db.Create(&shiftMember)
	return tx.Error
}

func (c *Club) GetShifts() ([]Shift, error) {
	var shifts []Shift
	tx := database.Db.Model(&Shift{}).Where("club_id = ?", c.ID).Find(&shifts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return shifts, nil
}

func (c *Club) GetShiftsByEvent(eventID string) ([]Shift, error) {
	var shifts []Shift
	tx := database.Db.Model(&Shift{}).Where("club_id = ? AND event_id = ?", c.ID, eventID).Find(&shifts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return shifts, nil
}

func GetShiftMembers(shiftID string) ([]ShiftMember, error) {
	var shiftMembers []ShiftMember
	err := database.Db.Where("shift_id = ?", shiftID).Find(&shiftMembers).Error
	if err != nil {
		return nil, err
	}

	return shiftMembers, nil
}

func RemoveMemberFromShift(shiftID, userID string) error {
	tx := database.Db.Where("shift_id = ? AND user_id = ?", shiftID, userID).Delete(&ShiftMember{})
	return tx.Error
}
