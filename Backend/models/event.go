package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type Event struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClubID      string `json:"club_id" gorm:"type:uuid"`
	Date        string `json:"date"`
	BeginTime   string `json:"begin_time"`
	EndTime     string `json:"end_time"`
}

func GetClubEvents(clubID string) ([]Event, error) {
	var events []Event
	err := database.Db.Where("club_id = ?", clubID).Find(&events).Error
	return events, err
}

func CreateEvent(event *Event, clubID string) error {
	event.ID = uuid.New().String()
	event.ClubID = clubID
	return database.Db.Create(event).Error
}

func DeleteEvent(eventID, clubID string) (int64, error) {
	result := database.Db.Where("id = ? AND club_id = ?", eventID, clubID).Delete(&Event{})
	return result.RowsAffected, result.Error
}

func (e *Event) Validate() bool {
	return e.Name != "" && e.Date != "" && e.BeginTime != "" && e.EndTime != ""
}
