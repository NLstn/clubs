package odata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

// TestRecurringEvents_ExpandRecurrence tests the recurring event expansion function
func TestRecurringEvents_ExpandRecurrence(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create test user and club
	user, _ := handlers.CreateTestUser(t, "user@example.com")
	club := handlers.CreateTestClub(t, user, "Test Club")

	// Initialize OData service
	service := &Service{
		db: database.Db,
	}

	t.Run("expand_weekly_recurring_event", func(t *testing.T) {
		// Create a recurring event (every week for 4 weeks)
		startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(2 * time.Hour)
		recurrenceEnd := startTime.AddDate(0, 0, 28) // 4 weeks

		pattern := "weekly"
		event := models.Event{
			ClubID:             club.ID,
			Name:               "Weekly Meeting",
			StartTime:          startTime,
			EndTime:            endTime,
			CreatedBy:          user.ID,
			UpdatedBy:          user.ID,
			IsRecurring:        true,
			RecurrencePattern:  &pattern,
			RecurrenceInterval: 1,
			RecurrenceEnd:      &recurrenceEnd,
		}

		err := database.Db.Create(&event).Error
		assert.NoError(t, err)

		// Expand for the entire period
		expandStart := startTime
		expandEnd := recurrenceEnd

		// Create request context with user ID
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(ctx)

		// Call the function
		params := map[string]interface{}{
			"startDate": expandStart,
			"endDate":   expandEnd,
		}

		result, err := service.expandRecurrenceFunction(httptest.NewRecorder(), req, &event, params)
		assert.NoError(t, err)

		// Verify results
		instances, ok := result.([]models.Event)
		assert.True(t, ok)
		assert.Len(t, instances, 5) // Parent + 4 weekly occurrences

		// Verify first instance is the parent event
		assert.Equal(t, event.ID, instances[0].ID)
		assert.Equal(t, event.StartTime, instances[0].StartTime)

		// Verify subsequent instances have correct dates
		for i := 1; i < len(instances); i++ {
			expectedStart := startTime.AddDate(0, 0, 7*i)
			assert.Equal(t, expectedStart, instances[i].StartTime, "Instance %d start time", i)

			// Verify duration is preserved
			duration := instances[i].EndTime.Sub(instances[i].StartTime)
			assert.Equal(t, 2*time.Hour, duration, "Instance %d duration", i)

			// Verify parent link
			assert.NotNil(t, instances[i].ParentEventID)
			assert.Equal(t, event.ID, *instances[i].ParentEventID)
		}
	})

	t.Run("expand_monthly_recurring_event", func(t *testing.T) {
		// Create a monthly recurring event (3 months)
		startTime := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.AddDate(0, 3, 0)

		pattern := "monthly"
		event := models.Event{
			ClubID:             club.ID,
			Name:               "Monthly Review",
			StartTime:          startTime,
			EndTime:            endTime,
			CreatedBy:          user.ID,
			UpdatedBy:          user.ID,
			IsRecurring:        true,
			RecurrencePattern:  &pattern,
			RecurrenceInterval: 1,
			RecurrenceEnd:      &recurrenceEnd,
		}

		err := database.Db.Create(&event).Error
		assert.NoError(t, err)

		// Expand
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		params := map[string]interface{}{
			"startDate": startTime,
			"endDate":   recurrenceEnd,
		}

		result, err := service.expandRecurrenceFunction(httptest.NewRecorder(), req, &event, params)
		assert.NoError(t, err)

		instances, ok := result.([]models.Event)
		assert.True(t, ok)
		assert.Len(t, instances, 4) // Parent + 3 monthly occurrences

		// Verify monthly intervals
		for i := 1; i < len(instances); i++ {
			expectedStart := startTime.AddDate(0, i, 0)
			assert.Equal(t, expectedStart, instances[i].StartTime, "Instance %d", i)
		}
	})

	t.Run("expand_non_recurring_event_in_range", func(t *testing.T) {
		// Create a non-recurring event
		startTime := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)

		event := models.Event{
			ClubID:      club.ID,
			Name:        "One-time Event",
			StartTime:   startTime,
			EndTime:     endTime,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
			IsRecurring: false,
		}

		err := database.Db.Create(&event).Error
		assert.NoError(t, err)

		// Expand with date range that includes the event
		expandStart := startTime.AddDate(0, 0, -7)
		expandEnd := startTime.AddDate(0, 0, 7)

		ctx := context.WithValue(context.Background(), auth.UserIDKey, user.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		params := map[string]interface{}{
			"startDate": expandStart,
			"endDate":   expandEnd,
		}

		result, err := service.expandRecurrenceFunction(httptest.NewRecorder(), req, &event, params)
		assert.NoError(t, err)

		instances, ok := result.([]models.Event)
		assert.True(t, ok)
		assert.Len(t, instances, 1) // Just the event itself
		assert.Equal(t, event.ID, instances[0].ID)
	})

	t.Run("expand_non_recurring_event_out_of_range", func(t *testing.T) {
		// Create a non-recurring event
		startTime := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)

		event := models.Event{
			ClubID:      club.ID,
			Name:        "Future Event",
			StartTime:   startTime,
			EndTime:     endTime,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
			IsRecurring: false,
		}

		err := database.Db.Create(&event).Error
		assert.NoError(t, err)

		// Expand with date range that excludes the event
		expandStart := startTime.AddDate(0, 0, -30)
		expandEnd := startTime.AddDate(0, 0, -7)

		ctx := context.WithValue(context.Background(), auth.UserIDKey, user.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		params := map[string]interface{}{
			"startDate": expandStart,
			"endDate":   expandEnd,
		}

		result, err := service.expandRecurrenceFunction(httptest.NewRecorder(), req, &event, params)
		assert.NoError(t, err)

		instances, ok := result.([]models.Event)
		assert.True(t, ok)
		assert.Len(t, instances, 0) // Event is outside range
	})

	t.Run("expand_unauthorized_user", func(t *testing.T) {
		// Create a recurring event
		startTime := time.Date(2024, 4, 1, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.AddDate(0, 0, 14)

		pattern := "weekly"
		event := models.Event{
			ClubID:             club.ID,
			Name:               "Private Meeting",
			StartTime:          startTime,
			EndTime:            endTime,
			CreatedBy:          user.ID,
			UpdatedBy:          user.ID,
			IsRecurring:        true,
			RecurrencePattern:  &pattern,
			RecurrenceInterval: 1,
			RecurrenceEnd:      &recurrenceEnd,
		}

		err := database.Db.Create(&event).Error
		assert.NoError(t, err)

		// Create another user who is NOT a member of the club
		otherUser, _ := handlers.CreateTestUser(t, "other@example.com")

		// Try to expand with unauthorized user
		ctx := context.WithValue(context.Background(), auth.UserIDKey, otherUser.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		params := map[string]interface{}{
			"startDate": startTime,
			"endDate":   recurrenceEnd,
		}

		_, err = service.expandRecurrenceFunction(httptest.NewRecorder(), req, &event, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}

// TestCalculateNextOccurrence tests the date calculation logic
func TestCalculateNextOccurrence(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	t.Run("daily_pattern", func(t *testing.T) {
		next := calculateNextOccurrence(baseTime, "daily", 1)
		expected := baseTime.AddDate(0, 0, 1)
		assert.Equal(t, expected, next)

		next = calculateNextOccurrence(baseTime, "daily", 3)
		expected = baseTime.AddDate(0, 0, 3)
		assert.Equal(t, expected, next)
	})

	t.Run("weekly_pattern", func(t *testing.T) {
		next := calculateNextOccurrence(baseTime, "weekly", 1)
		expected := baseTime.AddDate(0, 0, 7)
		assert.Equal(t, expected, next)

		next = calculateNextOccurrence(baseTime, "weekly", 2)
		expected = baseTime.AddDate(0, 0, 14)
		assert.Equal(t, expected, next)
	})

	t.Run("monthly_pattern", func(t *testing.T) {
		next := calculateNextOccurrence(baseTime, "monthly", 1)
		expected := baseTime.AddDate(0, 1, 0)
		assert.Equal(t, expected, next)

		next = calculateNextOccurrence(baseTime, "monthly", 3)
		expected = baseTime.AddDate(0, 3, 0)
		assert.Equal(t, expected, next)
	})

	t.Run("invalid_pattern", func(t *testing.T) {
		next := calculateNextOccurrence(baseTime, "yearly", 1)
		// Should return unchanged time for invalid pattern
		assert.Equal(t, baseTime, next)
	})
}
