package core

import (
	"fmt"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
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

func GetFineTemplateByID(templateID string) (FineTemplate, error) {
	var template FineTemplate
	err := database.Db.Where("id = ?", templateID).First(&template).Error
	if err != nil {
		return FineTemplate{}, err
	}
	return template, nil
}

func (c *Club) UpdateFineTemplate(templateID, description string, amount float64, updatedBy string) (FineTemplate, error) {
	template, err := GetFineTemplateByID(templateID)
	if err != nil {
		return FineTemplate{}, err
	}

	if template.ClubID != c.ID {
		return FineTemplate{}, fmt.Errorf("template does not belong to this club")
	}

	template.Description = description
	template.Amount = amount
	template.UpdatedBy = updatedBy

	err = database.Db.Save(&template).Error
	if err != nil {
		return FineTemplate{}, err
	}

	return template, nil
}

func (c *Club) DeleteFineTemplate(templateID, deletedBy string) error {
	template, err := GetFineTemplateByID(templateID)
	if err != nil {
		return err
	}

	if template.ClubID != c.ID {
		return fmt.Errorf("template does not belong to this club")
	}

	err = database.Db.Delete(&template).Error
	return err
}
