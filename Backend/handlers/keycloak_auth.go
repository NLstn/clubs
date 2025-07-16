package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
)

// Store for tracking used authorization codes to prevent reuse
var (
	usedCodes = make(map[string]time.Time)
	usedCodesMutex sync.RWMutex
)

// Clean up expired codes every 10 minutes
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cleanupUsedCodes()
		}
	}()
}

func cleanupUsedCodes() {
	usedCodesMutex.Lock()
	defer usedCodesMutex.Unlock()
	
	cutoff := time.Now().Add(-1 * time.Hour) // Remove codes older than 1 hour
	for code, timestamp := range usedCodes {
		if timestamp.Before(cutoff) {
			delete(usedCodes, code)
		}
	}
}

func isCodeUsed(code string) bool {
	usedCodesMutex.RLock()
	defer usedCodesMutex.RUnlock()
	_, exists := usedCodes[code]
	return exists
}

func markCodeAsUsed(code string) {
	usedCodesMutex.Lock()
	defer usedCodesMutex.Unlock()
	usedCodes[code] = time.Now()
}

func registerKeycloakAuthRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/auth/keycloak/login", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleKeycloakLogin(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/keycloak/callback", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleKeycloakCallback(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/keycloak/validate", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleKeycloakTokenValidation(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/keycloak/refresh", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleKeycloakRefresh(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/keycloak/user", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleKeycloakUserCreation(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/auth/keycloak/logout", RateLimitMiddleware(authLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleKeycloakLogout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/auth/keycloak/login
func handleKeycloakLogin(w http.ResponseWriter, r *http.Request) {
	keycloakAuth := auth.GetKeycloakAuth()
	if keycloakAuth == nil {
		http.Error(w, "Keycloak not initialized", http.StatusInternalServerError)
		return
	}

	// Generate a random state for CSRF protection
	state := generateRandomState()

	// Store state in a secure way (you might want to use a session store)
	// For now, we'll return it to the frontend to include in the callback
	authURL := keycloakAuth.GetAuthURL(state)

	json.NewEncoder(w).Encode(map[string]string{
		"authURL": authURL,
		"state":   state,
	})
}

// endpoint: POST /api/v1/auth/keycloak/callback
func handleKeycloakCallback(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Authorization code required", http.StatusBadRequest)
		return
	}

	// Check if this authorization code has already been used
	if isCodeUsed(req.Code) {
		log.Printf("Authorization code has already been used: %s", req.Code[:8]+"...")
		http.Error(w, "Authorization code has already been used", http.StatusBadRequest)
		return
	}

	// Mark the code as used immediately to prevent concurrent requests
	markCodeAsUsed(req.Code)

	keycloakAuth := auth.GetKeycloakAuth()
	if keycloakAuth == nil {
		http.Error(w, "Keycloak not initialized", http.StatusInternalServerError)
		return
	}

	// Exchange code for tokens
	token, keycloakUser, idToken, err := keycloakAuth.ExchangeCodeForTokens(context.Background(), req.Code)
	if err != nil {
		log.Printf("Failed to exchange code for tokens: %v", err)
		http.Error(w, "Failed to authenticate with Keycloak", http.StatusUnauthorized)
		return
	}

	// Find or create user in our database using Keycloak subject as user ID
	user, err := models.FindOrCreateUserWithKeycloakID(keycloakUser.Sub, keycloakUser.Email, keycloakUser.Name)
	if err != nil {
		log.Printf("Failed to find or create user for Keycloak ID %s: %v", keycloakUser.Sub, err)
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Generate our own JWT tokens for the application
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
	userAgent, ipAddress := models.GetDeviceInfo(r)
	if err := user.StoreRefreshToken(refreshToken, userAgent, ipAddress); err != nil {
		log.Printf("Failed to store refresh token for user %s: %v", user.ID, err)
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	// Return our application tokens and user info
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access":          accessToken,
		"refresh":         refreshToken,
		"profileComplete": user.IsProfileComplete(),
		"keycloakTokens": map[string]interface{}{
			"accessToken":  token.AccessToken,
			"refreshToken": token.RefreshToken,
			"idToken":      idToken,
			"expiresIn":    token.Expiry.Unix(),
		},
	})
}

// endpoint: POST /api/v1/auth/keycloak/validate
func handleKeycloakTokenValidation(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Printf("Missing or invalid Authorization header: %s", authHeader)
		http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
		return
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Received ID token for validation (length: %d)", len(idToken))

	keycloakAuth := auth.GetKeycloakAuth()
	if keycloakAuth == nil {
		log.Printf("Keycloak not initialized")
		http.Error(w, "Keycloak not initialized", http.StatusInternalServerError)
		return
	}

	// Validate the ID token and extract user information
	keycloakUser, err := keycloakAuth.ValidateIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("Failed to validate ID token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("Successfully validated token for user: %s (%s)", keycloakUser.Email, keycloakUser.Sub)

	// Find or create user in our database using Keycloak subject as user ID
	user, err := models.FindOrCreateUserWithKeycloakID(keycloakUser.Sub, keycloakUser.Email, keycloakUser.Name)
	if err != nil {
		log.Printf("Failed to find or create user for Keycloak ID %s: %v", keycloakUser.Sub, err)
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Generate our own JWT tokens for the application
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
	userAgent, ipAddress := models.GetDeviceInfo(r)
	if err := user.StoreRefreshToken(refreshToken, userAgent, ipAddress); err != nil {
		log.Printf("Failed to store refresh token for user %s: %v", user.ID, err)
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	// Return our application tokens and user info
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access":          accessToken,
		"refresh":         refreshToken,
		"profileComplete": user.IsProfileComplete(),
	})
}

// endpoint: POST /api/v1/auth/keycloak/refresh
func handleKeycloakRefresh(w http.ResponseWriter, r *http.Request) {
	// This endpoint is similar to the regular refresh token endpoint
	// but could be enhanced to also refresh Keycloak tokens if needed
	handleRefreshToken(w, r)
}

// endpoint: POST /api/v1/auth/keycloak/user
func handleKeycloakUserCreation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Sub        string `json:"sub"`
		Email      string `json:"email"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Sub == "" || req.Email == "" {
		log.Printf("Missing required fields: sub=%s, email=%s", req.Sub, req.Email)
		http.Error(w, "Missing required fields (sub, email)", http.StatusBadRequest)
		return
	}

	log.Printf("Creating/finding user for Keycloak ID: %s, Email: %s", req.Sub, req.Email)

	// Find or create user in our database using Keycloak subject as user ID
	user, err := models.FindOrCreateUserWithKeycloakID(req.Sub, req.Email, req.Name)
	if err != nil {
		log.Printf("Failed to find or create user for Keycloak ID %s: %v", req.Sub, err)
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Generate our own JWT tokens for the application
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
	userAgent, ipAddress := models.GetDeviceInfo(r)
	if err := user.StoreRefreshToken(refreshToken, userAgent, ipAddress); err != nil {
		log.Printf("Failed to store refresh token for user %s: %v", user.ID, err)
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created tokens for user %s (%s)", user.Email, user.ID)

	// Return our application tokens and user info
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access":          accessToken,
		"refresh":         refreshToken,
		"profileComplete": user.IsProfileComplete(),
	})
}

// endpoint: POST /api/v1/auth/keycloak/logout
func handleKeycloakLogout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostLogoutRedirectURI string `json:"post_logout_redirect_uri,omitempty"`
		IDToken               string `json:"id_token,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keycloakAuth := auth.GetKeycloakAuth()
	if keycloakAuth == nil {
		log.Printf("Keycloak not initialized")
		http.Error(w, "Keycloak not initialized", http.StatusInternalServerError)
		return
	}

	// Get the Keycloak logout URL with ID token hint
	logoutURL := keycloakAuth.GetLogoutURL(req.PostLogoutRedirectURI, req.IDToken)

	// Also perform the regular logout (clear refresh tokens from our database)
	refreshToken := r.Header.Get("Authorization")
	if refreshToken != "" {
		userID, err := auth.ValidateRefreshToken(refreshToken)
		if err == nil {
			if user, err := models.GetUserByID(userID); err == nil {
				if err := user.DeleteAllRefreshTokens(); err != nil {
					log.Printf("Failed to delete refresh tokens for user ID %s: %v", userID, err)
				}
			}
		}
	}

	// Return the Keycloak logout URL for the frontend to redirect to
	json.NewEncoder(w).Encode(map[string]string{
		"logoutURL": logoutURL,
	})
}

func generateRandomState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple timestamp-based state
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)
}
