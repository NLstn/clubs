package models

type Club struct {
	ID          string   `json:"id" gorm:"type:uuid;primary_key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []Member `json:"members,omitempty" gorm:"foreignKey:ClubID"`
}

type Member struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
}
