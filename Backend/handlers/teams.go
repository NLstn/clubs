package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerTeamRoutes(mux *http.ServeMux) {
	// GET /api/v1/clubs/{clubid}/teams
	mux.Handle("/api/v1/clubs/{clubid}/teams", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request for user teams via query param
		if r.URL.Query().Has("user") {
			// Handle get user teams
			handleGetUserTeams(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetClubTeams(w, r)
		case http.MethodPost:
			handleCreateTeam(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// GET/PUT/DELETE /api/v1/clubs/{clubid}/teams/{teamid}
	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeam(w, r)
		case http.MethodPut:
			handleUpdateTeam(w, r)
		case http.MethodDelete:
			handleDeleteTeam(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// GET/POST /api/v1/clubs/{clubid}/teams/{teamid}/members
	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/members", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamMembers(w, r)
		case http.MethodPost:
			handleAddTeamMember(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// PATCH/DELETE /api/v1/clubs/{clubid}/teams/{teamid}/members/{memberid}
	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/members/{memberid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handleUpdateTeamMemberRole(w, r)
		case http.MethodDelete:
			handleRemoveTeamMember(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

}

// endpoint: GET /api/v1/clubs/{clubid}/teams
func handleGetClubTeams(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is a member of the club
	if !club.IsMember(user) {
		http.Error(w, "User is not a member of this club", http.StatusForbidden)
		return
	}

	teams, err := club.GetTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// endpoint: POST /api/v1/clubs/{clubid}/teams
func handleCreateTeam(w http.ResponseWriter, r *http.Request) {
	type CreateTeamRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		http.Error(w, "Team name is required", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user can create teams (only club admins)
	if !club.CanUserCreateTeam(user) {
		http.Error(w, "Only club admins can create teams", http.StatusForbidden)
		return
	}

	team, err := club.CreateTeam(payload.Name, payload.Description, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

// endpoint: GET /api/v1/clubs/{clubid}/teams/{teamid}
func handleGetTeam(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is a member of the club
	if !club.IsMember(user) {
		http.Error(w, "User is not a member of this club", http.StatusForbidden)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// endpoint: PUT /api/v1/clubs/{clubid}/teams/{teamid}
func handleUpdateTeam(w http.ResponseWriter, r *http.Request) {
	type UpdateTeamRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		http.Error(w, "Team name is required", http.StatusBadRequest)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	// Check if user can edit the team
	if !team.CanUserEditTeam(user) {
		http.Error(w, "Insufficient permissions to edit this team", http.StatusForbidden)
		return
	}

	err = team.Update(payload.Name, payload.Description, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/teams/{teamid}
func handleDeleteTeam(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	// Check if user can delete the team
	if !team.CanUserDeleteTeam(user) {
		http.Error(w, "Insufficient permissions to delete this team", http.StatusForbidden)
		return
	}

	err = team.SoftDelete(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/clubs/{clubid}/teams/{teamid}/members
func handleGetTeamMembers(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is a member of the club
	if !club.IsMember(user) {
		http.Error(w, "User is not a member of this club", http.StatusForbidden)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	members, err := team.GetTeamMembersWithUserInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// endpoint: POST /api/v1/clubs/{clubid}/teams/{teamid}/members
func handleAddTeamMember(w http.ResponseWriter, r *http.Request) {
	type AddMemberRequest struct {
		UserID string `json:"userId"`
		Role   string `json:"role"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(payload.UserID); err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Default role to member if not specified
	if payload.Role == "" {
		payload.Role = "member"
	}

	// Validate role
	if payload.Role != "admin" && payload.Role != "member" {
		http.Error(w, "Invalid role. Must be 'admin' or 'member'", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	// Check if user can edit the team
	if !team.CanUserEditTeam(user) {
		http.Error(w, "Insufficient permissions to manage team members", http.StatusForbidden)
		return
	}

	// Check if the user being added is a member of the club
	targetUser, err := models.GetUserByID(payload.UserID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsMember(targetUser) {
		http.Error(w, "User is not a member of the club", http.StatusBadRequest)
		return
	}

	// Check if user is already a member of the team
	if team.IsMember(targetUser) {
		http.Error(w, "User is already a member of this team", http.StatusBadRequest)
		return
	}

	err = team.AddMember(payload.UserID, payload.Role, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Member added to team successfully"}`))
}

// endpoint: PATCH /api/v1/clubs/{clubid}/teams/{teamid}/members/{memberid}
func handleUpdateTeamMemberRole(w http.ResponseWriter, r *http.Request) {
	type UpdateRoleRequest struct {
		Role string `json:"role"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	memberID := extractPathParam(r, "members")
	if _, err := uuid.Parse(memberID); err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Role != "admin" && payload.Role != "member" {
		http.Error(w, "Invalid role. Must be 'admin' or 'member'", http.StatusBadRequest)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	err = team.UpdateMemberRole(user, memberID, payload.Role)
	if err != nil {
		if err == models.ErrNotTeamAdmin {
			http.Error(w, "Insufficient permissions to change member roles", http.StatusForbidden)
			return
		}
		if err == models.ErrLastTeamAdminDemotion {
			http.Error(w, "Cannot demote the last admin of the team", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/teams/{teamid}/members/{memberid}
func handleRemoveTeamMember(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	teamID := extractPathParam(r, "teams")
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	memberID := extractPathParam(r, "members")
	if _, err := uuid.Parse(memberID); err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return
	}

	// Check if user can edit the team
	if !team.CanUserEditTeam(user) {
		http.Error(w, "Insufficient permissions to remove team members", http.StatusForbidden)
		return
	}

	err = team.RemoveMember(memberID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/clubs/{clubid}/teams
func handleGetUserTeams(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is a member of the club
	if !club.IsMember(user) {
		http.Error(w, "User is not a member of this club", http.StatusForbidden)
		return
	}

	teams, err := models.GetUserTeams(user.ID, clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
