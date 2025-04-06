package models

import "time"

type Club struct {
	ID          string   `json:"id" gorm:"type:uuid;primary_key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []Member `json:"members,omitempty" gorm:"foreignKey:ClubID"`
	Events      []Event  `json:"events,omitempty" gorm:"foreignKey:ClubID"`
}

type Member struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
}

type Event struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClubID      string `json:"club_id" gorm:"type:uuid"`
	Date        string `json:"date"`
	BeginTime   string `json:"begin_time"`
	EndTime     string `json:"end_time"`
}

type MagicLink struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
