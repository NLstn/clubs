package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
			description TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			role TEXT
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			club_id TEXT NOT NULL,
			reason TEXT,
			amount REAL,
			paid BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testDB.Exec(`
		CREATE TABLE IF NOT EXISTS shifts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			club_id TEXT NOT NULL,
			start_time DATETIME,
			end_time DATETIME,
			description TEXT
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

	// Handle club by ID - this pattern should match /api/v1/clubs/{id}
	mux.Handle("/api/v1/clubs/", withAuth(func(w http.ResponseWriter, r *http.Request) {
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
}

// Placeholder functions for other route registrations
func registerMemberRoutesForTest(mux *http.ServeMux) {
	// Skip member routes to avoid conflicts with club routes
	// Member functionality can be tested through integration tests
}

func registerShiftRoutesForTest(mux *http.ServeMux) {
	// TODO: Add shift routes when implementing those tests
}

func registerJoinRequestRoutesForTest(mux *http.ServeMux) {
	// TODO: Add join request routes when implementing those tests
}

func registerFineRoutesForTest(mux *http.ServeMux) {
	// TODO: Add fine routes when implementing those tests
}