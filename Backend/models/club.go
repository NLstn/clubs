package models

import (
	"log"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/notifications"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Club struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Member struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
	UserID string `json:"user_id" gorm:"type:uuid"`
	Role   string `json:"role"`
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

func CreateClub(name, description, ownerID string) (Club, error) {
	var club Club
	club.ID = uuid.New().String()
	club.Name = name
	club.Description = description

	err := database.Db.Transaction(func(tx *gorm.DB) error {
		if dbErr := tx.Create(&club).Error; dbErr != nil {
			return dbErr
		}
		var member Member
		member.ID = uuid.New().String()
		member.ClubID = club.ID
		member.UserID = ownerID
		member.Role = "owner"
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

func (c *Club) Update(name, description string) error {
	return database.Db.Model(c).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
	}).Error
}

func (c *Club) IsOwner(user User) bool {
	role, err := c.GetMemberRole(user)
	if err != nil {
		return false
	}
	if role == "owner" {
		return true
	}
	return false
}

func (c *Club) IsAdmin(user User) bool {
	role, err := c.GetMemberRole(user)
	if err != nil {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	return false
}

func (c *Club) IsMember(user User) bool {
	if user.ID == "" {
		log.Fatal("User ID is empty")
		return false
	}

	result := database.Db.Where("club_id = ? AND user_id = ?", c.ID, user.ID).Limit(1).Find(&Member{})
	if result.Error != nil {
		log.Default().Println("Error checking membership:", result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func (c *Club) GetClubMembers() ([]Member, error) {
	var members []Member
	err := database.Db.Where("club_id = ?", c.ID).Find(&members).Error
	return members, err
}

func (c *Club) AddMember(userId, role string) error {
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

func (c *Club) GetMemberRole(user User) (string, error) {
	var member Member
	result := database.Db.Where("club_id = ? AND user_id = ?", c.ID, user.ID).First(&member)
	if result.Error != nil {
		return "", result.Error
	}
	return member.Role, nil
}

func (c *Club) UpdateMemberRole(changingUser User, memberID, role string) error {
	var member Member
	result := database.Db.Where("id = ? AND club_id = ?", memberID, c.ID).First(&member)
	if result.Error != nil {
		return result.Error
	}
	if role != "owner" && role != "admin" && role != "member" {
		return gorm.ErrInvalidData
	}

	if canChange, err := c.canChangeRole(changingUser, member.Role, role); err != nil || !canChange {
		return gorm.ErrInvalidData
	}

	member.Role = role
	return database.Db.Save(&member).Error
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

func (c *Club) canChangeRole(changingUser User, oldRole, newRole string) (bool, error) {
	changingUserRole, err := c.GetMemberRole(changingUser)
	if err != nil {
		return false, err
	}
	if changingUserRole == "owner" {
		return true, nil
	}
	if changingUserRole == "admin" && (oldRole == "member" || newRole == "admin") {
		return true, nil
	}
	return false, nil
}
