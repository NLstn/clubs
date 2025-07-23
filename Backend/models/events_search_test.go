package models_test

import (
	"testing"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSearchEventsForUser(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create test data
	userID := uuid.New().String()
	club1ID := uuid.New().String()
	club2ID := uuid.New().String()

	// Create test user
	user := models.User{
		ID:         userID,
		FirstName:  "Test",
		LastName:   "User",
		Email:      "test@example.com",
		KeycloakID: userID + "-keycloak", // Make it unique
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := database.Db.Create(&user).Error
	assert.NoError(t, err)

	// Create test clubs
	club1 := models.Club{
		ID:          club1ID,
		Name:        "Test Club 1",
		Description: "First test club",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	club2 := models.Club{
		ID:          club2ID,
		Name:        "Test Club 2",
		Description: "Second test club",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	err = database.Db.Create(&club1).Error
	assert.NoError(t, err)
	err = database.Db.Create(&club2).Error
	assert.NoError(t, err)

	// Add user as member of club1 only
	member := models.Member{
		ID:        uuid.New().String(),
		UserID:    userID,
		ClubID:    club1ID,
		Role:      "member",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	err = database.Db.Create(&member).Error
	assert.NoError(t, err)

	// Create events in both clubs
	event1 := models.Event{
		ID:        uuid.New().String(),
		ClubID:    club1ID,
		Name:      "Search Test Event",
		StartTime: time.Now().Add(24 * time.Hour),
		EndTime:   time.Now().Add(26 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	event2 := models.Event{
		ID:        uuid.New().String(),
		ClubID:    club1ID,
		Name:      "Another Event",
		StartTime: time.Now().Add(48 * time.Hour),
		EndTime:   time.Now().Add(50 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	event3 := models.Event{
		ID:        uuid.New().String(),
		ClubID:    club2ID,
		Name:      "Search Private Event",
		StartTime: time.Now().Add(72 * time.Hour),
		EndTime:   time.Now().Add(74 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err = database.Db.Create(&event1).Error
	assert.NoError(t, err)
	err = database.Db.Create(&event2).Error
	assert.NoError(t, err)
	err = database.Db.Create(&event3).Error
	assert.NoError(t, err)

	t.Run("SearchEventsUserIsMemberOf", func(t *testing.T) {
		results, err := models.SearchEventsForUser(userID, "Search")
		assert.NoError(t, err)

		// Should find event1 (user is member of club1)
		// Should NOT find event3 (user is not member of club2)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, event1.ID, results[0].ID)
		assert.Equal(t, "Search Test Event", results[0].Name)
		assert.Equal(t, club1ID, results[0].ClubID)
		assert.Equal(t, "Test Club 1", results[0].ClubName)
	})

	t.Run("SearchEventsNoMatch", func(t *testing.T) {
		results, err := models.SearchEventsForUser(userID, "NonExistent")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("SearchEventsCaseInsensitive", func(t *testing.T) {
		results, err := models.SearchEventsForUser(userID, "search")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, event1.ID, results[0].ID)
	})

	t.Run("SearchEventsPartialMatch", func(t *testing.T) {
		results, err := models.SearchEventsForUser(userID, "Test")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, event1.ID, results[0].ID)
	})

	t.Run("SearchEventsUserNotMemberOfAnyClub", func(t *testing.T) {
		// Create a new user who is not a member of any club
		newUserID := uuid.New().String()
		newUser := models.User{
			ID:         newUserID,
			FirstName:  "New",
			LastName:   "User",
			Email:      "new@example.com",
			KeycloakID: newUserID + "-keycloak", // Make it unique
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := database.Db.Create(&newUser).Error
		assert.NoError(t, err)

		results, err := models.SearchEventsForUser(newUserID, "Search")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("SearchEventsFromDeletedClubs", func(t *testing.T) {
		// Mark club1 as deleted
		err := database.Db.Model(&club1).Update("deleted", true).Error
		assert.NoError(t, err)

		results, err := models.SearchEventsForUser(userID, "Search")
		assert.NoError(t, err)

		// Should not find events from deleted clubs
		assert.Equal(t, 0, len(results))

		// Clean up - mark club as not deleted for other tests
		err = database.Db.Model(&club1).Update("deleted", false).Error
		assert.NoError(t, err)
	})
}
