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

type UserShiftDetails struct {
	ID         string    `json:"id"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	EventID    string    `json:"eventId"`
	EventName  string    `json:"eventName"`
	Location   string    `json:"location"`
	ClubID     string    `json:"clubId"`
	ClubName   string    `json:"clubName"`
	Members    []string  `json:"members"` // Array of member names
}

func GetUserFutureShifts(userID string) ([]UserShiftDetails, error) {
	var shifts []UserShiftDetails
	
	query := `
		SELECT DISTINCT 
			s.id, s.start_time, s.end_time, s.event_id,
			e.name as event_name, e.location, e.club_id,
			c.name as club_name
		FROM shifts s
		INNER JOIN shift_members sm ON s.id = sm.shift_id  
		INNER JOIN events e ON s.event_id = e.id
		INNER JOIN clubs c ON e.club_id = c.id
		WHERE sm.user_id = ? AND s.start_time > NOW()
		ORDER BY s.start_time ASC
	`
	
	err := database.Db.Raw(query, userID).Scan(&shifts).Error
	if err != nil {
		return nil, err
	}
	
	// For each shift, get all the members assigned to it
	for i := range shifts {
		shiftMembers, err := GetShiftMembers(shifts[i].ID)
		if err != nil {
			return nil, err
		}
		
		var memberNames []string
		for _, shiftMember := range shiftMembers {
			user, err := GetUserByID(shiftMember.UserID)
			if err != nil {
				return nil, err
			}
			memberNames = append(memberNames, user.GetFullName())
		}
		shifts[i].Members = memberNames
	}
	
	return shifts, nil
}
