package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

// endpoint: GET /api/v1/clubs
func handleGetAllClubs(w http.ResponseWriter, r *http.Request) {

	user := extractUser(r)

	clubs, err := models.GetAllClubs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var authorizedClubs []models.Club
	for _, club := range clubs {
		if club.IsMember(user) {
			authorizedClubs = append(authorizedClubs, club)
		}
	}

	clubs = authorizedClubs

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubs)
}

// endpoint: GET /api/v1/clubs/{clubid}
func handleGetClubByID(w http.ResponseWriter, r *http.Request) {

	user := extractUser(r)

	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(club)
}

// endpoint: POST /api/v1/clubs
func handleCreateClub(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	user := extractUser(r)

	var payload Body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	club, err := models.CreateClub(payload.Name, payload.Description, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(club)
}

// endpoint: PATCH /api/v1/clubs/{clubid}
func handleUpdateClub(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := club.Update(payload.Name, payload.Description); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(club)
}
