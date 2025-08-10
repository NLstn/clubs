package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDatabase holds the test database instance
var testDB *gorm.DB

// SetupTestDB initializes an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) {
	// Set test environment variable
	os.Setenv("GO_ENV", "test")

	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set the global database reference for the application
	database.Db = testDB

	// Set up SQLite-compatible tables
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS magic_links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT NOT NULL UNIQUE,
			keycloak_id TEXT UNIQUE,
			birth_date DATE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME,
			user_agent TEXT,
			ip_address TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS clubs (
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
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			role TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS join_requests (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fines (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			team_id TEXT,
			reason TEXT,
			amount REAL,
			paid BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fine_templates (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			description TEXT,
			amount REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS shifts (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			event_id TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS shift_members (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			shift_id TEXT NOT NULL,
			user_id TEXT NOT NULL
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			team_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			location TEXT,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT,
			is_recurring BOOLEAN DEFAULT FALSE,
			recurrence_pattern TEXT,
			recurrence_interval INTEGER DEFAULT 1,
			recurrence_end DATETIME,
			parent_event_id TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS event_rsvps (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			response TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS news (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS club_settings (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL UNIQUE,
			fines_enabled BOOLEAN DEFAULT TRUE,
			shifts_enabled BOOLEAN DEFAULT TRUE,
			teams_enabled BOOLEAN DEFAULT TRUE,
			members_list_visible BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			read BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			club_id TEXT,
			event_id TEXT,
			fine_id TEXT,
			invite_id TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS activities (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			actor_id TEXT,
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			metadata TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS user_notification_preferences (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			member_added_in_app BOOLEAN DEFAULT TRUE,
			member_added_email BOOLEAN DEFAULT TRUE,
			invite_received_in_app BOOLEAN DEFAULT TRUE,
			invite_received_email BOOLEAN DEFAULT TRUE,
			event_created_in_app BOOLEAN DEFAULT TRUE,
			event_created_email BOOLEAN DEFAULT FALSE,
			fine_assigned_in_app BOOLEAN DEFAULT TRUE,
			fine_assigned_email BOOLEAN DEFAULT TRUE,
			news_created_in_app BOOLEAN DEFAULT TRUE,
			news_created_email BOOLEAN DEFAULT FALSE,
			role_changed_in_app BOOLEAN DEFAULT TRUE,
			role_changed_email BOOLEAN DEFAULT TRUE,
			join_request_in_app BOOLEAN DEFAULT TRUE,
			join_request_email BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS user_privacy_settings (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			club_id TEXT,
			share_birth_date BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			club_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT,
			deleted BOOLEAN DEFAULT FALSE,
			deleted_at DATETIME,
			deleted_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT DEFAULT 'member',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T) {
	if testDB != nil {
		// Clear all data from tables to ensure clean state
		// Only clear tables that exist
		testDB.Exec("DELETE FROM activities")
		testDB.Exec("DELETE FROM refresh_tokens")
		testDB.Exec("DELETE FROM magic_links")
		testDB.Exec("DELETE FROM user_notification_preferences")
		testDB.Exec("DELETE FROM user_privacy_settings")
		testDB.Exec("DELETE FROM notifications")
		testDB.Exec("DELETE FROM fines")
		testDB.Exec("DELETE FROM members")
		testDB.Exec("DELETE FROM events")
		testDB.Exec("DELETE FROM clubs")
		testDB.Exec("DELETE FROM users")

		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// CreateTestUser creates a test user and returns the user and access token
func CreateTestUser(t *testing.T, email string) (models.User, string) {
	// Generate a UUID-like string for SQLite
	userID := uuid.New().String()
	keycloakID := uuid.New().String() // Generate unique KeycloakID for test users

	// Create user directly in database
	user := models.User{
		ID:         userID,
		Email:      email,
		FirstName:  "Test",
		LastName:   "User",
		KeycloakID: keycloakID,
	}

	if err := testDB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	return user, accessToken
}

// CreateTestClub creates a test club with the given user as owner
func CreateTestClub(t *testing.T, user models.User, clubName string) models.Club {
	clubID := uuid.New().String()

	club := models.Club{
		ID:          clubID,
		Name:        clubName,
		Description: "Test club description",
	}

	if err := testDB.Create(&club).Error; err != nil {
		t.Fatalf("Failed to create test club: %v", err)
	}

	// Add the owner as a member with owner role
	memberID := uuid.New().String()
	member := models.Member{
		ID:     memberID,
		UserID: user.ID,
		ClubID: club.ID,
		Role:   "owner",
	}
	if err := testDB.Create(&member).Error; err != nil {
		t.Fatalf("Failed to add owner as member: %v", err)
	}

	return club
}

// MakeRequest creates an HTTP request for testing
func MakeRequest(t *testing.T, method, url string, body interface{}, token string) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req
}

// ExecuteRequest executes an HTTP request against the test server
func ExecuteRequest(t *testing.T, handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// CheckResponseCode verifies the HTTP response code
func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

// MockEnvironmentVariables sets up test environment variables
func MockEnvironmentVariables(t *testing.T) {
	// Set minimal environment variables for testing
	// We don't need real Azure/database configs for unit tests
	os.Setenv("DATABASE_URL", "localhost")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("DATABASE_USER", "test")
	os.Setenv("DATABASE_USER_PASSWORD", "test")
	os.Setenv("DATABASE_NAME", "test_clubs")
	os.Setenv("DATABASE_SSL_MODE", "disable")

	// Set Azure environment variables to avoid initialization errors
	os.Setenv("AZURE_CLIENT_ID", "test-client-id")
	os.Setenv("AZURE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("AZURE_TENANT_ID", "test-tenant-id")
	os.Setenv("ACS_CONNECTION_STRING", "test-connection-string")

	// Set JWT secret for token generation
	os.Setenv("JWT_SECRET", "test-secret")
	if err := auth.Init(); err != nil {
		t.Fatalf("Failed to initialize auth: %v", err)
	}
}

// CreateTestFine creates a test fine for a user in a club
func CreateTestFine(t *testing.T, user models.User, club models.Club, reason string, amount float64, paid bool) models.Fine {
	fineID := uuid.New().String()

	fine := models.Fine{
		ID:        fineID,
		UserID:    user.ID,
		ClubID:    club.ID,
		Reason:    reason,
		Amount:    amount,
		Paid:      paid,
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}

	if err := testDB.Create(&fine).Error; err != nil {
		t.Fatalf("Failed to create test fine: %v", err)
	}

	return fine
}

// CreateTestFineWithCreator creates a test fine with a specific creator
func CreateTestFineWithCreator(t *testing.T, user models.User, club models.Club, creator models.User, reason string, amount float64, paid bool) models.Fine {
	fineID := uuid.New().String()

	fine := models.Fine{
		ID:        fineID,
		UserID:    user.ID,
		ClubID:    club.ID,
		Reason:    reason,
		Amount:    amount,
		Paid:      paid,
		CreatedBy: creator.ID,
		UpdatedBy: creator.ID,
	}

	if err := testDB.Create(&fine).Error; err != nil {
		t.Fatalf("Failed to create test fine: %v", err)
	}

	return fine
}

// CreateTestMember creates a test member directly in the database without notifications
func CreateTestMember(t *testing.T, user models.User, club models.Club, role string) models.Member {
	memberID := uuid.New().String()

	member := models.Member{
		ID:        memberID,
		UserID:    user.ID,
		ClubID:    club.ID,
		Role:      role,
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}

	if err := testDB.Create(&member).Error; err != nil {
		t.Fatalf("Failed to create test member: %v", err)
	}

	return member
}

// ParseJSONResponse parses a JSON response body into the provided interface
func ParseJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, v interface{}) {
	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, str, substr string) {
	if !bytes.Contains([]byte(str), []byte(substr)) {
		t.Errorf("Expected '%s' to contain '%s'", str, substr)
	}
}

// GetTestHandler returns a test HTTP handler with minimal middleware
func GetTestHandler() http.Handler {
	// Return the handler without rate limiting for tests
	mux := http.NewServeMux()

	// Register routes without rate limiting middleware for testing
	registerAuthRoutesForTest(mux)
	registerClubRoutesForTest(mux)
	registerClubSettingsRoutesForTest(mux)
	registerUserRoutesForTest(mux)
	registerMemberRoutesForTest(mux)
	registerTeamRoutesForTest(mux)
	registerShiftRoutesForTest(mux)
	registerJoinRequestRoutesForTest(mux)
	registerFineRoutesForTest(mux)
	registerNotificationRoutesForTest(mux)
	registerInviteRoutesForTest(mux)

	return mux
}

// Test route registration functions without rate limiting
func registerAuthRoutesForTest(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/auth/requestMagicLink", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleRequestMagicLinkForTest(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/auth/verifyMagicLink", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			verifyMagicLink(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/auth/refreshToken", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleRefreshToken(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleLogout(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// Test version of handleRequestMagicLink that doesn't send emails
func handleRequestMagicLinkForTest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	_, err := models.CreateMagicLink(req.Email)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// Skip sending email in tests
	w.WriteHeader(http.StatusNoContent)
}

func registerClubRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAllClubs(w, r)
		case http.MethodPost:
			handleCreateClub(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Handle club by ID and club fines - this pattern should match /api/v1/clubs/{id} and /api/v1/clubs/{id}/fines
	mux.Handle("/api/v1/clubs/", withAuth(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a fines endpoint
		if strings.HasSuffix(r.URL.Path, "/fines") {
			switch r.Method {
			case http.MethodGet:
				handleGetFines(w, r)
			case http.MethodPost:
				handleCreateFine(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Check if this is an individual fine endpoint (fines/{fineid})
		if strings.Contains(r.URL.Path, "/fines/") && !strings.Contains(r.URL.Path, "/fine-templates") {
			switch r.Method {
			case http.MethodDelete:
				handleDeleteFine(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Check if this is a fine-templates endpoint
		if strings.Contains(r.URL.Path, "/fine-templates") {
			if strings.Contains(r.URL.Path, "/fine-templates/") {
				// This is for specific template operations (PUT/DELETE)
				switch r.Method {
				case http.MethodPut:
					handleUpdateFineTemplate(w, r)
				case http.MethodDelete:
					handleDeleteFineTemplate(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			} else {
				// This is for general template operations (GET/POST)
				switch r.Method {
				case http.MethodGet:
					handleGetFineTemplates(w, r)
				case http.MethodPost:
					handleCreateFineTemplate(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}
			return
		}

		// Handle regular club operations
		switch r.Method {
		case http.MethodGet:
			handleGetClubByID(w, r)
		case http.MethodPatch:
			handleUpdateClub(w, r)
		case http.MethodDelete:
			// Check if this is hard delete
			if strings.HasSuffix(r.URL.Path, "/hard-delete") {
				handleHardDeleteClub(w, r)
			} else {
				handleDeleteClub(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func registerUserRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/me", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetMe(w, r)
		case http.MethodPut:
			handleUpdateMe(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/me/fines", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetMyFines(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

// Placeholder functions for other route registrations
func registerMemberRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/members", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubMembers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/isAdmin", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleCheckAdminRights(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/ownerCount", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetOwnerCount(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/members/{memberid}", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleClubMemberDelete(w, r)
		case http.MethodPatch:
			handleUpdateMemberRole(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/leave", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleLeaveClub(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func registerShiftRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/shifts", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetShifts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/events/{eventid}/shifts", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetEventShifts(w, r)
		case http.MethodPost:
			handleCreateEventShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/shifts/{shiftid}/members", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetShiftMembers(w, r)
		case http.MethodPost:
			handleAddMemberToShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/shifts/{shiftid}/members/{memberid}", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleRemoveMemberFromShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func registerJoinRequestRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/joinRequests", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetJoinRequests(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/inviteLink", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetInviteLink(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/join", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleJoinClubViaLink(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/info", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetClubInfo(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/joinRequests/{requestid}/accept", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleAcceptJoinRequest(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/joinRequests/{requestid}/reject", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleRejectJoinRequest(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func registerFineRoutesForTest(mux *http.ServeMux) {
	// Fines routes are handled in registerClubRoutesForTest and registerUserRoutesForTest
	// This function is left empty to avoid conflicts
}

func registerClubSettingsRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/settings", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubSettings(w, r)
		case http.MethodPost:
			handleUpdateClubSettings(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func registerNotificationRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/notifications", withAuth(GetNotifications))
	mux.Handle("/api/v1/notifications/count", withAuth(GetNotificationCount))
	mux.Handle("/api/v1/notifications/", withAuth(handleNotificationByID))
	mux.Handle("/api/v1/notifications/mark-all-read", withAuth(MarkAllNotificationsRead))
	mux.Handle("/api/v1/notification-preferences", withAuth(handleNotificationPreferences))
}

func registerInviteRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/invites", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateInvite(w, r)
		case http.MethodGet:
			handleGetClubInvites(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/invites", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetUserInvites(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/invites/{inviteid}/accept", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleAcceptInvite(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/invites/{inviteid}/reject", withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleRejectInvite(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

// Test route registration for team endpoints
func registerTeamRoutesForTest(mux *http.ServeMux) {
	// Register specific team routes
	mux.Handle("/api/v1/clubs/{clubid}/teams", withAuth(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request for user teams via query param
		if userID := r.URL.Query().Get("user"); userID != "" {
			// Handle get user teams
			handleGetUserTeams(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetClubTeams(w, r)
		case http.MethodPost:
			handleCreateTeam(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeam(w, r)
		case http.MethodPut:
			handleUpdateTeam(w, r)
		case http.MethodDelete:
			handleDeleteTeam(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/members", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTeamMembers(w, r)
		case http.MethodPost:
			handleAddTeamMember(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/clubs/{clubid}/teams/{teamid}/members/{memberid}", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handleUpdateTeamMemberRole(w, r)
		case http.MethodDelete:
			handleRemoveTeamMember(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}
