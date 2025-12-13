package odata_test

import (
	"testing"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/NLstn/clubs/odata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSearchGlobal(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create test user
	user, _ := handlers.CreateTestUser(t, "testuser@example.com")

	// Create test club with "Test Club" name
	club := handlers.CreateTestClub(t, user, "Test Club")

	// Create test event in club
	eventDesc := "Test event description"
	eventLoc := "Test location"
	event := models.Event{
		ID:          uuid.New().String(),
		ClubID:      club.ID,
		Name:        "Test Event",
		Description: &eventDesc,
		Location:    &eventLoc,
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(26 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err := database.Db.Create(&event).Error
	assert.NoError(t, err)

	// Create OData service
	service, err := odata.NewService(database.Db)
	assert.NoError(t, err)

	t.Run("SearchClubsByPartialName", func(t *testing.T) {
		// Search for "Test" should find "Test Club"
		results, err := service.SearchClubsForTest(user, "Test")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1, "Should find at least one club with 'Test' in name")

		// Verify we found the Test Club
		found := false
		for _, result := range results {
			if result.Name == "Test Club" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find 'Test Club' when searching for 'Test'")
	})

	t.Run("SearchEventsByPartialName", func(t *testing.T) {
		// Search for "Test" should find "Test Event"
		results, err := service.SearchEventsForTest(user, "Test")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1, "Should find at least one event with 'Test' in name")

		// Verify we found the Test Event
		found := false
		for _, result := range results {
			if result.Name == "Test Event" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find 'Test Event' when searching for 'Test'")
	})

	t.Run("SearchClubsCaseInsensitive", func(t *testing.T) {
		// Search with lowercase should still find "Test Club"
		results, err := service.SearchClubsForTest(user, "test")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1, "Should find club with case-insensitive search")
	})

	t.Run("SearchEventsOnlyUpcoming", func(t *testing.T) {
		// Create a past event
		pastEvent := models.Event{
			ID:          uuid.New().String(),
			ClubID:      club.ID,
			Name:        "Past Test Event",
			Description: &eventDesc,
			Location:    &eventLoc,
			StartTime:   time.Now().Add(-48 * time.Hour),
			EndTime:     time.Now().Add(-46 * time.Hour),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := database.Db.Create(&pastEvent).Error
		assert.NoError(t, err)

		// Search should only find upcoming events
		results, err := service.SearchEventsForTest(user, "Test")
		assert.NoError(t, err)

		// Verify past event is not in results
		for _, result := range results {
			assert.NotEqual(t, "Past Test Event", result.Name, "Past events should not be returned")
		}
	})
}
