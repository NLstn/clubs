package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/database"
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

	// Team overview/stats endpoint
	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/overview", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamOverview(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Team events endpoints
	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/events", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamEvents(w, r)
		case http.MethodPost:
			handleCreateTeamEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/events/{eventid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamEvent(w, r)
		case http.MethodPut:
			handleUpdateTeamEvent(w, r)
		case http.MethodDelete:
			handleDeleteTeamEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/events/upcoming", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetTeamUpcomingEvents(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Team fines endpoints
	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/fines", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamFines(w, r)
		case http.MethodPost:
			handleCreateTeamFine(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/fines/{fineid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleDeleteTeamFine(w, r)
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

	// Hard delete the team
	err = database.Db.Delete(&team).Error
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

// Team Overview endpoint
func handleGetTeamOverview(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	teamID := extractPathParam(r, "teams")
	user := extractUser(r)

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
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

	stats, err := team.GetTeamStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user's role in the team
	userRole := ""
	if team.IsMember(user) {
		userRole, _ = team.GetUserRole(user)
	}

	response := map[string]interface{}{
		"team":      team,
		"stats":     stats,
		"user_role": userRole,
		"is_admin":  team.IsAdmin(user),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Team Events handlers
func handleGetTeamEvents(w http.ResponseWriter, r *http.Request) {
	user, club, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	// Check if user is a member of the team or club admin
	if !team.IsMember(user) && !club.IsAdmin(user) {
		http.Error(w, "User is not a member of this team", http.StatusForbidden)
		return
	}

	events, err := team.GetEvents()
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func handleGetTeamUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	user, club, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	// Check if user is a member of the team or club admin
	if !team.IsMember(user) && !club.IsAdmin(user) {
		http.Error(w, "User is not a member of this team", http.StatusForbidden)
		return
	}

	events, err := team.GetUpcomingEvents()
	if err != nil {
		http.Error(w, "Failed to get upcoming events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func handleCreateTeamEvent(w http.ResponseWriter, r *http.Request) {
	type CreateEventRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	user, _, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	if !team.IsAdmin(user) {
		http.Error(w, "Unauthorized - team admin access required", http.StatusForbidden)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse timestamps
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format. Expected RFC3339 timestamp", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format. Expected RFC3339 timestamp", http.StatusBadRequest)
		return
	}

	if startTime.After(endTime) || startTime.Equal(endTime) {
		http.Error(w, "Start time must be before end time", http.StatusBadRequest)
		return
	}

	event, err := team.CreateEvent(req.Name, req.Description, req.Location, startTime, endTime, user.ID)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func handleGetTeamEvent(w http.ResponseWriter, r *http.Request) {
	user, club, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	// Check if user is a member of the team or club admin
	if !team.IsMember(user) && !club.IsAdmin(user) {
		http.Error(w, "User is not a member of this team", http.StatusForbidden)
		return
	}

	event, err := team.GetEventByID(eventID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		return
	}

	// Get user's RSVP if it exists
	userRSVP, _ := user.GetUserRSVP(eventID)

	response := struct {
		*models.Event
		UserRSVP *models.EventRSVP `json:"user_rsvp,omitempty"`
	}{
		Event:    event,
		UserRSVP: userRSVP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleUpdateTeamEvent(w http.ResponseWriter, r *http.Request) {
	type UpdateEventRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	user, _, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	if !team.IsAdmin(user) {
		http.Error(w, "Unauthorized - team admin access required", http.StatusForbidden)
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse timestamps
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format. Expected RFC3339 timestamp", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format. Expected RFC3339 timestamp", http.StatusBadRequest)
		return
	}

	if startTime.After(endTime) || startTime.Equal(endTime) {
		http.Error(w, "Start time must be before end time", http.StatusBadRequest)
		return
	}

	event, err := team.UpdateEvent(eventID, req.Name, req.Description, req.Location, startTime, endTime, user.ID)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func handleDeleteTeamEvent(w http.ResponseWriter, r *http.Request) {
	user, _, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	if !team.IsAdmin(user) {
		http.Error(w, "Unauthorized - team admin access required", http.StatusForbidden)
		return
	}

	err = team.DeleteEvent(eventID)
	if err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Team Fines handlers
func handleGetTeamFines(w http.ResponseWriter, r *http.Request) {
	user, club, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	// Check if user is a member of the team or club admin
	if !team.IsMember(user) && !club.IsAdmin(user) {
		http.Error(w, "User is not a member of this team", http.StatusForbidden)
		return
	}

	fines, err := team.GetFines()
	if err != nil {
		http.Error(w, "Failed to retrieve fines", http.StatusInternalServerError)
		return
	}

	type Fine struct {
		ID        string  `json:"id"`
		UserID    string  `json:"userId"`
		UserName  string  `json:"userName"`
		Reason    string  `json:"reason"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt string  `json:"updatedAt"`
		Paid      bool    `json:"paid"`
	}

	var fineList []Fine
	for _, fine := range fines {
		user, err := models.GetUserByID(fine.UserID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		fineList = append(fineList, Fine{
			ID:        fine.ID,
			UserID:    fine.UserID,
			UserName:  user.GetFullName(),
			Reason:    fine.Reason,
			Amount:    fine.Amount,
			CreatedAt: fine.CreatedAt.Format(time.RFC3339),
			UpdatedAt: fine.UpdatedAt.Format(time.RFC3339),
			Paid:      fine.Paid,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fineList)
}

func handleCreateTeamFine(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		UserID string  `json:"userId"`
		Reason string  `json:"reason"`
		Amount float64 `json:"amount"`
	}

	user, _, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	if !team.IsAdmin(user) {
		http.Error(w, "Unauthorized - team admin access required", http.StatusForbidden)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.UserID == "" || payload.Reason == "" || payload.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	fine, err := team.CreateFine(payload.UserID, payload.Reason, user.ID, payload.Amount)
	if err != nil {
		http.Error(w, "Failed to create fine", http.StatusInternalServerError)
		return
	}

	type FineResponse struct {
		ID        string  `json:"id"`
		UserID    string  `json:"userId"`
		Reason    string  `json:"reason"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt string  `json:"updatedAt"`
		Paid      bool    `json:"paid"`
	}

	resp := FineResponse{
		ID:        fine.ID,
		UserID:    fine.UserID,
		Reason:    fine.Reason,
		Amount:    fine.Amount,
		CreatedAt: fine.CreatedAt.Format(time.RFC3339),
		UpdatedAt: fine.UpdatedAt.Format(time.RFC3339),
		Paid:      fine.Paid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func handleDeleteTeamFine(w http.ResponseWriter, r *http.Request) {
	user, _, team, err := validateTeamAccess(w, r)
	if err != nil {
		return
	}

	fineID := extractPathParam(r, "fines")
	if _, err := uuid.Parse(fineID); err != nil {
		http.Error(w, "Invalid fine ID format", http.StatusBadRequest)
		return
	}

	if !team.IsAdmin(user) {
		http.Error(w, "Unauthorized - team admin access required", http.StatusForbidden)
		return
	}

	err = team.DeleteFine(fineID)
	if err != nil {
		http.Error(w, "Failed to delete fine", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to validate team access
func validateTeamAccess(w http.ResponseWriter, r *http.Request) (models.User, models.Club, models.Team, error) {
	clubID := extractPathParam(r, "clubs")
	teamID := extractPathParam(r, "teams")
	user := extractUser(r)

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return models.User{}, models.Club{}, models.Team{}, err
	}

	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return models.User{}, models.Club{}, models.Team{}, err
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return models.User{}, models.Club{}, models.Team{}, err
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.User{}, models.Club{}, models.Team{}, err
	}

	// Check if user is a member of the club
	if !club.IsMember(user) {
		http.Error(w, "User is not a member of this club", http.StatusForbidden)
		return models.User{}, models.Club{}, models.Team{}, err
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return models.User{}, models.Club{}, models.Team{}, err
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.User{}, models.Club{}, models.Team{}, err
	}

	// Verify team belongs to the club
	if team.ClubID != clubID {
		http.Error(w, "Team does not belong to this club", http.StatusBadRequest)
		return models.User{}, models.Club{}, models.Team{}, err
	}

	return user, club, team, nil
}
