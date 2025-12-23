package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type News struct {
	ID        string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
	ClubID    string    `json:"ClubID" gorm:"type:uuid;not null" odata:"required"`
	Title     string    `json:"Title" gorm:"not null" odata:"required"`
	Content   string    `json:"Content" gorm:"type:text;not null" odata:"required"`
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	UpdatedBy string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`
}

// EntitySetName returns the custom entity set name for the News entity.
// By default, "News" would be pluralized to "Newses", but we want to keep it as "News".
func (News) EntitySetName() string {
	return "News"
}

// CreateNews creates a new news post for the club
func (c *Club) CreateNews(title, content, createdBy string) (*News, error) {
	news := News{
		ID:        uuid.New().String(),
		ClubID:    c.ID,
		Title:     title,
		Content:   content,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}

	tx := database.Db.Create(&news)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &news, nil
}

// GetNews returns all news posts for the club
func (c *Club) GetNews() ([]News, error) {
	var news []News
	err := database.Db.Where("club_id = ?", c.ID).Order("created_at DESC").Find(&news).Error
	return news, err
}

// UpdateNews updates an existing news post
func (c *Club) UpdateNews(newsID string, title, content, updatedBy string) (*News, error) {
	var news News
	err := database.Db.Where("id = ? AND club_id = ?", newsID, c.ID).First(&news).Error
	if err != nil {
		return nil, err
	}

	news.Title = title
	news.Content = content
	news.UpdatedBy = updatedBy

	err = database.Db.Save(&news).Error
	if err != nil {
		return nil, err
	}

	return &news, nil
}

// DeleteNews deletes a news post
func (c *Club) DeleteNews(newsID string) error {
	return database.Db.Where("id = ? AND club_id = ?", newsID, c.ID).Delete(&News{}).Error
}

// ODataBeforeReadCollection filters news to only those in clubs the user belongs to
func (n News) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see news of clubs they belong to and where news feature is enabled
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?) AND club_id IN (SELECT club_id FROM club_settings WHERE news_enabled = true)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific news post
func (n News) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see news of clubs they belong to and where news feature is enabled
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?) AND club_id IN (SELECT club_id FROM club_settings WHERE news_enabled = true)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates news creation permissions
func (n *News) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if news feature is enabled for the club
	if err := CheckFeatureEnabled(n.ClubID, "news"); err != nil {
		return err
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", n.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create news")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	n.CreatedBy = userID
	n.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates news update permissions
func (n *News) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if news feature is enabled for the club
	if err := CheckFeatureEnabled(n.ClubID, "news"); err != nil {
		return err
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", n.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update news")
	}

	// Set UpdatedBy
	now := time.Now()
	n.UpdatedAt = now
	n.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates news deletion permissions
func (n *News) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if news feature is enabled for the club
	if err := CheckFeatureEnabled(n.ClubID, "news"); err != nil {
		return err
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", n.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can delete news")
	}

	return nil
}
