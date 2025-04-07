package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/database"
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

	var events []models.Event
	if result := database.Db.Where("club_id = ?", clubID).Find(&events); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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

	if event.Name == "" || event.Date == "" || event.BeginTime == "" || event.EndTime == "" {
		http.Error(w, "Name, date, begin time, and end time are required", http.StatusBadRequest)
		return
	}

	event.ID = uuid.New().String()
	event.ClubID = clubID

	if result := database.Db.Create(&event); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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
	result := database.Db.Where("id = ? AND club_id = ?", eventID, clubID).Delete(&models.Event{})
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
