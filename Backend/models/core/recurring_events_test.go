package core_test

import (
	"testing"
	"time"

	"github.com/NLstn/clubs/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecurringEvent(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	user, _ := handlers.CreateTestUser(t, "eventcreator@example.com")
	club := handlers.CreateTestClub(t, user, "Test Club")

	t.Run("create weekly recurring event", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		recurrenceEnd := startTime.Add(30 * 24 * time.Hour) // 30 days from start

		events, err := club.CreateRecurringEvent(
			"Weekly Meeting",
			"Every Monday meeting",
			"Conference Room",
			startTime,
			endTime,
			"weekly",
			1,
			recurrenceEnd,
			user.ID,
		)

		assert.NoError(t, err)
		assert.True(t, len(events) > 1, "Should create multiple events")

		// Check parent event
		parentEvent := events[0]
		assert.True(t, parentEvent.IsRecurring)
		assert.NotNil(t, parentEvent.RecurrencePattern)
		assert.Equal(t, "weekly", *parentEvent.RecurrencePattern)
		assert.Equal(t, 1, parentEvent.RecurrenceInterval)
		assert.NotNil(t, parentEvent.RecurrenceEnd)
		assert.Nil(t, parentEvent.ParentEventID)

		// Check child events
		for i := 1; i < len(events); i++ {
			childEvent := events[i]
			assert.False(t, childEvent.IsRecurring)
			assert.Nil(t, childEvent.RecurrencePattern)
			assert.NotNil(t, childEvent.ParentEventID)
			assert.Equal(t, parentEvent.ID, *childEvent.ParentEventID)

			// Check time progression
			expectedStartTime := startTime.AddDate(0, 0, 7*i)
			assert.True(t, childEvent.StartTime.Equal(expectedStartTime) || childEvent.StartTime.After(expectedStartTime.Add(-time.Minute)))
		}
	})

	t.Run("create daily recurring event", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.Add(7 * 24 * time.Hour) // 7 days from start

		events, err := club.CreateRecurringEvent(
			"Daily Standup",
			"Daily team standup",
			"Online",
			startTime,
			endTime,
			"daily",
			1,
			recurrenceEnd,
			user.ID,
		)

		assert.NoError(t, err)
		assert.Equal(t, 8, len(events), "Should create 8 events (7 days + initial)")

		// Check that events are created daily
		for i := 1; i < len(events); i++ {
			expectedStartTime := startTime.AddDate(0, 0, i)
			assert.True(t, events[i].StartTime.Equal(expectedStartTime) || events[i].StartTime.After(expectedStartTime.Add(-time.Minute)))
		}
	})

	t.Run("invalid recurrence pattern", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.Add(7 * 24 * time.Hour)

		events, err := club.CreateRecurringEvent(
			"Invalid Pattern",
			"Test invalid pattern",
			"Location",
			startTime,
			endTime,
			"invalid",
			1,
			recurrenceEnd,
			user.ID,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported recurrence pattern")
		assert.Len(t, events, 1, "Should only create parent event before failing")
	})

	t.Run("invalid recurrence interval", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)
		recurrenceEnd := startTime.Add(7 * 24 * time.Hour)

		events, err := club.CreateRecurringEvent(
			"Invalid Interval",
			"Test invalid interval",
			"Location",
			startTime,
			endTime,
			"weekly",
			0, // Invalid interval
			recurrenceEnd,
			user.ID,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recurrence parameters")
		assert.Nil(t, events)
	})
}
