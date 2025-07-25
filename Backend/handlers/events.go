package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerEventRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/events", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetEvents(w, r)
		case http.MethodPost:
			handleCreateEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/events/{eventid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetEvent(w, r)
		case http.MethodPut:
			handleUpdateEvent(w, r)
		case http.MethodDelete:
			handleDeleteEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/events/upcoming", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetUpcomingEvents(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/events/{eventid}/rsvp", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleCreateOrUpdateRSVP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/events/{eventid}/rsvps", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetEventRSVPs(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// GET /api/v1/clubs/{clubid}/events
func handleGetEvents(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := club.GetEvents()
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// GET /api/v1/clubs/{clubid}/events/{eventid}
func handleGetEvent(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	eventID := extractPathParam(r, "events")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	event, err := club.GetEventByID(eventID)
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

// POST /api/v1/clubs/{clubid}/events
func handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	type CreateEventRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
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

	event, err := club.CreateEvent(req.Name, req.Description, req.Location, startTime, endTime, user.ID)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// PUT /api/v1/clubs/{clubid}/events/{eventid}
func handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	type UpdateEventRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	eventID := extractPathParam(r, "events")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
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

	event, err := club.UpdateEvent(eventID, req.Name, req.Description, req.Location, startTime, endTime, user.ID)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DELETE /api/v1/clubs/{clubid}/events/{eventid}
func handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	eventID := extractPathParam(r, "events")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
		return
	}

	err = club.DeleteEvent(eventID)
	if err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/clubs/{clubid}/events/upcoming
func handleGetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := club.GetUpcomingEvents()
	if err != nil {
		http.Error(w, "Failed to get upcoming events", http.StatusInternalServerError)
		return
	}

	// Add user's RSVP status to each event
	type EventWithRSVP struct {
		models.Event
		UserRSVP *models.EventRSVP `json:"user_rsvp,omitempty"`
	}

	var eventsWithRSVP []EventWithRSVP
	for _, event := range events {
		eventWithRSVP := EventWithRSVP{Event: event}

		userRSVP, err := user.GetUserRSVP(event.ID)
		if err == nil {
			eventWithRSVP.UserRSVP = userRSVP
		}

		eventsWithRSVP = append(eventsWithRSVP, eventWithRSVP)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventsWithRSVP)
}

// POST /api/v1/clubs/{clubid}/events/{eventid}/rsvp
func handleCreateOrUpdateRSVP(w http.ResponseWriter, r *http.Request) {
	type RSVPRequest struct {
		Response string `json:"response"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	eventID := extractPathParam(r, "events")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Verify event exists and belongs to this club
	_, err = club.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	var req RSVPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Response != "yes" && req.Response != "no" && req.Response != "maybe" {
		http.Error(w, "Response must be 'yes', 'no', or 'maybe'", http.StatusBadRequest)
		return
	}

	err = user.CreateOrUpdateRSVP(eventID, req.Response)
	if err != nil {
		http.Error(w, "Failed to update RSVP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GET /api/v1/clubs/{clubid}/events/{eventid}/rsvps
func handleGetEventRSVPs(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	eventID := extractPathParam(r, "events")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
		return
	}

	// Verify event exists and belongs to this club
	_, err = club.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	counts, err := models.GetEventRSVPCounts(eventID)
	if err != nil {
		http.Error(w, "Failed to get RSVP counts", http.StatusInternalServerError)
		return
	}

	rsvps, err := models.GetEventRSVPs(eventID)
	if err != nil {
		http.Error(w, "Failed to get RSVPs", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"counts": counts,
		"rsvps":  rsvps,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
