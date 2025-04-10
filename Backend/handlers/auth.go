package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	frontend "github.com/NLstn/clubs/tools"
)

func requestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	token := auth.GenerateToken()
	expiresAt := time.Now().Add(15 * time.Minute)

	tx := database.Db.Exec(`INSERT INTO magic_links (email, token, expires_at) VALUES ($1, $2, $3)`, req.Email, token, expiresAt)
	if tx.Error != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	link := frontend.MakeMagicLink(token)

	go auth.SendMagicLinkEmail(req.Email, link)

	w.WriteHeader(http.StatusNoContent)
}

func verifyMagicLink(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	var result struct {
		Email     string
		ExpiresAt time.Time
	}

	err := database.Db.Raw(`SELECT email, expires_at FROM magic_links WHERE token = ?`, token).
		Scan(&result).Error

	if err != nil || time.Now().After(result.ExpiresAt) {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Find or create user
	var userID string
	err = database.Db.Raw(`SELECT id FROM users WHERE email = $1`, result.Email).Scan(&userID).Error
	if userID == "" {
		err = database.Db.Raw(`INSERT INTO users (email) VALUES ($1) RETURNING id`, result.Email).Scan(&userID).Error
	}

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

	tx := database.Db.Exec(`DELETE FROM magic_links WHERE token = ?`, token)
	if tx.Error != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
}
