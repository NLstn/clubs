package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

// endpoint: GET /api/v1/clubs/{clubid}/events
func handleGetClubEvents(w http.ResponseWriter, r *http.Request) {

	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	userID := extractUserID(r)

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if !club.IsMember(userID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := models.GetClubEvents(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// endpoint: POST /api/v1/clubs/{clubid}/events
func handleCreateClubEvent(w http.ResponseWriter, r *http.Request) {

	userID := extractUserID(r)

	clubID := extractPathParam(r, "clubs")
	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	if !club.IsOwner(userID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !event.Validate() {
		http.Error(w, "Name, date, begin time, and end time are required", http.StatusBadRequest)
		return
	}

	if err := models.CreateEvent(&event, clubID); err != nil {
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

	userID := extractUserID(r)
	if !club.IsOwner(userID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Validate clubID as a UUID
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	// Validate eventID as a UUID
	if _, err := uuid.Parse(eventID); err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}

	if eventID == "" {
		http.Error(w, "Event ID parameter is required", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.DeleteEvent(eventID, clubID)
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
