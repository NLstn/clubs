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
	userID, err := models.FindOrCreateUser(email)
	if err != nil {
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Create session token or JWT
	jwt, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "JWT error", http.StatusInternalServerError)
		return
	}

	// Return token
	json.NewEncoder(w).Encode(map[string]string{
		"token": jwt,
	})

	if err := models.DeleteMagicLink(token); err != nil {
		// Log the error for debugging and monitoring purposes
		// Replace this with your logging framework if applicable
		http.Error(w, "Failed to delete magic link", http.StatusInternalServerError)
		return
	}
}
