package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NLstn/clubs/models"
)

// endpoint: GET /api/v1/clubs/{clubid}/fines
func handleGetFines(w http.ResponseWriter, r *http.Request) {

	type Fine struct {
		ID        string  `json:"id"`
		UserID    string  `json:"userId"`
		UserName  string  `json:"userName"`
		Reason    string  `json:"reason"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
		Paid      bool    `json:"paid"`
	}

	clubID := extractPathParam(r, "clubs")
	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)
	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	fines, err := club.GetFines()
	if err != nil {
		http.Error(w, "Failed to retrieve fines", http.StatusInternalServerError)
		return
	}

	var fineList []Fine
	for _, fine := range fines {

		user, err := models.GetUserByID(fine.UserID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Default().Printf("Fine: %v, user %s", fine, user.Name)

		fineList = append(fineList, Fine{
			ID:        fine.ID,
			UserID:    fine.UserID,
			UserName:  user.Name,
			Reason:    fine.Reason,
			Amount:    fine.Amount,
			CreatedAt: fine.CreatedAt,
			UpdatedAt: fine.UpdatedAt,
			Paid:      fine.Paid,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fineList)
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

	clubId := extractQueryParam(r, "clubId")

	fines, err := user.GetFines(clubId)
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

// endpoint: POST /api/v1/clubs/{clubid}/fines
func handleCreateFine(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		UserID string  `json:"userId"`
		Reason string  `json:"reason"`
		Amount float64 `json:"amount"`
	}

	clubID := extractPathParam(r, "clubs")
	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)
	if !club.IsAdmin(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.UserID == "" || payload.Reason == "" || payload.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	fine, err := club.CreateFine(payload.UserID, payload.Reason, payload.Amount)
	if err != nil {
		http.Error(w, "Failed to create fine", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fine)
}

// endpoint: PATCH /api/v1/clubs/{clubid}/fines/{fineid}
func handleUpdateFine(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Paid bool `json:"paid"`
	}

	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	fineID := extractPathParam(r, "fines")
	fine, err := club.GetFineByID(fineID)
	if err != nil {
		http.Error(w, "Fine not found", http.StatusNotFound)
		return
	}

	user := extractUser(r)
	if !club.IsAdmin(user) && fine.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = fine.SetPaid(payload.Paid)
	if err != nil {
		http.Error(w, "Failed to update fine", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fine)
}
