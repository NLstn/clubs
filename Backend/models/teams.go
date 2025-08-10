package models

import (
	"errors"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotTeamAdmin = errors.New("user is not a team admin")
var ErrNotClubAdminOrTeamAdmin = errors.New("user is not a club admin or team admin")
var ErrLastTeamAdminDemotion = errors.New("cannot demote the last admin of the team")

type Team struct {
	ID          string     `json:"id" gorm:"type:uuid;primary_key"`
	ClubID      string     `json:"clubId" gorm:"type:uuid;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	CreatedBy   string     `json:"createdBy" gorm:"type:uuid"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	UpdatedBy   string     `json:"updatedBy" gorm:"type:uuid"`
	Deleted     bool       `json:"deleted" gorm:"default:false"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *string    `json:"deletedBy,omitempty" gorm:"type:uuid"`
}

type TeamMember struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	TeamID    string    `json:"teamId" gorm:"type:uuid;not null"`
	UserID    string    `json:"userId" gorm:"type:uuid;not null"`
	Role      string    `json:"role" gorm:"default:member"` // admin, member
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy" gorm:"type:uuid"`
}

// BeforeCreate generates UUID for new team
func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}

// BeforeCreate generates UUID for new team member
func (tm *TeamMember) BeforeCreate(tx *gorm.DB) (err error) {
	if tm.ID == "" {
		tm.ID = uuid.New().String()
	}
	return
}

// CreateTeam creates a new team within a club
func (c *Club) CreateTeam(name, description, createdByUserID string) (Team, error) {
	team := Team{
		ClubID:      c.ID,
		Name:        name,
		Description: description,
		CreatedBy:   createdByUserID,
		UpdatedBy:   createdByUserID,
	}

	err := database.Db.Create(&team).Error
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

// GetTeams returns all teams for a club
func (c *Club) GetTeams() ([]Team, error) {
	var teams []Team
	err := database.Db.Where("club_id = ? AND deleted = false", c.ID).Find(&teams).Error
	return teams, err
}

// GetTeamByID returns a team by ID
func GetTeamByID(teamID string) (Team, error) {
	var team Team
	err := database.Db.Where("id = ? AND deleted = false", teamID).First(&team).Error
	return team, err
}

// Update updates team information
func (t *Team) Update(name, description, updatedBy string) error {
	return database.Db.Model(t).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
		"updated_by":  updatedBy,
		"updated_at":  time.Now(),
	}).Error
}

// SoftDelete soft deletes a team
func (t *Team) SoftDelete(deletedBy string) error {
	now := time.Now()
	return database.Db.Model(t).Updates(map[string]interface{}{
		"deleted":    true,
		"deleted_at": &now,
		"deleted_by": &deletedBy,
		"updated_at": now,
	}).Error
}

// AddMember adds a user to the team
func (t *Team) AddMember(userID, role, addedBy string) error {
	teamMember := TeamMember{
		TeamID:    t.ID,
		UserID:    userID,
		Role:      role,
		CreatedBy: addedBy,
		UpdatedBy: addedBy,
	}

	return database.Db.Create(&teamMember).Error
}

// GetMembers returns all members of the team
func (t *Team) GetMembers() ([]TeamMember, error) {
	var members []TeamMember
	err := database.Db.Where("team_id = ?", t.ID).Find(&members).Error
	return members, err
}

