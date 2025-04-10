package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

// handleGetAllClubs retrieves all clubs
func handleGetAllClubs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	clubs, err := models.GetAllClubs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var authorizedClubs []models.Club
	for _, club := range clubs {
		if auth.IsAuthorizedForClub(userID, club.ID) {
			authorizedClubs = append(authorizedClubs, club)
		}
	}

	clubs = authorizedClubs

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubs)
}

// handleGetClubByID retrieves a specific club by ID
func handleGetClubByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Extract club ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/clubs/")
	id := path

	club, err := models.GetClubByID(id)

	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !auth.IsAuthorizedForClub(userID, club.ID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(club)
}

// handleCreateClub creates a new club
func handleCreateClub(w http.ResponseWriter, r *http.Request) {
	var club models.Club
	if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ownerID := r.Context().Value(auth.UserIDKey).(string)

	if err := models.CreateClub(&club, ownerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(club)
}
