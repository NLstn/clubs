package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
)

// endpoint: GET /api/v1/clubs/{clubid}/events
func handleGetClubEvents(w http.ResponseWriter, r *http.Request) {

	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := club.GetEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// endpoint: POST /api/v1/clubs/{clubid}/events
func handleCreateClubEvent(w http.ResponseWriter, r *http.Request) {

	type Payload struct {
		ID          string `json:"id" gorm:"type:uuid;primary_key"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ClubID      string `json:"club_id" gorm:"type:uuid"`
		Date        string `json:"date"`
		BeginTime   string `json:"begin_time"`
		EndTime     string `json:"end_time"`
	}

	user := extractUser(r)

	clubID := extractPathParam(r, "clubs")
	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := club.CreateEvent(payload.Name, payload.Description, payload.Date, payload.BeginTime, payload.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// endpoint: /api/v1/clubs/{clubid}/events/{eventid}
func handleDeleteClubEvent(w http.ResponseWriter, r *http.Request) {

	clubID := extractPathParam(r, "clubs")
	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	eventID := extractPathParam(r, "events")
	if eventID == "" {
		http.Error(w, "Event ID parameter is required", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	rowsAffected, err := club.DeleteEvent(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
