package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotTeamAdmin = errors.New("user is not a team admin")
var ErrNotClubAdminOrTeamAdmin = errors.New("user is not a club admin or team admin")
var ErrLastTeamAdminDemotion = errors.New("cannot demote the last admin of the team")

type Team struct {
	ID          string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	ClubID      string    `json:"ClubID" gorm:"type:uuid;not null" odata:"required"`
	Name        string    `json:"Name" gorm:"not null" odata:"required"`
	Description *string   `json:"Description,omitempty" odata:"nullable"`
	CreatedAt   time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy   string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	UpdatedBy   string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`

	// Navigation properties for OData
	Events      []Event      `gorm:"foreignKey:TeamID" json:"Events,omitempty" odata:"nav"`
	Fines       []Fine       `gorm:"foreignKey:TeamID" json:"Fines,omitempty" odata:"nav"`
	TeamMembers []TeamMember `gorm:"foreignKey:TeamID" json:"TeamMembers,omitempty" odata:"nav"`
}

type TeamMember struct {
	ID        string    `json:"ID" gorm:"type:uuid;primary_key" odata:"key"`
	TeamID    string    `json:"TeamID" gorm:"type:uuid;not null" odata:"required"`
	UserID    string    `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
	Role      string    `json:"Role" gorm:"default:member" odata:"required"` // admin, member
	CreatedAt time.Time `json:"CreatedAt" odata:"immutable"`
	CreatedBy string    `json:"CreatedBy" gorm:"type:uuid" odata:"required"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	UpdatedBy string    `json:"UpdatedBy" gorm:"type:uuid" odata:"required"`

	// Navigation properties for OData
	User User `gorm:"foreignKey:UserID" json:"User,omitempty" odata:"nav"`
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
		Description: &description,
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
	err := database.Db.Where("club_id = ?", c.ID).Find(&teams).Error
	return teams, err
}

// GetTeamByID returns a team by ID
func GetTeamByID(teamID string) (Team, error) {
	var team Team
	err := database.Db.Where("id = ?", teamID).First(&team).Error
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
			ID       string    `json:"ID"`
			UserID   string    `json:"UserID"`
			Role     string    `json:"Role"`
			JoinedAt time.Time `json:"JoinedAt"`
			Name     string    `json:"Name"`
		}

		err := database.Db.ScanRows(rows, &memberInfo)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"ID":       memberInfo.ID,
			"UserID":   memberInfo.UserID,
			"Role":     memberInfo.Role,
			"JoinedAt": memberInfo.JoinedAt.Format("2006-01-02T15:04:05Z"),
			"Name":     memberInfo.Name,
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

// ODataBeforeReadCollection filters teams to only those in clubs the user belongs to
func (t Team) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see teams of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific team
func (t Team) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see teams of clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates team creation permissions
func (t *Team) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", t.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can create teams")
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.CreatedBy = userID
	t.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates team update permissions
func (t *Team) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", t.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can update teams")
	}

	// Set UpdatedBy
	now := time.Now()
	t.UpdatedAt = now
	t.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates team deletion permissions
func (t *Team) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Check if user is an admin/owner of the club
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", t.ClubID, userID).First(&existingMember).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins and owners can delete teams")
	}

	return nil
}

// TeamMember authorization hooks
// ODataBeforeReadCollection filters team members to only those in teams the user can access
func (tm TeamMember) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see team members of teams in clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("team_id IN (SELECT id FROM teams WHERE club_id IN (SELECT club_id FROM members WHERE user_id = ?))", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeReadEntity validates access to a specific team member record
func (tm TeamMember) ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// User can only see team members of teams in clubs they belong to
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("team_id IN (SELECT id FROM teams WHERE club_id IN (SELECT club_id FROM members WHERE user_id = ?))", userID)
	}

	return []func(*gorm.DB) *gorm.DB{scope}, nil
}

// ODataBeforeCreate validates team member creation permissions
func (tm *TeamMember) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get team to find club ID
	var team Team
	if err := database.Db.Where("id = ?", tm.TeamID).First(&team).Error; err != nil {
		return fmt.Errorf("team not found")
	}

	// Check if user is an admin/owner of the club or team admin
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", team.ClubID, userID).First(&existingMember).Error; err != nil {
		// Check if user is team admin
		var teamMember TeamMember
		if err := database.Db.Where("team_id = ? AND user_id = ? AND role = 'admin'", tm.TeamID, userID).First(&teamMember).Error; err != nil {
			return fmt.Errorf("unauthorized: only club admins/owners or team admins can add team members")
		}
	}

	// Set CreatedBy and UpdatedBy
	now := time.Now()
	tm.CreatedAt = now
	tm.UpdatedAt = now
	tm.CreatedBy = userID
	tm.UpdatedBy = userID

	return nil
}

// ODataBeforeUpdate validates team member update permissions
func (tm *TeamMember) ODataBeforeUpdate(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get team to find club ID
	var team Team
	if err := database.Db.Where("id = ?", tm.TeamID).First(&team).Error; err != nil {
		return fmt.Errorf("team not found")
	}

	// Check if user is an admin/owner of the club or team admin
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", team.ClubID, userID).First(&existingMember).Error; err != nil {
		// Check if user is team admin
		var teamMember TeamMember
		if err := database.Db.Where("team_id = ? AND user_id = ? AND role = 'admin'", tm.TeamID, userID).First(&teamMember).Error; err != nil {
			return fmt.Errorf("unauthorized: only club admins/owners or team admins can update team members")
		}
	}

	// Set UpdatedBy
	now := time.Now()
	tm.UpdatedAt = now
	tm.UpdatedBy = userID

	return nil
}

// ODataBeforeDelete validates team member deletion permissions
func (tm *TeamMember) ODataBeforeDelete(ctx context.Context, r *http.Request) error {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get team to find club ID
	var team Team
	if err := database.Db.Where("id = ?", tm.TeamID).First(&team).Error; err != nil {
		return fmt.Errorf("team not found")
	}

	// Users can leave teams (delete their own membership)
	if tm.UserID == userID {
		return nil
	}

	// Check if user is an admin/owner of the club or team admin
	var existingMember Member
	if err := database.Db.Where("club_id = ? AND user_id = ? AND role IN ('admin', 'owner')", team.ClubID, userID).First(&existingMember).Error; err != nil {
		// Check if user is team admin
		var teamMember TeamMember
		if err := database.Db.Where("team_id = ? AND user_id = ? AND role = 'admin'", tm.TeamID, userID).First(&teamMember).Error; err != nil {
			return fmt.Errorf("unauthorized: only club admins/owners or team admins can remove team members")
		}
	}

	return nil
}
