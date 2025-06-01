package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
)

func registerFineTemplateRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/fine-templates", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetFineTemplates(w, r)
		case http.MethodPost:
			handleCreateFineTemplate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/fine-templates/{templateid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handleUpdateFineTemplate(w, r)
		case http.MethodDelete:
			handleDeleteFineTemplate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/clubs/{clubid}/fine-templates
func handleGetFineTemplates(w http.ResponseWriter, r *http.Request) {
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

	templates, err := club.GetFineTemplates()
	if err != nil {
		http.Error(w, "Failed to retrieve fine templates", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(templates)
}

// endpoint: POST /api/v1/clubs/{clubid}/fine-templates
func handleCreateFineTemplate(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
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

	if payload.Description == "" || payload.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	template, err := club.CreateFineTemplate(payload.Description, payload.Amount, user.ID)
	if err != nil {
		http.Error(w, "Failed to create fine template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

// endpoint: PUT /api/v1/clubs/{clubid}/fine-templates/{templateid}
func handleUpdateFineTemplate(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
	}

	clubID := extractPathParam(r, "clubs")
	templateID := extractPathParam(r, "fine-templates")
	
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

	if payload.Description == "" || payload.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	template, err := club.UpdateFineTemplate(templateID, payload.Description, payload.Amount, user.ID)
	if err != nil {
		http.Error(w, "Failed to update fine template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(template)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/fine-templates/{templateid}
func handleDeleteFineTemplate(w http.ResponseWriter, r *http.Request) {
	clubID := extractPathParam(r, "clubs")
	templateID := extractPathParam(r, "fine-templates")
	
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

	err = club.DeleteFineTemplate(templateID, user.ID)
	if err != nil {
		http.Error(w, "Failed to delete fine template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}