package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func TestGetDashboardNews(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create a test user
	user, _ := CreateTestUser(t, "test@example.com")

	// Create a test club
	club := CreateTestClub(t, user, "Test Club")

	// Create a news post
	newsID := uuid.New().String()
	news := models.News{
		ID:        newsID,
		ClubID:    club.ID,
		Title:     "Test News",
		Content:   "Test content",
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}
	if err := testDB.Create(&news).Error; err != nil {
		t.Fatalf("Failed to create news: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", "/api/v1/dashboard/news", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user context
	ctx := context.WithValue(req.Context(), auth.UserIDKey, user.ID)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handleGetDashboardNews(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var newsItems []NewsWithClub
	if err := json.Unmarshal(rr.Body.Bytes(), &newsItems); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(newsItems) != 1 {
		t.Errorf("Expected 1 news item, got %d", len(newsItems))
	}

	if newsItems[0].Title != "Test News" {
		t.Errorf("Expected news title 'Test News', got '%s'", newsItems[0].Title)
	}

	if newsItems[0].Club.Name != "Test Club" {
		t.Errorf("Expected club name 'Test Club', got '%s'", newsItems[0].Club.Name)
	}
}

func TestGetDashboardEvents(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create a test user
	user, _ := CreateTestUser(t, "test@example.com")

	// Create a test club
	club := CreateTestClub(t, user, "Test Club")

	// Create a future event
	eventID := uuid.New().String()
	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	
	event := models.Event{
		ID:        eventID,
		ClubID:    club.ID,
		Name:      "Test Event",
		StartTime: startTime,
		EndTime:   endTime,
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}
	if err := testDB.Create(&event).Error; err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", "/api/v1/dashboard/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user context
	ctx := context.WithValue(req.Context(), auth.UserIDKey, user.ID)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handleGetDashboardEvents(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var events []EventWithClub
	if err := json.Unmarshal(rr.Body.Bytes(), &events); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].Name != "Test Event" {
		t.Errorf("Expected event name 'Test Event', got '%s'", events[0].Name)
	}

	if events[0].Club.Name != "Test Club" {
		t.Errorf("Expected club name 'Test Club', got '%s'", events[0].Club.Name)
	}
}

func TestGetDashboardNewsNoClubs(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create a test user but no clubs
	user, _ := CreateTestUser(t, "test@example.com")

	// Create HTTP request
	req, err := http.NewRequest("GET", "/api/v1/dashboard/news", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user context
	ctx := context.WithValue(req.Context(), auth.UserIDKey, user.ID)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handleGetDashboardNews(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var news []NewsWithClub
	if err := json.Unmarshal(rr.Body.Bytes(), &news); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should return empty array when user has no clubs
	if len(news) != 0 {
		t.Errorf("Expected 0 news items, got %d", len(news))
	}
}

func TestGetDashboardEventsNoClubs(t *testing.T) {
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create a test user but no clubs
	user, _ := CreateTestUser(t, "test@example.com")

	// Create HTTP request
	req, err := http.NewRequest("GET", "/api/v1/dashboard/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user context
	ctx := context.WithValue(req.Context(), auth.UserIDKey, user.ID)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handleGetDashboardEvents(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var events []EventWithClub
	if err := json.Unmarshal(rr.Body.Bytes(), &events); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should return empty array when user has no clubs
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}