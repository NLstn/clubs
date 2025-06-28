package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
)

func registerNotificationSettingsRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/notificationSettings", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNotificationSettings(w, r)
		case http.MethodPost:
			handleUpdateNotificationSettings(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// GET /api/v1/notificationSettings
func handleGetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	settings, err := models.GetUserNotificationSettings(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// POST /api/v1/notificationSettings
func handleUpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	var payload []models.NotificationSetting
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i := range payload {
		payload[i].UserID = user.ID
		if err := models.UpsertNotificationSetting(&payload[i]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
