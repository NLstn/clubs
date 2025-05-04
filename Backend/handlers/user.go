package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
)

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
		ID        string  `json:"id" gorm:"type:uuid;primary_key"`
		ClubID    string  `json:"clubId" gorm:"type:uuid"`
		ClubName  string  `json:"clubName"`
		Reason    string  `json:"reason"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
		Paid      bool    `json:"paid"`
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fines, err := user.GetFines()
	if err != nil {
		http.Error(w, "Failed to get fines", http.StatusInternalServerError)
		return
	}

	// load club names
	var result []Fine
	for i := range fines {
		club, err := models.GetClubByID(fines[i].ClubID)
		if err != nil {
			http.Error(w, "Failed to get club", http.StatusInternalServerError)
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

		result = append(result, fine)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
