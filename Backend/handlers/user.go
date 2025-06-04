package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/models"
)

func registerUserRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/me", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetMe(w, r)
		case http.MethodPut:
			handleUpdateMe(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/me
func handleGetMe(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// endpoint: POST /api/v1/me
func handleUpdateMe(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	if err := user.UpdateUserName(req.Name); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/me/fines
func handleGetMyFines(w http.ResponseWriter, r *http.Request) {

	type Fine struct {
		ID            string    `json:"id" gorm:"type:uuid;primary_key"`
		ClubID        string    `json:"clubId" gorm:"type:uuid"`
		ClubName      string    `json:"clubName"`
		Reason        string    `json:"reason"`
		Amount        float64   `json:"amount"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
		Paid          bool      `json:"paid"`
		CreatedByName string    `json:"createdByName"`
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fines, err := user.GetUnpaidFines()
	if err != nil {
		http.Error(w, "Failed to get fines", http.StatusInternalServerError)
		return
	}

	// If user has no fines, return empty array early
	if len(fines) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Fine{})
		return
	}

	// Extract unique club IDs and creator IDs for batch queries
	clubIDSet := make(map[string]bool)
	creatorIDSet := make(map[string]bool)
	for _, fine := range fines {
		if fine.ClubID != "" {
			clubIDSet[fine.ClubID] = true
		}
		if fine.CreatedBy != "" {
			creatorIDSet[fine.CreatedBy] = true
		}
	}

	// Convert sets to slices
	var clubIDs []string
	for clubID := range clubIDSet {
		clubIDs = append(clubIDs, clubID)
	}
	var creatorIDs []string
	for creatorID := range creatorIDSet {
		creatorIDs = append(creatorIDs, creatorID)
	}

	// Fetch all clubs and creators in bulk
	clubs, err := models.GetClubsByIDs(clubIDs)
	if err != nil {
		http.Error(w, "Failed to get clubs", http.StatusInternalServerError)
		return
	}

	creators, err := models.GetUsersByIDs(creatorIDs)
	if err != nil {
		http.Error(w, "Failed to get fine creators", http.StatusInternalServerError)
		return
	}

	// Create lookup maps for quick access
	clubMap := make(map[string]models.Club)
	for _, club := range clubs {
		clubMap[club.ID] = club
	}

	creatorMap := make(map[string]models.User)
	for _, creator := range creators {
		creatorMap[creator.ID] = creator
	}

	// Build response using cached data
	var result []Fine
	for i := range fines {
		club, clubExists := clubMap[fines[i].ClubID]
		if !clubExists {
			http.Error(w, "Club not found for fine", http.StatusInternalServerError)
			return
		}

		creator, creatorExists := creatorMap[fines[i].CreatedBy]
		if !creatorExists {
			http.Error(w, "Creator not found for fine", http.StatusInternalServerError)
			return
		}

		var fine Fine
		fine.ID = fines[i].ID
		fine.ClubID = fines[i].ClubID
		fine.Reason = fines[i].Reason
		fine.Amount = fines[i].Amount
		fine.CreatedAt = fines[i].CreatedAt
		fine.UpdatedAt = fines[i].UpdatedAt
		fine.Paid = fines[i].Paid
		fine.ClubName = club.Name
		fine.CreatedByName = creator.Name

		result = append(result, fine)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
