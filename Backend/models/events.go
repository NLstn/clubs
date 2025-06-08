package models

import (
	"time"

	"github.com/NLstn/clubs/database"
)

type Event struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClubID    string    `gorm:"type:uuid;not null" json:"club_id"`
	Name      string    `gorm:"not null" json:"name"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
}

type EventRSVP struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	EventID   string    `gorm:"type:uuid;not null" json:"event_id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	Response  string    `gorm:"not null" json:"response"` // "yes" or "no"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Event Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CreateEvent creates a new event for the club
func (c *Club) CreateEvent(name string, startTime, endTime time.Time, createdBy string) (*Event, error) {
	event := Event{
		ClubID:    c.ID,
		Name:      name,
		StartTime: startTime,
		EndTime:   endTime,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}

	tx := database.Db.Create(&event)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &event, nil
}

// GetEvents returns all events for the club
func (c *Club) GetEvents() ([]Event, error) {
	var events []Event
	err := database.Db.Where("club_id = ?", c.ID).Order("start_time ASC").Find(&events).Error
	return events, err
}

// GetUpcomingEvents returns upcoming events for the club
func (c *Club) GetUpcomingEvents() ([]Event, error) {
	var events []Event
	now := time.Now()
	err := database.Db.Where("club_id = ? AND start_time >= ?", c.ID, now).
		Order("start_time ASC").Find(&events).Error
	return events, err
}

// UpdateEvent updates an existing event
func (c *Club) UpdateEvent(eventID string, name string, startTime, endTime time.Time, updatedBy string) (*Event, error) {
	var event Event
	err := database.Db.Where("id = ? AND club_id = ?", eventID, c.ID).First(&event).Error
	if err != nil {
		return nil, err
	}

	event.Name = name
	event.StartTime = startTime
	event.EndTime = endTime
	event.UpdatedBy = updatedBy

	err = database.Db.Save(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// DeleteEvent deletes an event
func (c *Club) DeleteEvent(eventID string) error {
	// First delete all RSVPs for this event
	err := database.Db.Where("event_id = ?", eventID).Delete(&EventRSVP{}).Error
	if err != nil {
		return err
	}

	// Then delete the event
	return database.Db.Where("id = ? AND club_id = ?", eventID, c.ID).Delete(&Event{}).Error
}

// GetEventByID returns an event by ID
func (c *Club) GetEventByID(eventID string) (*Event, error) {
	var event Event
	err := database.Db.Where("id = ? AND club_id = ?", eventID, c.ID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// CreateOrUpdateRSVP creates or updates an RSVP for an event
func (u *User) CreateOrUpdateRSVP(eventID string, response string) error {
	var rsvp EventRSVP
	err := database.Db.Where("event_id = ? AND user_id = ?", eventID, u.ID).First(&rsvp).Error
	
	if err != nil {
		// Create new RSVP
		rsvp = EventRSVP{
			EventID: eventID,
			UserID:  u.ID,
			Response: response,
		}
		return database.Db.Create(&rsvp).Error
	} else {
		// Update existing RSVP
		rsvp.Response = response
		return database.Db.Save(&rsvp).Error
	}
}

// GetEventRSVPs returns all RSVPs for an event
func GetEventRSVPs(eventID string) ([]EventRSVP, error) {
	var rsvps []EventRSVP
	err := database.Db.Where("event_id = ?", eventID).Preload("User").Find(&rsvps).Error
	return rsvps, err
}

// GetEventRSVPCounts returns RSVP counts for an event
func GetEventRSVPCounts(eventID string) (map[string]int, error) {
	var results []struct {
		Response string
		Count    int
	}
	
	err := database.Db.Table("event_rsvps").
		Select("response, COUNT(*) as count").
		Where("event_id = ?", eventID).
		Group("response").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	counts := make(map[string]int)
	for _, result := range results {
		counts[result.Response] = result.Count
	}
	
	return counts, nil
}

// GetUserRSVP returns a user's RSVP for an event
func (u *User) GetUserRSVP(eventID string) (*EventRSVP, error) {
	var rsvp EventRSVP
	err := database.Db.Where("event_id = ? AND user_id = ?", eventID, u.ID).First(&rsvp).Error
	if err != nil {
		return nil, err
	}
	return &rsvp, nil
}