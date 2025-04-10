package models

import (
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/notifications"
	"github.com/google/uuid"
)

type Member struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	ClubID string `json:"club_id" gorm:"type:uuid"`
}

func GetClubMembers(clubID string) ([]Member, error) {
	var members []Member
	err := database.Db.Where("club_id = ?", clubID).Find(&members).Error
	return members, err
}

func AddMember(member *Member, clubID string) error {
	member.ID = uuid.New().String()
	member.ClubID = clubID
	return database.Db.Create(member).Error
}

func DeleteMember(memberID, clubID string) (int64, error) {
	result := database.Db.Where("id = ? AND club_id = ?", memberID, clubID).Delete(&Member{})
	return result.RowsAffected, result.Error
}

func (m *Member) Validate() bool {
	return m.Email != "" && m.Name != ""
}

func (m *Member) NotifyAdded(clubName string) {
	notifications.SendMemberAddedNotification(m.Email, clubName, m.ClubID)
}
