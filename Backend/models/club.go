package models

import (
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Club struct {
	ID          string     `json:"id" gorm:"type:uuid;primary_key" odata:"key"`
	Name        string     `json:"name" odata:"required"`
	Description *string    `json:"description" odata:"nullable"`
	LogoURL     *string    `json:"logo_url,omitempty" odata:"nullable"`
	CreatedAt   time.Time  `json:"created_at" odata:"immutable"`
	CreatedBy   string     `json:"created_by" gorm:"type:uuid" odata:"required"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   string     `json:"updated_by" gorm:"type:uuid" odata:"required"`
	Deleted     bool       `json:"deleted" gorm:"default:false"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" odata:"nullable"`
	DeletedBy   *string    `json:"deleted_by,omitempty" gorm:"type:uuid" odata:"nullable"`

	// Navigation properties for OData
	Members []Member `gorm:"foreignKey:ClubID" json:"members,omitempty" odata:"nav"`
	Teams   []Team   `gorm:"foreignKey:ClubID" json:"teams,omitempty" odata:"nav"`
	Events  []Event  `gorm:"foreignKey:ClubID" json:"events,omitempty" odata:"nav"`
	News    []News   `gorm:"foreignKey:ClubID" json:"news,omitempty" odata:"nav"`
	Fines   []Fine   `gorm:"foreignKey:ClubID" json:"fines,omitempty" odata:"nav"`
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
