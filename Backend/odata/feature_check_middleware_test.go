package odata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFeatureCheckTest(t *testing.T) (*gorm.DB, string, string, string) {
	// Create SQLite database with shared cache
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	// Set database for the models package
	database.Db = db

	// Create tables manually with SQLite-compatible SQL
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		email TEXT NOT NULL,
		keycloak_id TEXT,
		birth_date DATE,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS clubs (
		id TEXT PRIMARY KEY,
		name TEXT,
		description TEXT,
		logo_url TEXT,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT,
		deleted BOOLEAN DEFAULT FALSE,
		deleted_at DATETIME,
		deleted_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS members (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		club_id TEXT NOT NULL,
		role TEXT,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS club_settings (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL UNIQUE,
		fines_enabled BOOLEAN DEFAULT FALSE,
		shifts_enabled BOOLEAN DEFAULT FALSE,
		teams_enabled BOOLEAN DEFAULT FALSE,
		members_list_visible BOOLEAN DEFAULT FALSE,
		discoverable_by_non_members BOOLEAN DEFAULT FALSE,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS fines (
		id TEXT PRIMARY KEY,
		club_id TEXT,
		team_id TEXT,
		user_id TEXT,
		reason TEXT,
		amount REAL,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT,
		paid BOOLEAN DEFAULT FALSE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS fine_templates (
		id TEXT PRIMARY KEY,
		club_id TEXT,
		description TEXT,
		amount REAL,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		team_id TEXT,
		name TEXT NOT NULL,
		description TEXT,
		location TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS shifts (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		event_id TEXT NOT NULL,
		start_time DATETIME,
		end_time DATETIME,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS shift_members (
		id TEXT PRIMARY KEY,
		shift_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS teams (
		id TEXT PRIMARY KEY,
		club_id TEXT NOT NULL,
		name TEXT,
		description TEXT,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS team_members (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		role TEXT,
		created_at DATETIME,
		created_by TEXT,
		updated_at DATETIME,
		is_recurring BOOLEAN DEFAULT FALSE,
		recurrence_pattern TEXT,
		recurrence_interval INTEGER DEFAULT 1,
		recurrence_end DATETIME,
		parent_event_id TEXT,
		updated_by TEXT
	)`)

	// Create test user
	userID := uuid.New().String()
	err = db.Exec(`
		INSERT INTO users (id, email, first_name, last_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, userID, "test@example.com", "Test", "User").Error
	require.NoError(t, err)

	// Create test club
	clubID := uuid.New().String()
	club := models.Club{
		ID:        clubID,
		Name:      "Test Club",
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&club).Error)

	// Create member relationship
	memberID := uuid.New().String()
	member := models.Member{
		ID:        memberID,
		ClubID:    clubID,
		UserID:    userID,
		Role:      "owner",
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&member).Error)

	// Create club settings with all features disabled by default
	settingsID := uuid.New().String()
	settings := models.ClubSettings{
		ID:            settingsID,
		ClubID:        clubID,
		FinesEnabled:  false,
		ShiftsEnabled: false,
		TeamsEnabled:  false,
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}
	require.NoError(t, db.Create(&settings).Error)

	return db, userID, clubID, memberID
}

func TestFeatureCheckMiddleware_FinesDisabled(t *testing.T) {
	db, userID, clubID, _ := setupFeatureCheckTest(t)

	// Create a fine entity
	fineID := uuid.New().String()
	fine := models.Fine{
		ID:        fineID,
		ClubID:    clubID,
		UserID:    userID,
		Reason:    "Test fine",
		Amount:    10.0,
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&fine).Error)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler that should not be called
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request to access the fine
	req := httptest.NewRequest("GET", "/Fines('"+fineID+"')", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 when fines feature is disabled")
	assert.False(t, handlerCalled, "Handler should not be called when feature is disabled")
	assert.Contains(t, w.Body.String(), "fines feature is disabled", "Error message should mention disabled feature")
}

func TestFeatureCheckMiddleware_FinesEnabled(t *testing.T) {
	db, userID, clubID, _ := setupFeatureCheckTest(t)

	// Enable fines feature
	var settings models.ClubSettings
	require.NoError(t, db.Where("club_id = ?", clubID).First(&settings).Error)
	settings.FinesEnabled = true
	require.NoError(t, db.Save(&settings).Error)

	// Create a fine entity
	fineID := uuid.New().String()
	fine := models.Fine{
		ID:        fineID,
		ClubID:    clubID,
		UserID:    userID,
		Reason:    "Test fine",
		Amount:    10.0,
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&fine).Error)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request to access the fine
	req := httptest.NewRequest("GET", "/Fines('"+fineID+"')", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 when fines feature is enabled")
	assert.True(t, handlerCalled, "Handler should be called when feature is enabled")
}

func TestFeatureCheckMiddleware_TeamsDisabled(t *testing.T) {
	db, userID, clubID, _ := setupFeatureCheckTest(t)

	// Create a team entity
	teamID := uuid.New().String()
	team := models.Team{
		ID:        teamID,
		ClubID:    clubID,
		Name:      "Test Team",
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&team).Error)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request to access the team
	req := httptest.NewRequest("GET", "/Teams('"+teamID+"')", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 when teams feature is disabled")
	assert.False(t, handlerCalled, "Handler should not be called when feature is disabled")
	assert.Contains(t, w.Body.String(), "teams feature is disabled", "Error message should mention disabled feature")
}

func TestFeatureCheckMiddleware_ShiftsDisabled(t *testing.T) {
	db, userID, clubID, _ := setupFeatureCheckTest(t)

	// Create an event first (required for shifts)
	eventID := uuid.New().String()
	event := models.Event{
		ID:        eventID,
		ClubID:    clubID,
		Name:      "Test Event",
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&event).Error)

	// Create a shift entity
	shiftID := uuid.New().String()
	shift := models.Shift{
		ID:        shiftID,
		ClubID:    clubID,
		EventID:   eventID,
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	require.NoError(t, db.Create(&shift).Error)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request to access the shift
	req := httptest.NewRequest("GET", "/Shifts('"+shiftID+"')", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 when shifts feature is disabled")
	assert.False(t, handlerCalled, "Handler should not be called when feature is disabled")
	assert.Contains(t, w.Body.String(), "shifts feature is disabled", "Error message should mention disabled feature")
}

func TestFeatureCheckMiddleware_FineTemplateDisabled(t *testing.T) {
	db, userID, clubID, _ := setupFeatureCheckTest(t)

	// Create a fine template entity
	templateID := uuid.New().String()
	template := models.FineTemplate{
		ID:          templateID,
		ClubID:      clubID,
		Description: "Test template",
		Amount:      15.0,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	require.NoError(t, db.Create(&template).Error)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request to access the template
	req := httptest.NewRequest("GET", "/FineTemplates('"+templateID+"')", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 when fines feature is disabled")
	assert.False(t, handlerCalled, "Handler should not be called when feature is disabled")
}

func TestFeatureCheckMiddleware_NonFeatureEntity(t *testing.T) {
	_, userID, clubID, _ := setupFeatureCheckTest(t)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request to access a non-feature entity (e.g., Clubs)
	req := httptest.NewRequest("GET", "/Clubs('"+clubID+"')", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response - should pass through to handler
	assert.Equal(t, http.StatusOK, w.Code, "Should pass through for non-feature entities")
	assert.True(t, handlerCalled, "Handler should be called for non-feature entities")
}

func TestFeatureCheckMiddleware_CollectionQuery(t *testing.T) {
	_, userID, _, _ := setupFeatureCheckTest(t)

	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request for collection query (no specific entity ID)
	req := httptest.NewRequest("GET", "/Fines", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, userID))
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response - should pass through to handler (middleware can't check for collection queries)
	assert.Equal(t, http.StatusOK, w.Code, "Collection queries should pass through to OData hooks")
	assert.True(t, handlerCalled, "Handler should be called for collection queries")
}

func TestFeatureCheckMiddleware_Metadata(t *testing.T) {
	// Create middleware
	middleware := FeatureCheckMiddleware()

	// Create a mock handler
	handlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(mockHandler)

	// Create request for metadata
	req := httptest.NewRequest("GET", "/$metadata", nil)
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response - should pass through
	assert.Equal(t, http.StatusOK, w.Code, "Metadata requests should pass through")
	assert.True(t, handlerCalled, "Handler should be called for metadata requests")
}
