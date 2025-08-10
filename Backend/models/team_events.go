package models

import (
	"time"

	"github.com/NLstn/clubs/database"
)

// TeamEvent represents an event that belongs to a team within a club
// It mirrors the Event model but references a team instead of a club
// and uses TeamEventRSVP for RSVP tracking.
type TeamEvent struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeamID      string    `gorm:"type:uuid;not null" json:"team_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Location    string    `gorm:"type:varchar(255)" json:"location"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by" gorm:"type:uuid"`
}

// TeamEventRSVP stores RSVP responses for team events
// mirroring EventRSVP but for TeamEvent.
type TeamEventRSVP struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	EventID   string    `gorm:"type:uuid;not null" json:"event_id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	Response  string    `gorm:"not null" json:"response"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Event TeamEvent `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User  User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CreateEvent creates a new event for the team
func (t *Team) CreateEvent(name, description, location string, startTime, endTime time.Time, createdBy string) (*TeamEvent, error) {
	event := TeamEvent{
		TeamID:      t.ID,
		Name:        name,
		Description: description,
		Location:    location,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}

	if err := database.Db.Create(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEvents returns all events for the team
func (t *Team) GetEvents() ([]TeamEvent, error) {
	var events []TeamEvent
	err := database.Db.Where("team_id = ?", t.ID).Order("start_time ASC").Find(&events).Error
	return events, err
}

// GetUpcomingEvents returns upcoming events for the team
func (t *Team) GetUpcomingEvents() ([]TeamEvent, error) {
	var events []TeamEvent
	now := time.Now()
	err := database.Db.Where("team_id = ? AND start_time >= ?", t.ID, now).Order("start_time ASC").Find(&events).Error
	return events, err
}

// UpdateEvent updates an existing team event
func (t *Team) UpdateEvent(eventID, name, description, location string, startTime, endTime time.Time, updatedBy string) (*TeamEvent, error) {
	var event TeamEvent
	if err := database.Db.Where("id = ? AND team_id = ?", eventID, t.ID).First(&event).Error; err != nil {
		return nil, err
	}

	event.Name = name
	event.Description = description
	event.Location = location
	event.StartTime = startTime
	event.EndTime = endTime
	event.UpdatedBy = updatedBy

	if err := database.Db.Save(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// DeleteEvent deletes a team event and its RSVPs
func (t *Team) DeleteEvent(eventID string) error {
	if err := database.Db.Where("event_id = ?", eventID).Delete(&TeamEventRSVP{}).Error; err != nil {
		return err
	}
	return database.Db.Where("id = ? AND team_id = ?", eventID, t.ID).Delete(&TeamEvent{}).Error
}

// GetEventByID returns a team event by ID
func (t *Team) GetEventByID(eventID string) (*TeamEvent, error) {
	var event TeamEvent
	if err := database.Db.Where("id = ? AND team_id = ?", eventID, t.ID).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// CreateOrUpdateTeamEventRSVP creates or updates an RSVP for a team event
func (u *User) CreateOrUpdateTeamEventRSVP(eventID, response string) error {
	var rsvp TeamEventRSVP
	err := database.Db.Where("event_id = ? AND user_id = ?", eventID, u.ID).First(&rsvp).Error
	if err != nil {
		rsvp = TeamEventRSVP{EventID: eventID, UserID: u.ID, Response: response}
		return database.Db.Create(&rsvp).Error
	}
	rsvp.Response = response
	return database.Db.Save(&rsvp).Error
}

// GetTeamEventRSVPs returns all RSVPs for a team event
func GetTeamEventRSVPs(eventID string) ([]TeamEventRSVP, error) {
	var rsvps []TeamEventRSVP
	err := database.Db.Where("event_id = ?", eventID).Preload("User").Find(&rsvps).Error
	return rsvps, err
}

// GetTeamEventRSVPCounts returns RSVP counts for a team event
func GetTeamEventRSVPCounts(eventID string) (map[string]int, error) {
	var results []struct {
		Response string
		Count    int
	}

	err := database.Db.Table("team_event_rsvps").
		Select("response, COUNT(*) as count").
		Where("event_id = ?", eventID).
		Group("response").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Response] = r.Count
	}
	return counts, nil
}

// GetUserTeamEventRSVP returns a user's RSVP for a team event
func (u *User) GetUserTeamEventRSVP(eventID string) (*TeamEventRSVP, error) {
	var rsvp TeamEventRSVP
	if err := database.Db.Where("event_id = ? AND user_id = ?", eventID, u.ID).First(&rsvp).Error; err != nil {
		return nil, err
	}
	return &rsvp, nil
}
