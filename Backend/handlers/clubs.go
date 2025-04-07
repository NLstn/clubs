package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// handleGetAllClubs retrieves all clubs
func handleGetAllClubs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var clubs []models.Club
	if result := database.Db.Find(&clubs); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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

	var club models.Club
	result := database.Db.First(&club, "id = ?", id)

	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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

	club.ID = uuid.New().String()
	club.OwnerID = r.Context().Value(auth.UserIDKey).(string)

	if result := database.Db.Create(&club); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(club)
}
