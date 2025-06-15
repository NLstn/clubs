package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	frontend "github.com/NLstn/clubs/tools"
)

// setSecureCookie sets a secure HTTP cookie with the given name and value
func setSecureCookie(w http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// clearCookie clears a cookie by setting its MaxAge to -1
func clearCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// getTokenFromCookieOrHeader tries to get token from cookie first, then from Authorization header
func getTokenFromCookieOrHeader(r *http.Request, cookieName string) string {
	// Try cookie first
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	
	// Fall back to Authorization header for backwards compatibility
	return r.Header.Get("Authorization")
}

func registerAuthRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/auth/requestMagicLink", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleRequestMagicLink(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/verifyMagicLink", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			verifyMagicLink(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/refreshToken", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleRefreshToken(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/logout", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleLogout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

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
		log.Printf("Failed to create magic link for email %s: %v", req.Email, err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	link := frontend.MakeMagicLink(token)

	err = auth.SendMagicLinkEmail(req.Email, link)
	if err != nil {
		log.Printf("Failed to send magic link email to %s: %v", req.Email, err)
		http.Error(w, "Failed to send magic link email", http.StatusInternalServerError)
		return
	}

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
		if err != nil {
			log.Printf("Magic link verification error: %v", err)
		}
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Find or create user
	user, err := models.FindOrCreateUser(email)
	if err != nil {
		log.Printf("Failed to find or create user for email %s: %v", email, err)
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Create session token or JWT
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate access token for user %s: %v", user.ID, err)
		http.Error(w, "JWT error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate refresh token for user %s: %v", user.ID, err)
		http.Error(w, "Refresh token error", http.StatusInternalServerError)
		return
	}

	// Store refresh token in the database
	if err := user.StoreRefreshToken(refreshToken); err != nil {
		log.Printf("Failed to store refresh token for user %s: %v", user.ID, err)
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	// Set secure cookies
	setSecureCookie(w, "access_token", accessToken, 15*60)        // 15 minutes
	setSecureCookie(w, "refresh_token", refreshToken, 30*24*3600) // 30 days

	// Return token (for backwards compatibility)
	json.NewEncoder(w).Encode(map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	})

	if err := models.DeleteMagicLink(token); err != nil {
		// Log the error for debugging and monitoring purposes
		log.Printf("Failed to delete magic link after successful verification: %v", err)
		http.Error(w, "Failed to delete magic link", http.StatusInternalServerError)
		return
	}
}

// endpoint: POST /api/v1/auth/refresh
func handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := getTokenFromCookieOrHeader(r, "refresh_token")
	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusUnauthorized)
		return
	}

	log.Default().Println("Refreshing access token")

	userID, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Set new access token cookie
	setSecureCookie(w, "access_token", newAccessToken, 15*60) // 15 minutes

	json.NewEncoder(w).Encode(map[string]string{
		"access": newAccessToken,
	})
}

// endpoint: POST /api/v1/auth/logout
func handleLogout(w http.ResponseWriter, r *http.Request) {
	refreshToken := getTokenFromCookieOrHeader(r, "refresh_token")
	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := user.DeleteAllRefreshTokens(); err != nil {
		log.Printf("Failed to delete refresh tokens for user ID %s: %v", userID, err)
		http.Error(w, "Failed to delete refresh tokens", http.StatusInternalServerError)
		return
	}

	// Clear cookies
	clearCookie(w, "access_token")
	clearCookie(w, "refresh_token")

	w.WriteHeader(http.StatusNoContent)
}
