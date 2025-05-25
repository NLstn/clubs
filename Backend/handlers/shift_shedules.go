package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
)

func registerShiftRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/shifts", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetShifts(w, r)
		case http.MethodPost:
			handleCreateShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/shifts/{shiftid}/members", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetShiftMembers(w, r)
		case http.MethodPost:
			handleAddMemberToShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/shifts/{shiftid}/members/{memberid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleRemoveMemberFromShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// GET /api/v1/clubs/{clubid}/shifts
func handleGetShifts(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shifts, err := club.GetShifts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shifts); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// POST /api/v1/clubs/{clubid}/shifts
func handleCreateShift(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var shift models.Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if shift.StartTime.IsZero() || shift.EndTime.IsZero() {
		http.Error(w, "StartTime and EndTime are required", http.StatusBadRequest)
		return
	}

	shiftID, err := club.CreateShift(shift.StartTime.Time, shift.EndTime.Time)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": shiftID}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GET /api/v1/clubs/{clubid}/shifts/{shiftid}/members
func handleGetShiftMembers(w http.ResponseWriter, r *http.Request) {

	type ApiMember struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	clubID := extractPathParam(r, "clubs")
	shiftID := extractPathParam(r, "shifts")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(shiftID); err != nil {
		http.Error(w, "Invalid shift ID format", http.StatusBadRequest)
		return
	}

	members, err := models.GetShiftMembers(shiftID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var apiMembers []ApiMember
	for _, member := range members {
		user, err := models.GetUserByID(member.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		apiMembers = append(apiMembers, ApiMember{
			ID:   member.UserID,
			Name: user.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiMembers); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// POST /api/v1/clubs/{clubid}/shifts/{shiftid}/members
func handleAddMemberToShift(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	shiftID := extractPathParam(r, "shifts")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(shiftID); err != nil {
		http.Error(w, "Invalid shift ID format", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(requestBody.UserID); err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	if err := models.AddMemberToShift(shiftID, requestBody.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Member added to shift successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DELETE /api/v1/clubs/{clubid}/shifts/{shiftid}/members/{memberid}
func handleRemoveMemberFromShift(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	shiftID := extractPathParam(r, "shifts")
	memberID := extractPathParam(r, "members")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(shiftID); err != nil {
		http.Error(w, "Invalid shift ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(memberID); err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	if err := models.RemoveMemberFromShift(shiftID, memberID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Member removed from shift successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
