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
	"github.com/NLstn/clubs/models/core"
	"github.com/google/uuid"
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
		event := core.Event{
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
		instances, ok := result.([]core.Event)
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
		event := core.Event{
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

		instances, ok := result.([]core.Event)
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

		event := core.Event{
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

		instances, ok := result.([]core.Event)
		assert.True(t, ok)
		assert.Len(t, instances, 1) // Just the event itself
		assert.Equal(t, event.ID, instances[0].ID)
	})

	t.Run("expand_non_recurring_event_out_of_range", func(t *testing.T) {
		// Create a non-recurring event
		startTime := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)

		event := core.Event{
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

		instances, ok := result.([]core.Event)
		assert.True(t, ok)
		assert.Len(t, instances, 0) // Event is outside range
	})

	t.Run("expand_unauthorized_user", func(t *testing.T) {
		// Create a recurring event
		startTime := time.Date(2024, 4, 1, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.AddDate(0, 0, 14)

		pattern := "weekly"
		event := core.Event{
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

// TestGetRSVPCounts tests the RSVP count aggregation function
func TestGetRSVPCounts(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create test user and club
	user, _ := handlers.CreateTestUser(t, "user@example.com")
	club := handlers.CreateTestClub(t, user, "Test Club")

	// Create additional users for RSVPs
	user2, _ := handlers.CreateTestUser(t, "user2@example.com")
	user3, _ := handlers.CreateTestUser(t, "user3@example.com")
	user4, _ := handlers.CreateTestUser(t, "user4@example.com")

	// Add additional users to club
	club.AddMember(user2.ID, "member")
	club.AddMember(user3.ID, "member")
	club.AddMember(user4.ID, "member")

	// Initialize OData service
	service := &Service{
		db: database.Db,
	}

	t.Run("counts_rsvps_by_response_type", func(t *testing.T) {
		// Create an event
		event1 := core.Event{
			ID:        uuid.New().String(),
			ClubID:    club.ID,
			Name:      "Test Event 1",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(2 * time.Hour),
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		}
		err := database.Db.Create(&event1).Error
		assert.NoError(t, err)

		// Create RSVPs: 2 Yes, 1 No, 1 Maybe
		rsvps := []core.EventRSVP{
			{EventID: event1.ID, UserID: user.ID, Response: "yes"},
			{EventID: event1.ID, UserID: user2.ID, Response: "yes"},
			{EventID: event1.ID, UserID: user3.ID, Response: "no"},
			{EventID: event1.ID, UserID: user4.ID, Response: "maybe"},
		}

		for _, rsvp := range rsvps {
			err := database.Db.Create(&rsvp).Error
			assert.NoError(t, err)
		}

		// Create request context with user ID
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(ctx)

		// Call the function
		result, err := service.getRSVPCountsFunction(httptest.NewRecorder(), req, &event1, map[string]interface{}{})
		assert.NoError(t, err)

		// Verify results
		counts, ok := result.(map[string]int64)
		assert.True(t, ok)
		assert.Equal(t, int64(2), counts["Yes"])
		assert.Equal(t, int64(1), counts["No"])
		assert.Equal(t, int64(1), counts["Maybe"])
	})

	t.Run("returns_zero_for_missing_response_types", func(t *testing.T) {
		// Create a new event with only "yes" responses
		event2 := core.Event{
			ID:        uuid.New().String(),
			ClubID:    club.ID,
			Name:      "Test Event 2",
			StartTime: time.Now().Add(24 * time.Hour),
			EndTime:   time.Now().Add(26 * time.Hour),
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		}
		err := database.Db.Create(&event2).Error
		assert.NoError(t, err)

		// Only one RSVP for this specific event
		rsvp := core.EventRSVP{
			EventID:  event2.ID,
			UserID:   user.ID,
			Response: "yes",
		}
		err = database.Db.Create(&rsvp).Error
		assert.NoError(t, err)

		// Create request context
		ctx := context.WithValue(context.Background(), auth.UserIDKey, user.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(ctx)

		// Call the function for event2
		result, err := service.getRSVPCountsFunction(httptest.NewRecorder(), req, &event2, map[string]interface{}{})
		assert.NoError(t, err)

		// Verify results - should have all response types even with zero count
		counts, ok := result.(map[string]int64)
		assert.True(t, ok)
		assert.Equal(t, int64(1), counts["Yes"])
		assert.Equal(t, int64(0), counts["No"])
		assert.Equal(t, int64(0), counts["Maybe"])
	})

	t.Run("unauthorized_user_cannot_get_counts", func(t *testing.T) {
		// Create event
		event3 := core.Event{
			ID:        uuid.New().String(),
			ClubID:    club.ID,
			Name:      "Test Event 3",
			StartTime: time.Now().Add(48 * time.Hour),
			EndTime:   time.Now().Add(50 * time.Hour),
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		}
		err := database.Db.Create(&event3).Error
		assert.NoError(t, err)

		// Create a user not in the club
		outsider, _ := handlers.CreateTestUser(t, "outsider@example.com")

		// Create request context with outsider user ID
		ctx := context.WithValue(context.Background(), auth.UserIDKey, outsider.ID)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(ctx)

		// Call the function - should fail
		_, err = service.getRSVPCountsFunction(httptest.NewRecorder(), req, &event3, map[string]interface{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}
