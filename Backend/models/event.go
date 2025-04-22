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

func (c *Club) GetEvents() ([]Event, error) {
	var events []Event
	err := database.Db.Where("club_id = ?", c.ID).Find(&events).Error
	return events, err
}

func (c *Club) CreateEvent(name, description, date, beginTime, endTime string) (Event, error) {
	var event Event
	event.ID = uuid.New().String()
	event.Name = name
	event.Description = description
	event.ClubID = c.ID
	event.Date = date
	event.BeginTime = beginTime
	event.EndTime = endTime

	err := database.Db.Create(event).Error
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

func (c *Club) DeleteEvent(eventID string) (int64, error) {
	result := database.Db.Where("id = ? AND club_id = ?", eventID, c.ID).Delete(&Event{})
	return result.RowsAffected, result.Error
}
