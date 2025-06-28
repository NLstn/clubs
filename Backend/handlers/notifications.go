package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/NLstn/clubs/models"
)

// registerNotificationRoutes registers all notification-related routes
func registerNotificationRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/notifications", withAuth(GetNotifications))
	mux.Handle("/api/v1/notifications/count", withAuth(GetNotificationCount))
	mux.Handle("/api/v1/notifications/", withAuth(handleNotificationByID)) // for marking specific notifications as read
	mux.Handle("/api/v1/notifications/mark-all-read", withAuth(MarkAllNotificationsRead))
	mux.Handle("/api/v1/notification-preferences", withAuth(handleNotificationPreferences))
}

// handleNotificationByID routes to the appropriate handler based on the notification ID
func handleNotificationByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		MarkNotificationRead(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleNotificationPreferences routes to the appropriate handler for notification preferences
func handleNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetNotificationPreferences(w, r)
	case http.MethodPut:
		UpdateNotificationPreferences(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetNotifications handles GET requests to retrieve user notifications
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	notifications, err := models.GetUserNotifications(user.ID, limit)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// GetNotificationCount handles GET requests to retrieve unread notification count
func GetNotificationCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := models.GetUnreadNotificationCount(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch notification count", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MarkNotificationRead handles PUT requests to mark a notification as read
func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get notification ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/notifications/")
	if path == "" || path == "mark-all-read" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}
	
	notificationID := path

	err := models.MarkNotificationAsRead(notificationID, user.ID)
	if err != nil {
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// MarkAllNotificationsRead handles PUT requests to mark all notifications as read
func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := models.MarkAllNotificationsAsRead(user.ID)
	if err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetNotificationPreferences handles GET requests to retrieve user notification preferences
func GetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	preferences, err := models.GetUserNotificationPreferences(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch notification preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preferences)
}

// UpdateNotificationPreferences handles PUT requests to update user notification preferences
func UpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateData struct {
		MemberAddedInApp  *bool `json:"memberAddedInApp"`
		MemberAddedEmail  *bool `json:"memberAddedEmail"`
		EventCreatedInApp *bool `json:"eventCreatedInApp"`
		EventCreatedEmail *bool `json:"eventCreatedEmail"`
		FineAssignedInApp *bool `json:"fineAssignedInApp"`
		FineAssignedEmail *bool `json:"fineAssignedEmail"`
		NewsCreatedInApp  *bool `json:"newsCreatedInApp"`
		NewsCreatedEmail  *bool `json:"newsCreatedEmail"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	preferences, err := models.GetUserNotificationPreferences(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch notification preferences", http.StatusInternalServerError)
		return
	}

	// Update only the fields that were provided
	if updateData.MemberAddedInApp != nil {
		preferences.MemberAddedInApp = *updateData.MemberAddedInApp
	}
	if updateData.MemberAddedEmail != nil {
		preferences.MemberAddedEmail = *updateData.MemberAddedEmail
	}
	if updateData.EventCreatedInApp != nil {
		preferences.EventCreatedInApp = *updateData.EventCreatedInApp
	}
	if updateData.EventCreatedEmail != nil {
		preferences.EventCreatedEmail = *updateData.EventCreatedEmail
	}
	if updateData.FineAssignedInApp != nil {
		preferences.FineAssignedInApp = *updateData.FineAssignedInApp
	}
	if updateData.FineAssignedEmail != nil {
		preferences.FineAssignedEmail = *updateData.FineAssignedEmail
	}
	if updateData.NewsCreatedInApp != nil {
		preferences.NewsCreatedInApp = *updateData.NewsCreatedInApp
	}
	if updateData.NewsCreatedEmail != nil {
		preferences.NewsCreatedEmail = *updateData.NewsCreatedEmail
	}

	err = preferences.Update()
	if err != nil {
		http.Error(w, "Failed to update notification preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preferences)
}