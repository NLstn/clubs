package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/notifications"
	"github.com/google/uuid"
)

type Member struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
	UserID string `json:"user_id" gorm:"type:uuid"`
}

func GetClubMembers(clubID string) ([]Member, error) {
	var members []Member
	err := database.Db.Where("club_id = ?", clubID).Find(&members).Error
	return members, err
}

func AddMember(clubID, userId string) error {
	var member Member
	member.ID = uuid.New().String()
	member.ClubID = clubID
	member.UserID = userId
	err := database.Db.Create(member).Error
	if err != nil {
		return err
	}
	member.notifyAdded()
	return nil
}

func DeleteMember(memberID, clubID string) (int64, error) {
	result := database.Db.Where("id = ? AND club_id = ?", memberID, clubID).Delete(&Member{})
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
