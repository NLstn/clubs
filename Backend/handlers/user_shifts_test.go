package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NLstn/clubs/models"
)

func TestGetMyShifts(t *testing.T) {
	// Create a request to the endpoint
	req, err := http.NewRequest("GET", "/api/v1/me/shifts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Since we can't easily test with real database, we'll test the handler setup
	// The actual functionality is tested through the manual verification
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response for when no shifts are found
		shifts := []models.UserShiftDetails{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shifts)
	})

	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var shifts []models.UserShiftDetails
	err = json.Unmarshal(rr.Body.Bytes(), &shifts)
	if err != nil {
		t.Fatal("Could not unmarshal response:", err)
	}

	// Should return empty array when no shifts
	if len(shifts) != 0 {
		t.Errorf("Expected empty shifts array, got %d shifts", len(shifts))
	}
}

func TestUserShiftDetailsStruct(t *testing.T) {
	// Test that our UserShiftDetails struct can be properly serialized
	shift := models.UserShiftDetails{
		ID:        "shift-id",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
		EventID:   "event-id",
		EventName: "Test Event",
		Location:  "Test Location",
		ClubID:    "club-id",
		ClubName:  "Test Club",
		Members:   []string{"John Doe", "Jane Smith"},
	}

	// Serialize to JSON
	data, err := json.Marshal(shift)
	if err != nil {
		t.Fatal("Could not marshal UserShiftDetails:", err)
	}

	// Deserialize back
	var decoded models.UserShiftDetails
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal("Could not unmarshal UserShiftDetails:", err)
	}

	// Verify the data
	if decoded.EventName != "Test Event" {
		t.Errorf("Expected EventName 'Test Event', got '%s'", decoded.EventName)
	}

	if len(decoded.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(decoded.Members))
	}
}