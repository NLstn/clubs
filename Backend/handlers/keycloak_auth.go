package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
)

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

	// Get the auth URL with PKCE parameters
	authURL, codeVerifier, err := keycloakAuth.GetAuthURLWithPKCE(state)
	if err != nil {
		http.Error(w, "Failed to generate authentication URL", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"authURL":      authURL,
		"state":        state,
		"codeVerifier": codeVerifier,
	})
}

// endpoint: POST /api/v1/auth/keycloak/callback
func handleKeycloakCallback(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code         string `json:"code"`
		State        string `json:"state"`
		CodeVerifier string `json:"codeVerifier"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Authorization code required", http.StatusBadRequest)
		return
	}

	if req.CodeVerifier == "" {
		http.Error(w, "Code verifier required", http.StatusBadRequest)
		return
	}

	if req.State == "" {
		http.Error(w, "State parameter required", http.StatusBadRequest)
		return
	}

	// TODO: Validate state parameter against stored state to prevent CSRF attacks
	// For now, we just ensure it's present. In production, implement proper state validation
	// using a secure session store or signed/encrypted state tokens.

	keycloakAuth := auth.GetKeycloakAuth()
	if keycloakAuth == nil {
		http.Error(w, "Keycloak not initialized", http.StatusInternalServerError)
		return
	}

	// Exchange code for tokens with PKCE
	token, keycloakUser, idToken, err := keycloakAuth.ExchangeCodeForTokensWithPKCE(context.Background(), req.Code, req.CodeVerifier)
	if err != nil {
		http.Error(w, "Failed to authenticate with Keycloak", http.StatusUnauthorized)
		return
	}

	// Find or create user in our database using Keycloak subject as user ID
	user, err := models.FindOrCreateUserWithKeycloakID(keycloakUser.Sub, keycloakUser.Email, keycloakUser.Name)
	if err != nil {
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	// Generate our own JWT tokens for the application
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
	userAgent, ipAddress := models.GetDeviceInfo(r)
	if err := user.StoreRefreshToken(refreshToken, userAgent, ipAddress); err != nil {
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

// endpoint: POST /api/v1/auth/keycloak/logout
func handleKeycloakLogout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostLogoutRedirectURI string `json:"post_logout_redirect_uri,omitempty"`
		IDToken               string `json:"id_token,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keycloakAuth := auth.GetKeycloakAuth()
	if keycloakAuth == nil {
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
				user.DeleteAllRefreshTokens()
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
