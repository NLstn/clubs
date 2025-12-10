package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/google/uuid"
	"github.com/nlstn/go-odata"
	"gorm.io/gorm"

	"github.com/NLstn/clubs/database"
)

type Club struct {
	ID          string     `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	Name        string     `json:"Name" odata:"required"`
	Description *string    `json:"Description,omitempty" odata:"nullable"`
	LogoURL     *string    `json:"LogoURL,omitempty" odata:"nullable"`
	CreatedAt   time.Time  `json:"CreatedAt" odata:"auto,immutable"`                  // Set server-side, immutable after creation
	CreatedBy   string     `json:"CreatedBy" gorm:"type:uuid" odata:"auto,immutable"` // Set server-side from context
	UpdatedAt   time.Time  `json:"UpdatedAt" odata:"auto"`                            // Set server-side automatically
	UpdatedBy   string     `json:"UpdatedBy" gorm:"type:uuid" odata:"auto"`           // Set server-side from context
	Deleted     bool       `json:"Deleted" gorm:"default:false"`
	DeletedAt   *time.Time `json:"DeletedAt,omitempty" odata:"nullable"`
	DeletedBy   *string    `json:"DeletedBy,omitempty" gorm:"type:uuid" odata:"nullable"`

	// Navigation properties for OData
	Members []Member `gorm:"foreignKey:ClubID" json:"Members,omitempty" odata:"nav"`
	Teams   []Team   `gorm:"foreignKey:ClubID" json:"Teams,omitempty" odata:"nav"`
	Events  []Event  `gorm:"foreignKey:ClubID" json:"Events,omitempty" odata:"nav"`
	News    []News   `gorm:"foreignKey:ClubID" json:"News,omitempty" odata:"nav"`
	Fines   []Fine   `gorm:"foreignKey:ClubID" json:"Fines,omitempty" odata:"nav"`
}

func GetAllClubs() ([]Club, error) {
	var clubs []Club
	err := database.Db.Where("deleted = false").Find(&clubs).Error
	return clubs, err
}

func GetAllClubsIncludingDeleted() ([]Club, error) {
	var clubs []Club
	err := database.Db.Find(&clubs).Error
	return clubs, err
}

func GetClubByID(id string) (Club, error) {
	var club Club
	result := database.Db.First(&club, "id = ?", id)
	return club, result.Error
}

func GetClubsByIDs(clubIDs []string) ([]Club, error) {
	var clubs []Club
	if len(clubIDs) == 0 {
		return clubs, nil
	}
	err := database.Db.Where("id IN ?", clubIDs).Find(&clubs).Error
	return clubs, err
}

func CreateClub(name, description, ownerID string) (Club, error) {
	var club Club
	club.ID = uuid.New().String()
	club.Name = name
	club.Description = &description
	club.CreatedBy = ownerID
	club.UpdatedBy = ownerID

	err := database.Db.Transaction(func(tx *gorm.DB) error {
		if dbErr := tx.Create(&club).Error; dbErr != nil {
			return dbErr
		}
		var member Member
		member.ID = uuid.New().String()
		member.ClubID = club.ID
		member.UserID = ownerID
		member.Role = "owner"
		member.CreatedBy = ownerID
		member.UpdatedBy = ownerID
		if dbErr := tx.Create(&member).Error; dbErr != nil {
			return dbErr
		}
		return nil
	})
	if err != nil {
		return Club{}, err
	}

	return club, nil
}

func (c *Club) Update(name, description, updatedBy string) error {
	return database.Db.Model(c).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
		"updated_by":  updatedBy,
	}).Error
}

func (c *Club) UpdateLogo(logoURL *string, updatedBy string) error {
	return database.Db.Model(c).Updates(map[string]interface{}{
		"logo_url":   logoURL,
		"updated_by": updatedBy,
	}).Error
}

func (c *Club) SoftDelete(deletedBy string) error {
	now := time.Now()
	return database.Db.Model(c).Updates(map[string]interface{}{
		"deleted":    true,
		"deleted_at": &now,
		"deleted_by": &deletedBy,
	}).Error
}

func DeleteClubPermanently(clubID string) error {
	// This will permanently delete the club and all related data
	// Note: This should cascade delete related records if foreign keys are set up properly
	return database.Db.Unscoped().Delete(&Club{}, "id = ?", clubID).Error
}

// GetAdminsAndOwners returns all users who are admins or owners of the club
func (c *Club) GetAdminsAndOwners() ([]User, error) {
	var users []User
	err := database.Db.Table("users").
		Joins("JOIN members ON users.id = members.user_id").
		Where("members.club_id = ? AND (members.role = ? OR members.role = ?)", c.ID, "admin", "owner").
		Find(&users).Error
	return users, err
}

// BeforeCreate GORM hook - sets UUID if not provided
func (c *Club) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	// Note: CreatedBy and UpdatedBy are set by OData hooks from HTTP context
	// GORM BeforeCreate runs after OData BeforeCreate

	return nil
}

// OData EntityHook implementation

// ODataBeforeCreate OData hook - sets audit fields from authenticated user context
func (c *Club) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	// Extract user ID from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Set audit fields
	now := time.Now()
	c.CreatedAt = now
	c.CreatedBy = userID
	c.UpdatedAt = now
	c.UpdatedBy = userID

	return nil
}

// AfterCreate OData hook - creates the creator as an owner member of the club
func (c *Club) AfterCreate(ctx context.Context, r *http.Request) error {
	// Get transaction from context
	tx, ok := odata.TransactionFromContext(ctx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	// Create member record with owner role
	member := Member{
		ClubID:    c.ID,
		UserID:    c.CreatedBy,
		Role:      "owner",
		CreatedBy: c.CreatedBy,
		UpdatedBy: c.CreatedBy,
	}

	if err := tx.Create(&member).Error; err != nil {
		return fmt.Errorf("failed to create owner member: %w", err)
	}

	return nil
}

// BeforeUpdate OData hook - sets UpdatedBy from authenticated user context
func (c *Club) BeforeUpdate(ctx context.Context, r *http.Request) error {
	// Extract user ID from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Set updated by and updated at
	c.UpdatedAt = time.Now()
	c.UpdatedBy = userID

	return nil
}
