package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestNotificationEndpoints(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Get Notifications - Unauthorized", func(t *testing.T) {
		req := MakeRequest(t, "GET", "/api/v1/notifications", nil, "")
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Get Notifications - Valid", func(t *testing.T) {
		_, token := CreateTestUser(t, "notifications@example.com")

		req := MakeRequest(t, "GET", "/api/v1/notifications", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var notifications []models.Notification
		ParseJSONResponse(t, rr, &notifications)
		assert.IsType(t, []models.Notification{}, notifications)
	})

	t.Run("Get Notification Count - Valid", func(t *testing.T) {
		_, token := CreateTestUser(t, "count@example.com")

		req := MakeRequest(t, "GET", "/api/v1/notifications/count", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		ParseJSONResponse(t, rr, &response)
		assert.Contains(t, response, "count")
	})

	t.Run("Get Notification Preferences - Valid", func(t *testing.T) {
		user, token := CreateTestUser(t, "preferences@example.com")

		req := MakeRequest(t, "GET", "/api/v1/notification-preferences", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var preferences models.UserNotificationPreferences
		ParseJSONResponse(t, rr, &preferences)
		assert.Equal(t, user.ID, preferences.UserID)
		assert.True(t, preferences.MemberAddedInApp) // Should create default preferences
	})

	t.Run("Update Notification Preferences - Valid", func(t *testing.T) {
		_, token := CreateTestUser(t, "update@example.com")

		updateData := map[string]interface{}{
			"memberAddedEmail":  false,
			"eventCreatedInApp": false,
		}

		req := MakeRequest(t, "PUT", "/api/v1/notification-preferences", updateData, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		var preferences models.UserNotificationPreferences
		ParseJSONResponse(t, rr, &preferences)
		assert.False(t, preferences.MemberAddedEmail)
		assert.False(t, preferences.EventCreatedInApp)
	})

	t.Run("Create and Mark Notification as Read", func(t *testing.T) {
		user, token := CreateTestUser(t, "readtest@example.com")

		// Create a test notification
		clubID := "test-club-id"
		err := models.CreateNotification(user.ID, "test_notification", "Test Title", "Test Message", &clubID, nil, nil)
		assert.NoError(t, err)

		// Get the notification
		notifications, err := models.GetUserNotifications(user.ID, 10)
		assert.NoError(t, err)
		assert.Greater(t, len(notifications), 0)

		notification := notifications[0]
		assert.False(t, notification.Read)

		// Mark as read
		req := MakeRequest(t, "PUT", "/api/v1/notifications/"+notification.ID, nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		// Verify it's marked as read
		updatedNotifications, err := models.GetUserNotifications(user.ID, 10)
		assert.NoError(t, err)

		found := false
		for _, n := range updatedNotifications {
			if n.ID == notification.ID {
				assert.True(t, n.Read)
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Mark All Notifications as Read", func(t *testing.T) {
		user, token := CreateTestUser(t, "markall@example.com")

		// Create multiple test notifications
		clubID := "test-club-id"
		err := models.CreateNotification(user.ID, "test1", "Test 1", "Message 1", &clubID, nil, nil)
		assert.NoError(t, err)
		err = models.CreateNotification(user.ID, "test2", "Test 2", "Message 2", &clubID, nil, nil)
		assert.NoError(t, err)

		// Mark all as read
		req := MakeRequest(t, "PUT", "/api/v1/notifications/mark-all-read", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusOK, rr.Code)

		// Verify all are marked as read
		count, err := models.GetUnreadNotificationCount(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		_, token := CreateTestUser(t, "method@example.com")

		req := MakeRequest(t, "POST", "/api/v1/notifications", nil, token)
		rr := ExecuteRequest(t, handler, req)
		CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
	})
}
