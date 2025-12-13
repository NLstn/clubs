package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"gorm.io/gorm"
)

type Event struct {
	ID          string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	ClubID      string    `json:"ClubID" gorm:"type:uuid;not null" odata:"required"`
	TeamID      *string   `json:"TeamID,omitempty" gorm:"type:uuid" odata:"nullable"` // Optional team association
	Name        string    `json:"Name" gorm:"not null" odata:"required"`
	Description *string   `json:"Description,omitempty" gorm:"type:text" odata:"nullable"`
	Location    *string   `json:"Location,omitempty" gorm:"type:varchar(255)" odata:"nullable"`
	StartTime   time.Time `json:"StartTime" gorm:"not null" odata:"required"`
	EndTime     time.Time `json:"EndTime" gorm:"not null" odata:"required"`
	CreatedAt   time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy   string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	UpdatedBy   string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
	// Recurring event fields
	IsRecurring        bool       `json:"IsRecurring" gorm:"column:is_recurring;default:false"`
	RecurrencePattern  *string    `json:"RecurrencePattern,omitempty" gorm:"column:recurrence_pattern;type:varchar(50)" odata:"nullable"` // "weekly", "daily", "monthly"
	RecurrenceInterval int        `json:"RecurrenceInterval,omitempty" gorm:"column:recurrence_interval;default:1"`                       // Every N weeks/days/months
	RecurrenceEnd      *time.Time `json:"RecurrenceEnd,omitempty" gorm:"column:recurrence_end" odata:"nullable"`                          // When recurrence stops
	ParentEventID      *string    `json:"ParentEventID,omitempty" gorm:"column:parent_event_id;type:uuid" odata:"nullable"`               // Links recurring event instances
}

type EventRSVP struct {
	ID        string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	EventID   string    `json:"EventID" gorm:"type:uuid;not null" odata:"required"`
	UserID    string    `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
	Response  string    `json:"Response" gorm:"not null" odata:"required"` // "yes" or "no"
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	UpdatedAt time.Time `json:"UpdatedAt"`

	// Relationships
	Event Event `gorm:"foreignKey:EventID" json:"Event,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"User,omitempty"`
}

// CreateEvent creates a new event for the club
func (c *Club) CreateEvent(name string, description string, location string, startTime, endTime time.Time, createdBy string) (*Event, error) {
	event := Event{
		ClubID:      c.ID,
		Name:        name,
		Description: &description,
		Location:    &location,
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
		ClubID:             c.ID,
		Name:               name,
		Description:        &description,
		Location:           &location,
		StartTime:          currentStart,
		EndTime:            currentEnd,
		CreatedBy:          createdBy,
		UpdatedBy:          createdBy,
		IsRecurring:        true,
		RecurrencePattern:  &recurrencePattern,
		RecurrenceInterval: recurrenceInterval,
		RecurrenceEnd:      &recurrenceEnd,
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
			Description:   &description,
			Location:      &location,
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
	event.Description = &description
	event.Location = &location
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

// GetUserRSVPsByEventIDs returns a map of event IDs to RSVPs for the user
func (u *User) GetUserRSVPsByEventIDs(eventIDs []string) (map[string]*EventRSVP, error) {
	if len(eventIDs) == 0 {
		return make(map[string]*EventRSVP), nil
	}

	var rsvps []EventRSVP
	err := database.Db.Where("event_id IN ? AND user_id = ?", eventIDs, u.ID).Find(&rsvps).Error
	if err != nil {
		return nil, err
	}

	// Build map for efficient lookup
	rsvpMap := make(map[string]*EventRSVP)
	for i := range rsvps {
		rsvpMap[rsvps[i].EventID] = &rsvps[i]
	}

	return rsvpMap, nil
}

// EventWithClub represents an event with its associated club name
type EventWithClub struct {
	Event
	ClubName string `json:"ClubName"`
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
		Description: &description,
		Location:    &location,
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
	event.Description = &description
	event.Location = &location
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

// ODataBeforeReadCollection filters events to only those in clubs the user belongs to
func (e Event) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see events of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific event
func (e Event) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see events of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates event creation permissions
func (e *Event) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", e.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create events")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	e.CreatedBy = userID
	e.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates event update permissions
func (e *Event) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", e.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update events")
	}

	// Set UpdatedBy
	now := time.Now()
	e.UpdatedAt = now
	e.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates event deletion permissions
func (e *Event) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", e.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can delete events")
	}

	return nil
}

// EventRSVP authorization hooks
// ODataBeforeReadCollection filters RSVPs to only those for events the user can access
func (er EventRSVP) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see RSVPs for events in clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("event_id IN (SELECT id FROM events WHERE club_id IN (SELECT club_id FROM members WHERE user_id = ?))", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific RSVP
func (er EventRSVP) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see RSVPs for events in clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("event_id IN (SELECT id FROM events WHERE club_id IN (SELECT club_id FROM members WHERE user_id = ?))", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates RSVP creation permissions
func (er *EventRSVP) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get event to check club membership
	var event Event
	if err := database.Db.Where("id = ?", er.EventID).First(&event).Error; err != nil {
		return fmt.Errorf("event not found")
	}

	// Check if user is a member of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ?", event.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only club members can RSVP to events")
	}

	// Users can only create RSVPs for themselves
	if er.UserID != "" && er.UserID != userID {
		return fmt.Errorf("unauthorized: cannot create RSVP for another user")
	}

	// Set UserID if not already set
	if er.UserID == "" {
		er.UserID = userID
	}

	// Set CreatedAt and UpdatedAt
	now := time.Now()
	er.CreatedAt = now
	er.UpdatedAt = now

	return nil
}

// ODataBeforeUpdate validates RSVP update permissions
func (er *EventRSVP) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Users can only update their own RSVPs
	if er.UserID != userID {
		return fmt.Errorf("unauthorized: can only update your own RSVPs")
	}

	// Set UpdatedAt
	er.UpdatedAt = time.Now()

	return nil
}

// ODataBeforeDelete validates RSVP deletion permissions
func (er *EventRSVP) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Users can only delete their own RSVPs
	if er.UserID != userID {
		// Check if user is club admin/owner
		var event Event
		if err := database.Db.Where("id = ?", er.EventID).First(&event).Error; err != nil {
			return fmt.Errorf("event not found")
		}

		var existingMember Member
		if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", event.ClubID, userID).First(&existingMember).Error; err != nil {
			return fmt.Errorf("unauthorized: can only delete your own RSVPs")
		}
	}

	return nil
}
