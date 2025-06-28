package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
)

func registerNotificationRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/notifications", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNotifications(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/notifications/", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handleUpdateNotification(w, r)
		case http.MethodDelete:
			handleDeleteNotification(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// GET /api/v1/notifications
func handleGetNotifications(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	notis, err := models.GetNotifications(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notis)
}

// PUT /api/v1/notifications/{id}
func handleUpdateNotification(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	id := r.URL.Path[len("/api/v1/notifications/"):]
	var body struct {
		Read bool `json:"read"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := models.UpdateNotificationRead(id, user.ID, body.Read); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/notifications/{id}
func handleDeleteNotification(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	id := r.URL.Path[len("/api/v1/notifications/"):]
	if err := models.DeleteNotification(id, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
