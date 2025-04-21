package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/notifications"
	"github.com/google/uuid"
)

type Club struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id" gorm:"type:uuid"`
}

type Member struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
	UserID string `json:"user_id" gorm:"type:uuid"`
}

func CreateClub(club *Club, ownerID string) error {
	club.ID = uuid.New().String()
	club.OwnerID = ownerID
	return database.Db.Create(club).Error
}

func GetAllClubs() ([]Club, error) {
	var clubs []Club
	err := database.Db.Find(&clubs).Error
	return clubs, err
}

func GetClubByID(id string) (Club, error) {
	var club Club
	result := database.Db.First(&club, "id = ?", id)
	return club, result.Error
}

func (c *Club) IsOwner(userID string) bool {
	return c.OwnerID == userID
}

func (c *Club) GetClubMembers() ([]Member, error) {
	var members []Member
	err := database.Db.Where("club_id = ?", c.ID).Find(&members).Error
	return members, err
}

func (c *Club) AddMember(userId string) error {
	var member Member
	member.ID = uuid.New().String()
	member.ClubID = c.ID
	member.UserID = userId
	err := database.Db.Create(member).Error
	if err != nil {
		return err
	}
	member.notifyAdded()
	return nil
}

func (c *Club) DeleteMember(memberID string) (int64, error) {
	result := database.Db.Where("id = ? AND club_id = ?", memberID, c.ID).Delete(&Member{})
	return result.RowsAffected, result.Error
}

func (m *Member) notifyAdded() {
	var club Club
	if err := database.Db.Where("id = ?", m.ClubID).First(&club).Error; err != nil {
		return
	}

	var user User
	if err := database.Db.Where("id = ?", m.UserID).First(&user).Error; err != nil {
		return
	}

	notifications.SendMemberAddedNotification(user.Email, club.ID, club.Name)
}
