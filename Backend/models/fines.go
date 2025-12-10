package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Fine struct {
	ID        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid" odata:"required"`
	TeamID    *string   `json:"TeamID,omitempty" gorm:"type:uuid" odata:"nullable"` // Optional team association
	UserID    string    `json:"UserID" gorm:"type:uuid" odata:"required"`
	Reason    string    `json:"Reason" odata:"required"`
	Amount    float64   `json:"Amount" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	UpdatedBy string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
	Paid      bool      `json:"Paid"`
}

func (c *Club) CreateFine(userID, reason, createdBy string, amount float64) (Fine, error) {

	user, err := GetUserByID(userID)
	if err != nil {
		return Fine{}, err
	}
	if !c.IsMember(user) {
		err = fmt.Errorf("user is not a member of the club")
		return Fine{}, err
	}

	var fine Fine
	fine.ID = uuid.New().String()
	fine.ClubID = c.ID
	fine.UserID = userID
	fine.Reason = reason
	fine.Amount = amount
	fine.CreatedBy = createdBy
	fine.UpdatedBy = createdBy

	err = database.Db.Create(&fine).Error
	if err != nil {
		return Fine{}, err
	}

	return fine, nil
}

func (c *Club) GetFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Where("club_id = ?", c.ID).Find(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}

func (c *Club) DeleteFine(fineID string) error {
	return database.Db.Where("id = ? AND club_id = ?", fineID, c.ID).Delete(&Fine{}).Error
}

func (t *Team) CreateFine(userID, reason, createdBy string, amount float64) (Fine, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return Fine{}, err
	}

	// Check if user is a member of the team
	if !t.IsMember(user) {
		err = fmt.Errorf("user is not a member of the team")
		return Fine{}, err
	}

	var fine Fine
	fine.ID = uuid.New().String()
	fine.ClubID = t.ClubID
	fine.TeamID = &t.ID
	fine.UserID = userID
	fine.Reason = reason
	fine.Amount = amount
	fine.CreatedBy = createdBy
	fine.UpdatedBy = createdBy

	err = database.Db.Create(&fine).Error
	if err != nil {
		return Fine{}, err
	}

	return fine, nil
}

func (t *Team) GetFines() ([]Fine, error) {
	var fines []Fine
	err := database.Db.Where("team_id = ?", t.ID).Find(&fines).Error
	if err != nil {
		return nil, err
	}
	return fines, nil
}

func (t *Team) DeleteFine(fineID string) error {
	return database.Db.Where("id = ? AND team_id = ?", fineID, t.ID).Delete(&Fine{}).Error
}

// ODataBeforeReadCollection filters fines to only those in clubs the user belongs to
func (f Fine) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see fines of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific fine
func (f Fine) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see fines of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates fine creation permissions
func (f *Fine) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", f.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create fines")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now
	f.CreatedBy = userID
	f.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates fine update permissions
func (f *Fine) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", f.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update fines")
	}

	// Set UpdatedBy
	now := time.Now()
	f.UpdatedAt = now
	f.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates fine deletion permissions
func (f *Fine) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", f.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can delete fines")
	}

	return nil
}