// GetTeamMembersWithUserInfo returns team members with user information
func (t *Team) GetTeamMembersWithUserInfo() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			tm.id,
			tm.user_id,
			tm.role,
			tm.created_at as joined_at,
			CONCAT(u.first_name, ' ', u.last_name) as name
		FROM team_members tm
		JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = ?
		ORDER BY 
			CASE tm.role 
				WHEN 'admin' THEN 1 
				WHEN 'member' THEN 2 
				ELSE 3 
			END,
			CONCAT(u.first_name, ' ', u.last_name)
	`

	rows, err := database.Db.Raw(query, t.ID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var memberInfo struct {
			ID       string    `json:"id"`
			UserID   string    `json:"userId"`
			Role     string    `json:"role"`
			JoinedAt time.Time `json:"joinedAt"`
			Name     string    `json:"name"`
		}

		err := database.Db.ScanRows(rows, &memberInfo)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":       memberInfo.ID,
			"userId":   memberInfo.UserID,
			"role":     memberInfo.Role,
			"joinedAt": memberInfo.JoinedAt.Format("2006-01-02T15:04:05Z"),
			"name":     memberInfo.Name,
		})
	}

	return results, nil
}

// RemoveMember removes a user from the team
func (t *Team) RemoveMember(memberID string) error {
	return database.Db.Where("id = ? AND team_id = ?", memberID, t.ID).Delete(&TeamMember{}).Error
}

// UpdateMemberRole updates a team member's role
func (t *Team) UpdateMemberRole(changingUser User, memberID, newRole string) error {
	var teamMember TeamMember
	err := database.Db.Where("id = ? AND team_id = ?", memberID, t.ID).First(&teamMember).Error
	if err != nil {
		return err
	}

	// Validate role
	if newRole != "admin" && newRole != "member" {
		return errors.New("invalid role")
	}

	// Check permissions
	canChange, err := t.canChangeRole(changingUser, teamMember, newRole)
	if err != nil {
		return err
	}
	if !canChange {
		return ErrNotTeamAdmin
	}

	// Update role
	teamMember.Role = newRole
	teamMember.UpdatedBy = changingUser.ID
	teamMember.UpdatedAt = time.Now()

	return database.Db.Save(&teamMember).Error
}

// canChangeRole checks if a user can change another member's role
func (t *Team) canChangeRole(changingUser User, targetMember TeamMember, newRole string) (bool, error) {
	// Get the club this team belongs to
	var club Club
	err := database.Db.Where("id = ?", t.ClubID).First(&club).Error
	if err != nil {
		return false, err
	}

	// Club owners can change any role
	if club.IsOwner(changingUser) {
		// Check if this would demote the last team admin
		if targetMember.Role == "admin" && newRole != "admin" {
			adminCount, err := t.CountAdmins()
			if err != nil {
				return false, err
			}
			if adminCount <= 1 {
				return false, ErrLastTeamAdminDemotion
			}
		}
		return true, nil
	}

	// Club admins can change any role
	if club.IsAdmin(changingUser) {
		// Check if this would demote the last team admin
		if targetMember.Role == "admin" && newRole != "admin" {
			adminCount, err := t.CountAdmins()
			if err != nil {
				return false, err
			}
			if adminCount <= 1 {
				return false, ErrLastTeamAdminDemotion
			}
		}
		return true, nil
	}

	// Team admins can promote members to admin or demote admins to member
	if t.IsAdmin(changingUser) {
		// Check if this would demote the last team admin (including themselves)
		if targetMember.Role == "admin" && newRole != "admin" {
			adminCount, err := t.CountAdmins()
			if err != nil {
				return false, err
			}
			if adminCount <= 1 {
				return false, ErrLastTeamAdminDemotion
			}
		}
		return true, nil
	}

	return false, nil
}

// IsAdmin checks if a user is an admin of the team
func (t *Team) IsAdmin(user User) bool {
	var count int64
	database.Db.Model(&TeamMember{}).Where("team_id = ? AND user_id = ? AND role = ?", t.ID, user.ID, "admin").Count(&count)
	return count > 0
}

// IsMember checks if a user is a member of the team
func (t *Team) IsMember(user User) bool {
	var count int64
	database.Db.Model(&TeamMember{}).Where("team_id = ? AND user_id = ?", t.ID, user.ID).Count(&count)
	return count > 0
}

// CountAdmins returns the number of admins in the team
func (t *Team) CountAdmins() (int64, error) {
	var count int64
	err := database.Db.Model(&TeamMember{}).Where("team_id = ? AND role = ?", t.ID, "admin").Count(&count).Error
	return count, err
}

// GetUserRole returns the role of a user in the team
func (t *Team) GetUserRole(user User) (string, error) {
	var teamMember TeamMember
	err := database.Db.Where("team_id = ? AND user_id = ?", t.ID, user.ID).First(&teamMember).Error
	if err != nil {
		return "", err
	}
	return teamMember.Role, nil
}

// GetUserTeams returns all teams a user belongs to within a specific club
func GetUserTeams(userID, clubID string) ([]Team, error) {
	var teams []Team

	query := `
		SELECT DISTINCT t.* 
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = ? AND t.club_id = ? AND t.deleted = false
	`

	err := database.Db.Raw(query, userID, clubID).Find(&teams).Error
	return teams, err
}

// CanUserCreateTeam checks if a user can create teams in a club
func (c *Club) CanUserCreateTeam(user User) bool {
	return c.IsAdmin(user) // Only club admins can create teams
}

// CanUserEditTeam checks if a user can edit a team
func (t *Team) CanUserEditTeam(user User) bool {
	// Get the club this team belongs to
	var club Club
	err := database.Db.Where("id = ?", t.ClubID).First(&club).Error
	if err != nil {
		return false
	}

	// Club owners/admins can edit any team
	if club.IsAdmin(user) {
		return true
	}

	// Team admins can edit their team
	return t.IsAdmin(user)
}

// CanUserDeleteTeam checks if a user can delete a team
func (t *Team) CanUserDeleteTeam(user User) bool {
	// Get the club this team belongs to
	var club Club
	err := database.Db.Where("id = ?", t.ClubID).First(&club).Error
	if err != nil {
		return false
	}

	// Only club owners/admins can delete teams
	return club.IsAdmin(user)
}

// GetTeamStats returns statistics for the team
func (t *Team) GetTeamStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Get member count
	var memberCount int64
	err := database.Db.Model(&TeamMember{}).Where("team_id = ?", t.ID).Count(&memberCount).Error
	if err != nil {
		return nil, err
	}
	stats["member_count"] = memberCount
	
	// Get admin count
	var adminCount int64
	err = database.Db.Model(&TeamMember{}).Where("team_id = ? AND role = ?", t.ID, "admin").Count(&adminCount).Error
	if err != nil {
		return nil, err
	}
	stats["admin_count"] = adminCount
	
	// Get upcoming events count
	var upcomingEventCount int64
	now := time.Now()
	err = database.Db.Model(&Event{}).Where("team_id = ? AND start_time >= ?", t.ID, now).Count(&upcomingEventCount).Error
	if err != nil {
		return nil, err
	}
	stats["upcoming_events"] = upcomingEventCount
	
	// Get total events count
	var totalEventCount int64
	err = database.Db.Model(&Event{}).Where("team_id = ?", t.ID).Count(&totalEventCount).Error
	if err != nil {
		return nil, err
	}
	stats["total_events"] = totalEventCount
	
	// Get unpaid fines count
	var unpaidFineCount int64
	err = database.Db.Model(&Fine{}).Where("team_id = ? AND paid = false", t.ID).Count(&unpaidFineCount).Error
	if err != nil {
		return nil, err
	}
	stats["unpaid_fines"] = unpaidFineCount
	
	// Get total fines count
	var totalFineCount int64
	err = database.Db.Model(&Fine{}).Where("team_id = ?", t.ID).Count(&totalFineCount).Error
	if err != nil {
		return nil, err
	}
	stats["total_fines"] = totalFineCount
	
	return stats, nil
}
