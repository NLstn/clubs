package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	frontend "github.com/NLstn/clubs/tools"
)

// endpoint: GET /api/v1/auth/requestMagicLink
func handleRequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	token, err := models.CreateMagicLink(req.Email)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	link := frontend.MakeMagicLink(token)

	go auth.SendMagicLinkEmail(req.Email, link)

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/auth/verifyMagicLink
func verifyMagicLink(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	email, valid, err := models.VerifyMagicLink(token)
	if err != nil || !valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Find or create user
	user, err := models.FindOrCreateUser(email)
	if err != nil {
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Create session token or JWT
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "JWT error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Refresh token error", http.StatusInternalServerError)
		return
	}

	// Store refresh token in the database
	if err := user.StoreRefreshToken(refreshToken); err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	// Return token
	json.NewEncoder(w).Encode(map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	})

	if err := models.DeleteMagicLink(token); err != nil {
		// Log the error for debugging and monitoring purposes
		// Replace this with your logging framework if applicable
		http.Error(w, "Failed to delete magic link", http.StatusInternalServerError)
		return
	}
}

// endpoint: GET /api/v1/auth/me
func handleGetMe(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// endpoint: POST /api/v1/auth/me
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
