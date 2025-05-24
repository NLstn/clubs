package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/NLstn/clubs/database"
)

// CustomTime handles both full RFC3339 format and the shortened format from frontend
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements custom JSON unmarshaling for time
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	timeStr := strings.Trim(string(data), "\"")

	// Try parsing different formats
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // "2006-01-02T15:04:05"
		"2006-01-02T15:04",    // "2006-01-02T15:04" (frontend format)
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			ct.Time = t
			return nil
		}
	}

	return fmt.Errorf("unable to parse time: %s", timeStr)
}

// MarshalJSON implements custom JSON marshaling for time
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format(time.RFC3339))
}

// Value implements the driver.Valuer interface for database storage
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	case []byte:
		return ct.Time.UnmarshalText(v)
	case string:
		return ct.Time.UnmarshalText([]byte(v))
	default:
		return fmt.Errorf("cannot scan %T into CustomTime", value)
	}
}

type Shift struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClubID    string     `gorm:"type:uuid;not null"`
	StartTime CustomTime `gorm:"not null" json:"startTime"`
	EndTime   CustomTime `gorm:"not null" json:"endTime"`
}

type ShiftMember struct {
	ID      string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShiftID string `gorm:"type:uuid;not null"`
	UserID  string `gorm:"type:uuid;not null"`
}

func (c *Club) CreateShift(startTime, endTime time.Time) (string, error) {
	shift := Shift{
		ClubID:    c.ID,
		StartTime: CustomTime{startTime},
		EndTime:   CustomTime{endTime},
	}

	tx := database.Db.Create(&shift)
	if tx.Error != nil {
		return "", tx.Error
	}

	return shift.ID, nil
}

func AddMemberToShift(shiftID, userID string) error {
	shiftMember := ShiftMember{
		ShiftID: shiftID,
		UserID:  userID,
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
