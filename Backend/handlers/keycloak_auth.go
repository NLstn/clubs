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
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/csrf"
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

	// Get client IP for state validation
	clientIP := getClientIP(r)
	ipHash := csrf.HashIP(clientIP)

	// Generate a cryptographically signed state token with HMAC
	// Format: nonce.timestamp.signature
	// This provides CSRF protection without requiring server-side storage
	state, err := csrf.GenerateStateToken(ipHash)
	if err != nil {
		log.Printf("Failed to generate state token: %v", err)
		http.Error(w, "Failed to generate authentication URL", http.StatusInternalServerError)
		return
	}

	// Store the state nonce in database for one-time use validation
	// Extract nonce from state token (first part before first dot)
	nonceParts := strings.Split(state, ".")
	if len(nonceParts) != 3 {
		http.Error(w, "Failed to generate authentication URL", http.StatusInternalServerError)
		return
	}
	nonce := nonceParts[0]

	// Get the auth URL with PKCE parameters
	authURL, codeVerifier, err := keycloakAuth.GetAuthURLWithPKCE(state)
	if err != nil {
		http.Error(w, "Failed to generate authentication URL", http.StatusInternalServerError)
		return
	}

	// Store state with nonce and code verifier for replay protection and PKCE validation
	if err := models.CreateOAuthState(nonce, codeVerifier); err != nil {
		log.Printf("Failed to store OAuth state: %v", err)
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

	// CSRF Protection: Validate state token with HMAC signature and timestamp
	clientIP := getClientIP(r)
	ipHash := csrf.HashIP(clientIP)

	// Validate the signed state token
	nonce, valid := csrf.ValidateStateToken(req.State, ipHash)
	if !valid {
		log.Printf("Invalid state token: signature verification failed or token expired")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Validate and consume the state nonce (one-time use, prevents replay attacks)
	oauthState, err := models.GetOAuthStateByState(nonce)
	if err != nil {
		log.Printf("OAuth state validation error: %v", err)
		http.Error(w, "Invalid or already used state parameter", http.StatusBadRequest)
		return
	}

	// Verify the code verifier matches (prevents PKCE downgrade attacks)
	if oauthState.CodeVerifier != req.CodeVerifier {
		log.Printf("OAuth state validation failed: code verifier mismatch")
		http.Error(w, "Invalid code verifier", http.StatusBadRequest)
		return
	}

	// Delete the state (one-time use)
	if err := models.DeleteOAuthState(nonce); err != nil {
		log.Printf("Warning: Failed to delete OAuth state: %v", err)
		// Continue anyway - state will be cleaned up by periodic job
	}

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

// getClientIP extracts the client IP address from the request
//
// Security Note: This function trusts X-Forwarded-For and X-Real-IP headers.
// In production environments behind a reverse proxy, ensure the proxy is configured
// to set these headers correctly and that direct client access to the backend is blocked.
// If X-Forwarded-For can be spoofed by clients, IP-based validation can be bypassed.
//
// For enhanced security in untrusted environments:
// - Configure your reverse proxy to strip/override client-provided headers
// - Use RemoteAddr only (comment out X-Forwarded-For logic)
// - Consider making IP validation optional via configuration
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (common in proxied environments)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// Fall back to RemoteAddr
	// RemoteAddr is in format "IP:port", extract just the IP
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// generateRandomState generates a random state parameter
// Deprecated: Use csrf.GenerateStateToken for CSRF-protected state tokens
func generateRandomState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple timestamp-based state
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)
}
