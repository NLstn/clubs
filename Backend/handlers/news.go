package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerNewsRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/news", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNews(w, r)
		case http.MethodPost:
			handleCreateNews(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/news/{newsid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handleUpdateNews(w, r)
		case http.MethodDelete:
			handleDeleteNews(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// GET /api/v1/clubs/{clubid}/news
func handleGetNews(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
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

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	news, err := club.GetNews()
	if err != nil {
		http.Error(w, "Failed to get news", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

// POST /api/v1/clubs/{clubid}/news
func handleCreateNews(w http.ResponseWriter, r *http.Request) {
	type CreateNewsRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
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

	var req CreateNewsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	news, err := club.CreateNews(req.Title, req.Content, user.ID)
	if err != nil {
		http.Error(w, "Failed to create news", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(news)
}

// PUT /api/v1/clubs/{clubid}/news/{newsid}
func handleUpdateNews(w http.ResponseWriter, r *http.Request) {
	type UpdateNewsRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	newsID := extractPathParam(r, "news")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(newsID); err != nil {
		http.Error(w, "Invalid news ID format", http.StatusBadRequest)
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

	var req UpdateNewsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	news, err := club.UpdateNews(newsID, req.Title, req.Content, user.ID)
	if err != nil {
		http.Error(w, "Failed to update news", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

// DELETE /api/v1/clubs/{clubid}/news/{newsid}
func handleDeleteNews(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")
	newsID := extractPathParam(r, "news")

	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(newsID); err != nil {
		http.Error(w, "Invalid news ID format", http.StatusBadRequest)
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

	err = club.DeleteNews(newsID)
	if err != nil {
		http.Error(w, "Failed to delete news", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}