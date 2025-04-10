package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

// Helper function to extract path parameters
func extractPathParam(r *http.Request, param string) string {
	parts := strings.Split(r.URL.Path, "/")
	for i, part := range parts {
		if part == param && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// endpoint: /api/v1/clubs/{clubid}/events
func handleGetClubEvents(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	// Validate clubID as a UUID
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
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

func handleCreateClubEvent(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	// Validate clubID as a UUID
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
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
	eventID := extractPathParam(r, "events")
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
