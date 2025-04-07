package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func handleGetClubEvents(w http.ResponseWriter, r *http.Request, clubID string) {
	var events []models.Event
	if result := database.Db.Where("club_id = ?", clubID).Find(&events); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func handleCreateClubEvent(w http.ResponseWriter, r *http.Request, clubID string) {
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

func handleDeleteClubEvent(w http.ResponseWriter, r *http.Request, clubID, eventID string) {
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
