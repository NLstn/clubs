package odata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// timelineTestContext holds test fixtures for timeline tests
type timelineTestContext struct {
	service   *Service
	handler   http.Handler
	user1     *core.User
	user2     *core.User
	club1     *core.Club
	club2     *core.Club
	member1   *core.Member
	member2   *core.Member
	token1    string
	token2    string
	activity1 *core.Activity
	event1    *core.Event
	news1     *core.News
	event2    *core.Event
}

// setupTimelineTestContext creates a test environment with timeline test data
func setupTimelineTestContext(t *testing.T) *timelineTestContext {
	// Set test environment
	os.Setenv("GO_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")

	// Initialize auth
	err := auth.Init()
	require.NoError(t, err, "Failed to initialize auth")

	// Set up unique in-memory SQLite database for each test
	// Using file::memory: with unique connection string for complete isolation
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err, "Failed to connect to test database")
	database.Db = testDB

	// Create tables manually with SQLite-compatible SQL to avoid UUID function issues
	testDB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		email TEXT NOT NULL UNIQUE,
		keycloak_id TEXT UNIQUE,
		birth_date DATE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS clubs (
		id TEXT PRIMARY KEY,
		name TEXT,
		description TEXT,
		logo_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT,
		deleted BOOLEAN DEFAULT FALSE,
		deleted_at DATETIME,
		deleted_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS members (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT NOT NULL,
		role TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS activities (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		user_id TEXT,
		actor_id TEXT,
		type TEXT,
		title TEXT,
		content TEXT,
		metadata TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		team_id TEXT,
		name TEXT,
		description TEXT,
		start_time DATETIME,
		end_time DATETIME,
		location TEXT,
		cancelled BOOLEAN DEFAULT FALSE,
		rsvp_enabled BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS news (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		title TEXT,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	)`)

	testDB.Exec(`CREATE TABLE IF NOT EXISTS event_rsvps (
		id TEXT PRIMARY KEY,
		event_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		response TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	// Create test users
	user1 := &core.User{
		ID:        uuid.New().String(),
		Email:     "user1@example.com",
		FirstName: "User",
		LastName:  "One",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := &core.User{
		ID:        uuid.New().String(),
		Email:     "user2@example.com",
		FirstName: "User",
		LastName:  "Two",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, testDB.Create(user1).Error)
	require.NoError(t, testDB.Create(user2).Error)

	// Create test clubs
	club1 := &core.Club{
		ID:        uuid.New().String(),
		Name:      "Test Club 1",
		CreatedAt: time.Now(),
		CreatedBy: user1.ID,
		UpdatedAt: time.Now(),
		UpdatedBy: user1.ID,
		Deleted:   false,
	}
	club2 := &core.Club{
		ID:        uuid.New().String(),
		Name:      "Test Club 2",
		CreatedAt: time.Now(),
		CreatedBy: user2.ID,
		UpdatedAt: time.Now(),
		UpdatedBy: user2.ID,
		Deleted:   false,
	}
	require.NoError(t, testDB.Create(club1).Error)
	require.NoError(t, testDB.Create(club2).Error)

	// Create memberships
	member1 := &core.Member{
		ID:        uuid.New().String(),
		UserID:    user1.ID,
		ClubID:    club1.ID,
		Role:      "owner",
		CreatedAt: time.Now(),
		CreatedBy: user1.ID,
		UpdatedAt: time.Now(),
		UpdatedBy: user1.ID,
	}
	member2 := &core.Member{
		ID:        uuid.New().String(),
		UserID:    user2.ID,
		ClubID:    club2.ID,
		Role:      "member",
		CreatedAt: time.Now(),
		CreatedBy: user2.ID,
		UpdatedAt: time.Now(),
		UpdatedBy: user2.ID,
	}
	require.NoError(t, testDB.Create(member1).Error)
	require.NoError(t, testDB.Create(member2).Error)

	// Create test activity
	activity1 := &core.Activity{
		ID:        uuid.New().String(),
		ClubID:    club1.ID,
		Type:      "member_joined",
		Title:     "New member joined",
		Content:   "John Doe joined the club",
		ActorID:   &user1.ID,
		Metadata:  `{"action": "join"}`,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}
	require.NoError(t, testDB.Create(activity1).Error)

	// Create test event
	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	description := "Annual general meeting"
	location := "Main Hall"
	event1 := &core.Event{
		ID:          uuid.New().String(),
		ClubID:      club1.ID,
		Name:        "Annual Meeting",
		Description: &description,
		Location:    &location,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, testDB.Create(event1).Error)

	// Create RSVP for user1 to event1
	rsvp := &core.EventRSVP{
		ID:        uuid.New().String(),
		EventID:   event1.ID,
		UserID:    user1.ID,
		Response:  "yes",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, testDB.Create(rsvp).Error)

	// Create another event for club2
	event2 := &core.Event{
		ID:        uuid.New().String(),
		ClubID:    club2.ID,
		Name:      "Club 2 Event",
		StartTime: time.Now().Add(48 * time.Hour),
		EndTime:   time.Now().Add(50 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, testDB.Create(event2).Error)

	// Create test news
	news1 := &core.News{
		ID:        uuid.New().String(),
		ClubID:    club1.ID,
		Title:     "Important Announcement",
		Content:   "Please read this important message",
		CreatedAt: time.Now().Add(-30 * time.Minute),
		UpdatedAt: time.Now().Add(-30 * time.Minute),
	}
	require.NoError(t, testDB.Create(news1).Error)

	// Create OData service
	service, err := NewService(testDB)
	require.NoError(t, err, "Failed to create OData service")

	// Create JWT tokens
	token1, err := auth.GenerateAccessToken(user1.ID)
	require.NoError(t, err)
	token2, err := auth.GenerateAccessToken(user2.ID)
	require.NoError(t, err)

	// Wrap service with authentication middleware and path prefix stripping
	handler := http.StripPrefix("/api/v2", AuthMiddleware(auth.GetJWTSecret())(service))

	return &timelineTestContext{
		service:   service,
		handler:   handler,
		user1:     user1,
		user2:     user2,
		club1:     club1,
		club2:     club2,
		member1:   member1,
		member2:   member2,
		token1:    token1,
		token2:    token2,
		activity1: activity1,
		event1:    event1,
		news1:     news1,
		event2:    event2,
	}
}

// TestGetTimelineCollection_Success tests successful timeline collection retrieval
func TestGetTimelineCollection_Success(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v2/TimelineItems", nil)
	req.Header.Set("Authorization", "Bearer "+ctx.token1)

	// Execute request
	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response struct {
		Value []core.TimelineItem `json:"value"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// User1 is member of club1, should see 3 items (activity, event, news)
	assert.Len(t, response.Value, 3)

	// Verify items are sorted by timestamp (most recent first)
	// Note: Events use StartTime for timestamp, so future events appear first
	if len(response.Value) >= 3 {
		// Event is most recent (StartTime is 24h in future)
		assert.Equal(t, "event", response.Value[0].Type)
		assert.Equal(t, ctx.event1.Name, response.Value[0].Title)
		assert.NotNil(t, response.Value[0].UserRSVP)
		assert.Equal(t, "yes", response.Value[0].UserRSVP.Response)

		// News is next (30 min ago)
		assert.Equal(t, "news", response.Value[1].Type)
		assert.Equal(t, ctx.news1.Title, response.Value[1].Title)
		assert.Equal(t, ctx.club1.Name, response.Value[1].ClubName)

		// Activity is oldest (2 hours ago)
		assert.Equal(t, "activity", response.Value[2].Type)
		assert.Equal(t, ctx.activity1.Title, response.Value[2].Title)
	}
}

// TestGetTimelineCollection_Unauthorized tests unauthorized access
func TestGetTimelineCollection_Unauthorized(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Create request without token
	req := httptest.NewRequest(http.MethodGet, "/api/v2/TimelineItems", nil)

	// Execute request
	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	// Assert unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetTimelineCollection_NoClubs tests user with no club memberships
func TestGetTimelineCollection_NoClubs(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Create a user with no club memberships
	user3 := &core.User{
		ID:        uuid.New().String(),
		Email:     "user3@example.com",
		FirstName: "User",
		LastName:  "Three",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, database.Db.Create(user3).Error)

	token3, err := auth.GenerateAccessToken(user3.ID)
	require.NoError(t, err)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v2/TimelineItems", nil)
	req.Header.Set("Authorization", "Bearer "+token3)

	// Execute request
	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response struct {
		Value []core.TimelineItem `json:"value"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return empty array
	assert.Len(t, response.Value, 0)
}

// TestGetTimelineCollection_MultipleClubs tests user with multiple club memberships
func TestGetTimelineCollection_MultipleClubs(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Add user1 to club2 as well
	member := &core.Member{
		ID:        uuid.New().String(),
		UserID:    ctx.user1.ID,
		ClubID:    ctx.club2.ID,
		Role:      "member",
		CreatedAt: time.Now(),
		CreatedBy: ctx.user1.ID,
		UpdatedAt: time.Now(),
		UpdatedBy: ctx.user1.ID,
	}
	require.NoError(t, database.Db.Create(member).Error)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v2/TimelineItems", nil)
	req.Header.Set("Authorization", "Bearer "+ctx.token1)

	// Execute request
	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response struct {
		Value []core.TimelineItem `json:"value"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// User1 is now member of both clubs, should see 4 items
	// (activity, event, news from club1 + event from club2)
	assert.Len(t, response.Value, 4)

	// Verify both clubs are represented
	clubIDs := make(map[string]bool)
	for _, item := range response.Value {
		clubIDs[item.ClubID] = true
	}
	assert.True(t, clubIDs[ctx.club1.ID])
	assert.True(t, clubIDs[ctx.club2.ID])
}

// TestGetUserClubs tests the getUserClubs helper function
func TestGetUserClubs(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Test getting clubs for user1
	clubIDs, clubNameMap, err := ctx.service.getUserClubs(ctx.user1.ID)
	require.NoError(t, err)

	assert.Len(t, clubIDs, 1)
	assert.Contains(t, clubIDs, ctx.club1.ID)
	assert.Equal(t, ctx.club1.Name, clubNameMap[ctx.club1.ID])
}

// TestGetUserClubs_NoClubs tests getUserClubs with user who has no clubs
func TestGetUserClubs_NoClubs(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Create a user with no memberships
	user3 := &core.User{
		ID:        uuid.New().String(),
		Email:     "user3@example.com",
		FirstName: "User",
		LastName:  "Three",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, database.Db.Create(user3).Error)

	// Test getting clubs
	clubIDs, clubNameMap, err := ctx.service.getUserClubs(user3.ID)
	require.NoError(t, err)

	assert.Len(t, clubIDs, 0)
	assert.Len(t, clubNameMap, 0)
}

// TestTimelineRSVPBatchFetch tests that RSVPs are fetched in batch
func TestTimelineRSVPBatchFetch(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Create multiple events with RSVPs
	for i := 0; i < 5; i++ {
		event := &core.Event{
			ID:        uuid.New().String(),
			ClubID:    ctx.club1.ID,
			Name:      fmt.Sprintf("Event %d", i),
			StartTime: time.Now().Add(time.Duration(i+1) * 24 * time.Hour),
			EndTime:   time.Now().Add(time.Duration(i+1)*24*time.Hour + 2*time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		require.NoError(t, database.Db.Create(event).Error)

		rsvp := &core.EventRSVP{
			ID:        uuid.New().String(),
			EventID:   event.ID,
			UserID:    ctx.user1.ID,
			Response:  "yes",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		require.NoError(t, database.Db.Create(rsvp).Error)
	}

	// Fetch timeline
	req := httptest.NewRequest(http.MethodGet, "/api/v2/TimelineItems", nil)
	req.Header.Set("Authorization", "Bearer "+ctx.token1)

	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Value []core.TimelineItem `json:"value"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Count events with RSVPs
	eventsWithRSVP := 0
	for _, item := range response.Value {
		if item.Type == "event" && item.UserRSVP != nil {
			eventsWithRSVP++
		}
	}

	// Should have 6 events with RSVPs (original event1 + 5 new events)
	assert.Equal(t, 6, eventsWithRSVP)
}

// TestGetTimelineCollection_WithContext tests that context is properly passed
func TestGetTimelineCollection_WithContext(t *testing.T) {
	ctx := setupTimelineTestContext(t)

	// Create request with context
	req := httptest.NewRequest(http.MethodGet, "/api/v2/TimelineItems", nil)
	req.Header.Set("Authorization", "Bearer "+ctx.token1)

	// Add user ID to context (this is typically done by middleware)
	reqCtx := context.WithValue(req.Context(), auth.UserIDKey, ctx.user1.ID)
	req = req.WithContext(reqCtx)

	// Execute request
	w := httptest.NewRecorder()
	ctx.handler.ServeHTTP(w, req)

	// Should work with context
	assert.Equal(t, http.StatusOK, w.Code)
}
