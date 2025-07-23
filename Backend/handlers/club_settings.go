package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

func registerClubSettingsRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/settings", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubSettings(w, r)
		case http.MethodPost:
			handleUpdateClubSettings(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/clubs/{clubid}/settings
func handleGetClubSettings(w http.ResponseWriter, r *http.Request) {
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

	if !club.IsAdmin(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	settings, err := models.GetClubSettings(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// endpoint: POST /api/v1/clubs/{clubid}/settings
func handleUpdateClubSettings(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		FinesEnabled  bool `json:"finesEnabled"`
		ShiftsEnabled bool `json:"shiftsEnabled"`
		TeamsEnabled  bool `json:"teamsEnabled"`
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

	if !club.IsAdmin(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := models.GetClubSettings(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := settings.Update(payload.FinesEnabled, payload.ShiftsEnabled, payload.TeamsEnabled, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
