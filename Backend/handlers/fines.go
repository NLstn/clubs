package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerFineRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/me/fines", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetMyFines(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/fines", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateFine(w, r)
		case http.MethodGet:
			handleGetFines(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/fines/{fineid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleDeleteFine(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/clubs/{clubid}/fines
func handleGetFines(w http.ResponseWriter, r *http.Request) {

	type Fine struct {
		ID        string  `json:"id"`
		UserID    string  `json:"userId"`
		UserName  string  `json:"userName"`
		Reason    string  `json:"reason"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt string  `json:"updatedAt"`
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

		fineList = append(fineList, Fine{
			ID:        fine.ID,
			UserID:    fine.UserID,
			UserName:  user.GetFullName(),
			Reason:    fine.Reason,
			Amount:    fine.Amount,
			CreatedAt: fine.CreatedAt.Format(time.RFC3339),
			UpdatedAt: fine.UpdatedAt.Format(time.RFC3339),
			Paid:      fine.Paid,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fineList)
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

	fine, err := club.CreateFine(payload.UserID, payload.Reason, user.ID, payload.Amount)
	if err != nil {
		http.Error(w, "Failed to create fine", http.StatusInternalServerError)
		return
	}

	type FineResponse struct {
		ID        string  `json:"id"`
		UserID    string  `json:"userId"`
		Reason    string  `json:"reason"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt string  `json:"updatedAt"`
		Paid      bool    `json:"paid"`
	}

	resp := FineResponse{
		ID:        fine.ID,
		UserID:    fine.UserID,
		Reason:    fine.Reason,
		Amount:    fine.Amount,
		CreatedAt: fine.CreatedAt.Format(time.RFC3339),
		UpdatedAt: fine.UpdatedAt.Format(time.RFC3339),
		Paid:      fine.Paid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/fines/{fineid}
func handleDeleteFine(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	fineID := extractPathParam(r, "fines")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(fineID); err != nil {
		http.Error(w, "Invalid fine ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get club information", http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized - admin access required", http.StatusForbidden)
		return
	}

	err = club.DeleteFine(fineID)
	if err != nil {
		http.Error(w, "Failed to delete fine", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
