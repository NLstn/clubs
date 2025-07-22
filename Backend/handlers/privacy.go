package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func registerPrivacyRoutes(mux *http.ServeMux) {
	// Privacy settings endpoints
	mux.Handle("/api/v1/me/privacy", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetPrivacySettings(w, r)
		case http.MethodPut:
			handleUpdatePrivacySettings(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/me/privacy/clubs", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetClubSpecificPrivacySettings(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/me/privacy
func handleGetPrivacySettings(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	settings, err := models.GetUserGlobalPrivacySettings(user.ID)
	if err != nil {
		http.Error(w, "Failed to get privacy settings", http.StatusInternalServerError)
		return
	}

	response := struct {
		ShareBirthDate bool `json:"shareBirthDate"`
	}{
		ShareBirthDate: settings.ShareBirthDate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// endpoint: PUT /api/v1/me/privacy
func handleUpdatePrivacySettings(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ShareBirthDate *bool   `json:"shareBirthDate"`
		ClubID         *string `json:"clubId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ShareBirthDate == nil {
		http.Error(w, "ShareBirthDate is required", http.StatusBadRequest)
		return
	}

	var clubID string
	if req.ClubID != nil {
		clubID = *req.ClubID
		// Validate club ID format if provided
		if clubID != "" {
			if _, err := uuid.Parse(clubID); err != nil {
				http.Error(w, "Invalid club ID format", http.StatusBadRequest)
				return
			}
		}
	}

	err := models.UpdateOrCreatePrivacySettings(user.ID, clubID, *req.ShareBirthDate)
	if err != nil {
		http.Error(w, "Failed to update privacy settings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/me/privacy/clubs
func handleGetClubSpecificPrivacySettings(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get global settings
	globalSettings, err := models.GetUserGlobalPrivacySettings(user.ID)
	if err != nil {
		http.Error(w, "Failed to get privacy settings", http.StatusInternalServerError)
		return
	}

	// Get club-specific settings
	clubSettings, err := models.GetUserClubSpecificPrivacySettings(user.ID)
	if err != nil {
		http.Error(w, "Failed to get club privacy settings", http.StatusInternalServerError)
		return
	}

	response := struct {
		Global struct {
			ShareBirthDate bool `json:"shareBirthDate"`
		} `json:"global"`
		Clubs []struct {
			ClubID         string `json:"clubId"`
			ShareBirthDate bool   `json:"shareBirthDate"`
		} `json:"clubs"`
	}{}

	response.Global.ShareBirthDate = globalSettings.ShareBirthDate

	for _, setting := range clubSettings {
		clubID := ""
		if setting.ClubID != nil {
			clubID = *setting.ClubID
		}
		response.Clubs = append(response.Clubs, struct {
			ClubID         string `json:"clubId"`
			ShareBirthDate bool   `json:"shareBirthDate"`
		}{
			ClubID:         clubID,
			ShareBirthDate: setting.ShareBirthDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
