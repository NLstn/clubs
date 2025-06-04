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
			name TEXT,
			email TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS clubs (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
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
			email TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fines (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
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
			start_time DATETIME,
			end_time DATETIME,
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
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T) {
	if testDB != nil {
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

	// Create user directly in database
	user := models.User{
		ID:    userID,
		Email: email,
		Name:  "Test User",
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

	// Set Azure environment variables to avoid initialization errors
	os.Setenv("AZURE_CLIENT_ID", "test-client-id")
	os.Setenv("AZURE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("AZURE_TENANT_ID", "test-tenant-id")
	os.Setenv("ACS_CONNECTION_STRING", "test-connection-string")
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
	registerUserRoutesForTest(mux)
	registerMemberRoutesForTest(mux)
	registerShiftRoutesForTest(mux)
	registerJoinRequestRoutesForTest(mux)
	registerFineRoutesForTest(mux)

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
}

func registerShiftRoutesForTest(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/shifts", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetShifts(w, r)
		case http.MethodPost:
			handleCreateShift(w, r)
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
		switch r.Method {
		case http.MethodPost:
			handleJoinRequestCreate(w, r)
		case http.MethodGet:
			handleGetJoinEvents(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/api/v1/joinRequests", withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUserJoinRequests(w, r)
		default:
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
