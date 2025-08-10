package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// registerTeamEventRoutes registers routes for team events
func registerTeamEventRoutes(mux *http.ServeMux) {
	base := "/api/v1/clubs/{clubid}/teams/{teamid}/events"

	mux.Handle(base, RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamEvents(w, r)
		case http.MethodPost:
			handleCreateTeamEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle(base+"/{eventid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
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

	mux.Handle(base+"/upcoming", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetTeamUpcomingEvents(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle(base+"/{eventid}/rsvp", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleCreateOrUpdateTeamEventRSVP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle(base+"/{eventid}/rsvps", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetTeamEventRSVPs(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// helper to fetch club and team and validate IDs
func getClubAndTeam(r *http.Request, w http.ResponseWriter) (models.Club, models.Team, bool) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	teamID := extractPathParam(r, "teams")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return models.Club{}, models.Team{}, false
	}
	if _, err := uuid.Parse(teamID); err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return models.Club{}, models.Team{}, false
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return models.Club{}, models.Team{}, false
	} else if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return models.Club{}, models.Team{}, false
	}

	team, err := models.GetTeamByID(teamID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Team not found", http.StatusNotFound)
		return models.Club{}, models.Team{}, false
	} else if err != nil {
		http.Error(w, "Failed to get team information", http.StatusInternalServerError)
		return models.Club{}, models.Team{}, false
	}

	if team.ClubID != club.ID {
		http.Error(w, "Team does not belong to club", http.StatusBadRequest)
		return models.Club{}, models.Team{}, false
	}

	// Check membership for general access
	if !team.IsMember(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return models.Club{}, models.Team{}, false
	}

	return club, team, true
}

// GET /api/v1/clubs/{clubid}/teams/{teamid}/events
func handleGetTeamEvents(w http.ResponseWriter, r *http.Request) {
	_, team, ok := getClubAndTeam(r, w)
	if !ok {
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

// GET /api/v1/clubs/{clubid}/teams/{teamid}/events/{eventid}
func handleGetTeamEvent(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	_, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}
	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	event, err := team.GetEventByID(eventID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		return
	}

	userRSVP, _ := user.GetUserTeamEventRSVP(eventID)
	response := struct {
		*models.TeamEvent
		UserRSVP *models.TeamEventRSVP `json:"user_rsvp,omitempty"`
	}{event, userRSVP}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/v1/clubs/{clubid}/teams/{teamid}/events
func handleCreateTeamEvent(w http.ResponseWriter, r *http.Request) {
	type CreateEventRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	user := extractUser(r)
	club, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}

	if !team.IsAdmin(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

	event, err := team.CreateEvent(req.Name, req.Description, req.Location, startTime, endTime, user.ID)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// PUT /api/v1/clubs/{clubid}/teams/{teamid}/events/{eventid}
func handleUpdateTeamEvent(w http.ResponseWriter, r *http.Request) {
	type UpdateEventRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	user := extractUser(r)
	club, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}

	if !team.IsAdmin(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
		return
	}

	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

	event, err := team.UpdateEvent(eventID, req.Name, req.Description, req.Location, startTime, endTime, user.ID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DELETE /api/v1/clubs/{clubid}/teams/{teamid}/events/{eventid}
func handleDeleteTeamEvent(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	club, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}

	if !team.IsAdmin(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
		return
	}

	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	if err := team.DeleteEvent(eventID); err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/clubs/{clubid}/teams/{teamid}/events/upcoming
func handleGetTeamUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	_, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}
	events, err := team.GetUpcomingEvents()
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// POST /api/v1/clubs/{clubid}/teams/{teamid}/events/{eventid}/rsvp
func handleCreateOrUpdateTeamEventRSVP(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	_, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}
	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	type RSVPRequest struct {
		Response string `json:"response"`
	}
	var req RSVPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !team.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := user.CreateOrUpdateTeamEventRSVP(eventID, req.Response); err != nil {
		http.Error(w, "Failed to update RSVP", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/clubs/{clubid}/teams/{teamid}/events/{eventid}/rsvps
func handleGetTeamEventRSVPs(w http.ResponseWriter, r *http.Request) {
	_, team, ok := getClubAndTeam(r, w)
	if !ok {
		return
	}
	eventID := extractPathParam(r, "events")
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}
	rsvps, err := models.GetTeamEventRSVPs(eventID)
	if err != nil {
		http.Error(w, "Failed to get RSVPs", http.StatusInternalServerError)
		return
	}

	counts, err := models.GetTeamEventRSVPCounts(eventID)
	if err != nil {
		http.Error(w, "Failed to get RSVP counts", http.StatusInternalServerError)
		return
	}

	response := struct {
		RSVPs  []models.TeamEventRSVP `json:"rsvps"`
		Counts map[string]int         `json:"counts"`
	}{rsvps, counts}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
