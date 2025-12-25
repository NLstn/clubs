package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/civo/auth"
	"github.com/google/uuid"
	"github.com/nlstn/go-odata"
	"gorm.io/gorm"

	"github.com/NLstn/civo/database"
)

// DefaultMaxActiveClubs is the default maximum number of active (non-deleted) clubs a user can create.
// This limit can be increased in the future through a paid subscription.
const DefaultMaxActiveClubs = 3

// ErrClubLimitExceeded returns an error when a user tries to create more clubs than their quota allows
func ErrClubLimitExceeded() error {
	return fmt.Errorf("club creation limit exceeded: you can only create up to %d active clubs", DefaultMaxActiveClubs)
}

type Club struct {
	ID          string     `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	Name        string     `json:"Name" odata:"required"`
	Description *string    `json:"Description,omitempty" odata:"nullable"`
	LogoURL     *string    `json:"LogoURL,omitempty" odata:"nullable"`
	CreatedAt   time.Time  `json:"CreatedAt" odata:"auto,immutable"`                                                         // Set server-side, immutable after creation
	CreatedBy   string     `json:"CreatedBy" gorm:"type:uuid;index:idx_clubs_created_by_deleted" odata:"auto,immutable"`     // Set server-side from context
	UpdatedAt   time.Time  `json:"UpdatedAt" odata:"auto"`                                                                   // Set server-side automatically
	UpdatedBy   string     `json:"UpdatedBy" gorm:"type:uuid" odata:"auto"`                                                  // Set server-side from context
	Deleted     bool       `json:"Deleted" gorm:"default:false;index:idx_clubs_created_by_deleted"`
	DeletedAt   *time.Time `json:"DeletedAt,omitempty" odata:"nullable"`
	DeletedBy   *string    `json:"DeletedBy,omitempty" gorm:"type:uuid" odata:"nullable"`

	// Navigation properties for OData
	Members  []Member      `gorm:"foreignKey:ClubID" json:"Members,omitempty" odata:"nav"`
	Teams    []Team        `gorm:"foreignKey:ClubID" json:"Teams,omitempty" odata:"nav"`
	Events   []Event       `gorm:"foreignKey:ClubID" json:"Events,omitempty" odata:"nav"`
	News     []News        `gorm:"foreignKey:ClubID" json:"News,omitempty" odata:"nav"`
	Fines    []Fine        `gorm:"foreignKey:ClubID" json:"Fines,omitempty" odata:"nav"`
	Settings *ClubSettings `gorm:"foreignKey:ClubID" json:"Settings,omitempty" odata:"nav"`
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

// GetAdminsAndOwners returns all users who are admins or owners of the club
func (c *Club) GetAdminsAndOwners() ([]User, error) {
	var users []User
	err := database.Db.Table("users").
		Joins("JOIN members ON users.id = members.user_id").
		Where("members.club_id = ? AND (members.role = ? OR members.role = ?)", c.ID, "admin", "owner").
		Find(&users).Error
	return users, err
}

// CountActiveClubsCreatedByUser returns the number of active (non-deleted) clubs created by a user.
// This is used to enforce the club creation quota.
// If tx is provided, it uses that transaction; otherwise it uses the global database connection.
func CountActiveClubsCreatedByUser(userID string, tx *gorm.DB) (int64, error) {
	var count int64
	db := tx
	if db == nil {
		db = database.Db
	}
	err := db.Model(&Club{}).
		Where("created_by = ? AND deleted = ?", userID, false).
		Count(&count).Error
	return count, err
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
// and enforces the club creation quota within the same transaction to prevent race conditions
func (c *Club) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	// Extract user ID from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get transaction from context for atomic quota check
	tx, _ := odata.TransactionFromContext(ctx)

	// Check club creation quota within the transaction to prevent race conditions
	activeClubCount, err := CountActiveClubsCreatedByUser(userID, tx)
	if err != nil {
		return fmt.Errorf("failed to check club creation quota: %w", err)
	}
	if activeClubCount >= DefaultMaxActiveClubs {
		return ErrClubLimitExceeded()
	}

	// Set audit fields
	now := time.Now()
	c.CreatedAt = now
	c.CreatedBy = userID
	c.UpdatedAt = now
	c.UpdatedBy = userID

	return nil
}

// ODataAfterCreate OData hook - creates the creator as an owner member of the club and default settings
func (c *Club) ODataAfterCreate(ctx context.Context, r *http.Request) error {
	// Get transaction from context
	tx, ok := odata.TransactionFromContext(ctx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	// Create member record with owner role
	now := time.Now()
	member := Member{
		ID:        uuid.New().String(),
		ClubID:    c.ID,
		UserID:    c.CreatedBy,
		Role:      "owner",
		CreatedAt: now,
		CreatedBy: c.CreatedBy,
		UpdatedAt: now,
		UpdatedBy: c.CreatedBy,
	}

	if err := tx.Create(&member).Error; err != nil {
		return fmt.Errorf("failed to create owner member: %w", err)
	}

	// Create default club settings with all features disabled
	if err := createClubSettingsWithTransaction(tx, c.ID, c.CreatedBy, now); err != nil {
		return fmt.Errorf("failed to create default club settings: %w", err)
	}

	return nil
}

// ODataBeforeUpdate OData hook - sets UpdatedBy from authenticated user context
func (c *Club) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	// Extract user ID from context
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Set updated by and updated at (these will always be updated)
	now := time.Now()
	c.UpdatedAt = now
	c.UpdatedBy = userID

	return nil
}

// ODataAfterUpdate OData hook - handles soft delete timestamp setting
func (c *Club) ODataAfterUpdate(ctx context.Context, r *http.Request) error {
	// If the club was just marked as deleted, set the soft delete fields
	if c.Deleted && c.DeletedAt == nil {
		// Get transaction from context
		tx, ok := odata.TransactionFromContext(ctx)
		if !ok {
			return nil // Log but don't fail the update
		}

		// Get the user ID from context
		userID, ok := ctx.Value(auth.UserIDKey).(string)
		if !ok || userID == "" {
			return nil // Log but don't fail
		}

		// Update the soft delete fields
		now := time.Now()
		if err := tx.Model(&Club{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
			"deleted_at": now,
			"deleted_by": userID,
		}).Error; err != nil {
			return fmt.Errorf("failed to set soft delete timestamp: %w", err)
		}

		// Update the in-memory struct as well
		c.DeletedAt = &now
		c.DeletedBy = &userID
	}

	return nil
}

// ODataBeforeReadCollection OData read hook - filters clubs based on membership and discoverability
// Users can see:
// 1. Non-deleted clubs they are members of
// 2. Non-deleted clubs where DiscoverableByNonMembers is enabled
// Note: Deleted clubs are only visible to owners when includeDeleted=true (handled elsewhere)
func (c Club) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	scope := func(db *gorm.DB) *gorm.DB {
		// Show only non-deleted clubs where:
		// 1. User is a member
		// OR
		// 2. Club is discoverable by non-members
		return db.Where(
			"clubs.deleted = ? AND (clubs.id IN (SELECT club_id FROM members WHERE user_id = ?) OR clubs.id IN (SELECT club_id FROM club_settings WHERE discoverable_by_non_members = ?))",
			false,
			userID,
			true,
		)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity OData read hook - allows access to clubs based on membership and discoverability
// Users can access:
// 1. Non-deleted clubs they are members of
// 2. Non-deleted clubs where DiscoverableByNonMembers is enabled
// Note: Deleted clubs are only accessible to owners when requested with includeDeleted=true
func (c Club) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	scope := func(db *gorm.DB) *gorm.DB {
		// Allow access to non-deleted clubs only where:
		// 1. User is a member
		// OR
		// 2. Club is discoverable by non-members
		return db.Where(
			"clubs.deleted = ? AND (clubs.id IN (SELECT club_id FROM members WHERE user_id = ?) OR clubs.id IN (SELECT club_id FROM club_settings WHERE discoverable_by_non_members = ?))",
			false,
			userID,
			true,
		)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}
