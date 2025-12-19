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

type Shift struct {
	ID        string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid;not null" odata:"required"`
	EventID   string    `json:"EventID" gorm:"type:uuid;not null" odata:"required"`
	StartTime time.Time `json:"StartTime" gorm:"not null" odata:"required"`
	EndTime   time.Time `json:"EndTime" gorm:"not null" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	UpdatedBy string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
}

type ShiftMember struct {
	ID        string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	ShiftID   string    `json:"ShiftID" gorm:"type:uuid;not null" odata:"required"`
	UserID    string    `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	UpdatedBy string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
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
	ID        string    `json:"ID"`
	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`
	EventID   string    `json:"EventID"`
	EventName string    `json:"EventName"`
	Location  string    `json:"Location"`
	ClubID    string    `json:"ClubID"`
	ClubName  string    `json:"ClubName"`
	Members   []string  `json:"Members"` // Array of member names
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

// ODataBeforeReadCollection filters shifts to only those in clubs the user belongs to
func (s Shift) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see shifts of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific shift
func (s Shift) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see shifts of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates shift creation permissions
func (s *Shift) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// SECURITY: Verify the EventID belongs to the specified ClubID
	// EventID is a required field, so we just check it's not empty
	if s.EventID == "" {
		return fmt.Errorf("event ID is required")
	}
	
	var event Event
	if err := database.Db.Where("id = ? AND club_id = ?", s.EventID, s.ClubID).First(&event).Error; err != nil {
		return fmt.Errorf("unauthorized: event does not belong to the specified club")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", s.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create shifts")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	s.CreatedBy = userID
	s.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates shift update permissions
func (s *Shift) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Load the existing shift to enforce immutable fields
	var existingShift Shift
	if err := database.Db.First(&existingShift, "id = ?", s.ID).Error; err != nil {
		return fmt.Errorf("shift not found")
	}

	// SECURITY: Prevent changing the club of an existing shift (ClubID is immutable)
	if s.ClubID != existingShift.ClubID {
		return fmt.Errorf("forbidden: club cannot be changed for an existing shift")
	}

	// SECURITY: Verify the EventID belongs to the (unchanged) ClubID
	// EventID is a required field, so we just check it's not empty
	if s.EventID == "" {
		return fmt.Errorf("event ID is required")
	}
	
	var event Event
	if err := database.Db.Where("id = ? AND club_id = ?", s.EventID, existingShift.ClubID).First(&event).Error; err != nil {
		return fmt.Errorf("unauthorized: event does not belong to the specified club")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", existingShift.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update shifts")
	}

	// Set UpdatedBy
	now := time.Now()
	s.UpdatedAt = now
	s.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates shift deletion permissions
func (s *Shift) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", s.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can delete shifts")
	}

	return nil
}

// ShiftMember authorization hooks
// ODataBeforeReadCollection filters shift members to only those in clubs the user belongs to
func (sm ShiftMember) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see shift members of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("shift_id IN (SELECT id FROM shifts WHERE club_id IN (SELECT club_id FROM members WHERE user_id = ?))", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific shift member record
func (sm ShiftMember) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see shift members of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("shift_id IN (SELECT id FROM shifts WHERE club_id IN (SELECT club_id FROM members WHERE user_id = ?))", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates shift member creation permissions
func (sm *ShiftMember) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get shift to find club ID
	var shift Shift
	if err := database.Db.Where("id = ?", sm.ShiftID).First(&shift).Error; err != nil {
		return fmt.Errorf("shift not found")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", shift.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can add shift members")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	sm.CreatedAt = now
	sm.UpdatedAt = now
	sm.CreatedBy = userID
	sm.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates shift member update permissions
func (sm *ShiftMember) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get shift to find club ID
	var shift Shift
	if err := database.Db.Where("id = ?", sm.ShiftID).First(&shift).Error; err != nil {
		return fmt.Errorf("shift not found")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", shift.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update shift members")
	}

	// Set UpdatedBy
	now := time.Now()
	sm.UpdatedAt = now
	sm.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates shift member deletion permissions
func (sm *ShiftMember) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get shift to find club ID
	var shift Shift
	if err := database.Db.Where("id = ?", sm.ShiftID).First(&shift).Error; err != nil {
		return fmt.Errorf("shift not found")
	}

	// Check if user is an admin/owner of the club, or removing themselves
	if sm.UserID == userID {
		return nil
	}

	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", shift.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can remove shift members")
	}

	return nil
}
