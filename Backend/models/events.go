package models

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
)

type Event struct {
	ID                string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClubID            string     `gorm:"type:uuid;not null" json:"club_id"`
	TeamID            *string    `gorm:"type:uuid" json:"team_id,omitempty"` // Optional team association
	Name              string     `gorm:"not null" json:"name"`
	Description       string     `gorm:"type:text" json:"description"`
	Location          string     `gorm:"type:varchar(255)" json:"location"`
	StartTime         time.Time  `gorm:"not null" json:"start_time"`
	EndTime           time.Time  `gorm:"not null" json:"end_time"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         string     `json:"created_by" gorm:"type:uuid"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         string     `json:"updated_by" gorm:"type:uuid"`
	// Recurring event fields
	IsRecurring       bool       `gorm:"default:false" json:"is_recurring"`
	RecurrencePattern string     `gorm:"type:varchar(50)" json:"recurrence_pattern,omitempty"` // "weekly", "daily", "monthly"
	RecurrenceInterval int       `gorm:"default:1" json:"recurrence_interval,omitempty"`        // Every N weeks/days/months
	RecurrenceEnd     *time.Time `json:"recurrence_end,omitempty"`                             // When recurrence stops
	ParentEventID     *string    `gorm:"type:uuid" json:"parent_event_id,omitempty"`           // Links recurring event instances
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
func (c *Club) CreateEvent(name string, description string, location string, startTime, endTime time.Time, createdBy string) (*Event, error) {
	event := Event{
		ClubID:      c.ID,
		Name:        name,
		Description: description,
		Location:    location,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}

	tx := database.Db.Create(&event)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &event, nil
}

// CreateRecurringEvent creates recurring events based on the recurrence pattern
func (c *Club) CreateRecurringEvent(name string, description string, location string, startTime, endTime time.Time, 
	recurrencePattern string, recurrenceInterval int, recurrenceEnd time.Time, createdBy string) ([]*Event, error) {
	
	if recurrencePattern == "" || recurrenceInterval < 1 {
		return nil, fmt.Errorf("invalid recurrence parameters")
	}

	var events []*Event
	currentStart := startTime
	currentEnd := endTime
	duration := endTime.Sub(startTime)

	// Create parent event (first occurrence)
	parentEvent := Event{
		ClubID:            c.ID,
		Name:              name,
		Description:       description,
		Location:          location,
		StartTime:         currentStart,
		EndTime:           currentEnd,
		CreatedBy:         createdBy,
		UpdatedBy:         createdBy,
		IsRecurring:       true,
		RecurrencePattern: recurrencePattern,
		RecurrenceInterval: recurrenceInterval,
		RecurrenceEnd:     &recurrenceEnd,
	}

	tx := database.Db.Create(&parentEvent)
	if tx.Error != nil {
		return nil, tx.Error
	}

	events = append(events, &parentEvent)

	// Generate recurring instances
	for {
		// Calculate next occurrence
		switch recurrencePattern {
		case "daily":
			currentStart = currentStart.AddDate(0, 0, recurrenceInterval)
		case "weekly":
			currentStart = currentStart.AddDate(0, 0, 7*recurrenceInterval)
		case "monthly":
			currentStart = currentStart.AddDate(0, recurrenceInterval, 0)
		default:
			return events, fmt.Errorf("unsupported recurrence pattern: %s", recurrencePattern)
		}

		currentEnd = currentStart.Add(duration)

		// Stop if we've passed the end date
		if currentStart.After(recurrenceEnd) {
			break
		}

		// Create recurring instance
		recurringEvent := Event{
			ClubID:        c.ID,
			Name:          name,
			Description:   description,
			Location:      location,
			StartTime:     currentStart,
			EndTime:       currentEnd,
			CreatedBy:     createdBy,
			UpdatedBy:     createdBy,
			IsRecurring:   false, // Individual instances are not marked as recurring
			ParentEventID: &parentEvent.ID,
		}

		tx := database.Db.Create(&recurringEvent)
		if tx.Error != nil {
			return events, tx.Error
		}

		events = append(events, &recurringEvent)
	}

	return events, nil
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
func (c *Club) UpdateEvent(eventID string, name string, description string, location string, startTime, endTime time.Time, updatedBy string) (*Event, error) {
	var event Event
	err := database.Db.Where("id = ? AND club_id = ?", eventID, c.ID).First(&event).Error
	if err != nil {
		return nil, err
	}

	event.Name = name
	event.Description = description
	event.Location = location
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
			EventID:  eventID,
			UserID:   u.ID,
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

// EventWithClub represents an event with its associated club name
type EventWithClub struct {
	Event
	ClubName string `json:"club_name"`
}

// SearchEventsForUser searches for events in clubs where the user is a member
func SearchEventsForUser(userID, query string) ([]EventWithClub, error) {
	var events []EventWithClub

	// Query events from clubs where user is a member and event name contains the query
	err := database.Db.Table("events e").
		Select("e.*, c.name as club_name").
		Joins("JOIN clubs c ON e.club_id = c.id").
		Joins("JOIN members m ON m.club_id = c.id").
		Where("m.user_id = ? AND LOWER(e.name) LIKE LOWER(?)", userID, "%"+query+"%").
		Where("c.deleted = false"). // Only from non-deleted clubs
		Order("e.start_time DESC").
		Scan(&events).Error

	return events, err
}

// CreateEvent creates a new event for the team
func (t *Team) CreateEvent(name string, description string, location string, startTime, endTime time.Time, createdBy string) (*Event, error) {
	event := Event{
		ClubID:      t.ClubID,
		TeamID:      &t.ID,
		Name:        name,
		Description: description,
		Location:    location,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}

	tx := database.Db.Create(&event)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &event, nil
}

// GetEvents returns all events for the team
func (t *Team) GetEvents() ([]Event, error) {
	var events []Event
	err := database.Db.Where("team_id = ?", t.ID).Order("start_time ASC").Find(&events).Error
	return events, err
}

// GetUpcomingEvents returns upcoming events for the team
func (t *Team) GetUpcomingEvents() ([]Event, error) {
	var events []Event
	now := time.Now()
	err := database.Db.Where("team_id = ? AND start_time >= ?", t.ID, now).
		Order("start_time ASC").Find(&events).Error
	return events, err
}

// UpdateEvent updates an existing team event
func (t *Team) UpdateEvent(eventID string, name string, description string, location string, startTime, endTime time.Time, updatedBy string) (*Event, error) {
	var event Event
	err := database.Db.Where("id = ? AND team_id = ?", eventID, t.ID).First(&event).Error
	if err != nil {
		return nil, err
	}

	event.Name = name
	event.Description = description
	event.Location = location
	event.StartTime = startTime
	event.EndTime = endTime
	event.UpdatedBy = updatedBy

	err = database.Db.Save(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// DeleteEvent deletes a team event
func (t *Team) DeleteEvent(eventID string) error {
	// First delete all RSVPs for this event
	err := database.Db.Where("event_id = ?", eventID).Delete(&EventRSVP{}).Error
	if err != nil {
		return err
	}

	// Then delete the event
	return database.Db.Where("id = ? AND team_id = ?", eventID, t.ID).Delete(&Event{}).Error
}

// GetEventByID returns a team event by ID
func (t *Team) GetEventByID(eventID string) (*Event, error) {
	var event Event
	err := database.Db.Where("id = ? AND team_id = ?", eventID, t.ID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}
