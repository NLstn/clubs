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

type FineTemplate struct {
	ID          string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID      string    `json:"ClubID" gorm:"type:uuid" odata:"required"`
	Description string    `json:"Description" odata:"required"`
	Amount      float64   `json:"Amount" odata:"required"`
	CreatedAt   time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy   string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	UpdatedBy   string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
}

func (c *Club) CreateFineTemplate(description string, amount float64, createdBy string) (FineTemplate, error) {
	var template FineTemplate
	template.ID = uuid.New().String()
	template.ClubID = c.ID
	template.Description = description
	template.Amount = amount
	template.CreatedBy = createdBy
	template.UpdatedBy = createdBy

	err := database.Db.Create(&template).Error
	if err != nil {
		return FineTemplate{}, err
	}

	return template, nil
}

func (c *Club) GetFineTemplates() ([]FineTemplate, error) {
	var templates []FineTemplate
	err := database.Db.Where("club_id = ?", c.ID).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// ODataBeforeReadCollection filters fine templates to only those in clubs the user belongs to
func (ft FineTemplate) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see fine templates of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific fine template
func (ft FineTemplate) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see fine templates of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates fine template creation permissions
func (ft *FineTemplate) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", ft.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create fine templates")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	ft.CreatedAt = now
	ft.UpdatedAt = now
	ft.CreatedBy = userID
	ft.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates fine template update permissions
func (ft *FineTemplate) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Load the existing template to enforce immutable fields
	var existingTemplate FineTemplate
	if err := database.Db.First(&existingTemplate, "id = ?", ft.ID).Error; err != nil {
		return fmt.Errorf("fine template not found")
	}

	// SECURITY: Prevent changing the club of an existing template (ClubID is immutable)
	if ft.ClubID != existingTemplate.ClubID {
		return fmt.Errorf("forbidden: club cannot be changed for an existing fine template")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", ft.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update fine templates")
	}

	// Set UpdatedBy and UpdatedAt
	ft.UpdatedAt = time.Now()
	ft.UpdatedBy = userID

	// Preserve immutable fields
	ft.CreatedAt = existingTemplate.CreatedAt
	ft.CreatedBy = existingTemplate.CreatedBy

	return nil
}

// ODataBeforeDelete validates fine template deletion permissions
func (ft *FineTemplate) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Load the existing template to check club membership
	var existingTemplate FineTemplate
	if err := database.Db.First(&existingTemplate, "id = ?", ft.ID).Error; err != nil {
		return fmt.Errorf("fine template not found")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", existingTemplate.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can delete fine templates")
	}

	return nil
}
