package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
)

type News struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClubID    string    `gorm:"type:uuid;not null" json:"club_id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by" gorm:"type:uuid"`
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
