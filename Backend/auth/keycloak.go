package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type KeycloakConfig struct {
	ServerURL   string
	Realm       string
	ClientID    string
	RedirectURL string
}

type KeycloakAuth struct {
	config   KeycloakConfig
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   oauth2.Config
}

type KeycloakUser struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Name              string `json:"name"`
}

var keycloakAuth *KeycloakAuth

func InitKeycloak() error {
	frontendURL := getEnvOrDefault("FRONTEND_URL", "http://localhost:5173")
	redirectURL := frontendURL + "/auth/callback"

	config := KeycloakConfig{
		ServerURL:   getEnvOrDefault("KEYCLOAK_SERVER_URL", "https://auth.clubsstaging.dev"),
		Realm:       getEnvOrDefault("KEYCLOAK_REALM", "clubs-dev"),
		ClientID:    getEnvOrDefault("KEYCLOAK_CLIENT_ID", "clubs-frontend"),
		RedirectURL: redirectURL,
	}

	issuerURL := fmt.Sprintf("%s/realms/%s", config.ServerURL, config.Realm)

	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		return fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: config.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	oauth2Config := oauth2.Config{
		ClientID:    config.ClientID,
		RedirectURL: config.RedirectURL,
		Endpoint:    provider.Endpoint(),
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
	}

	keycloakAuth = &KeycloakAuth{
		config:   config,
		provider: provider,
		verifier: verifier,
		oauth2:   oauth2Config,
	}

	return nil
}

func GetKeycloakAuth() *KeycloakAuth {
	return keycloakAuth
}

func (k *KeycloakAuth) GetAuthURL(state string) string {
	return k.oauth2.AuthCodeURL(state)
}

func (k *KeycloakAuth) ExchangeCodeForTokens(ctx context.Context, code string) (*oauth2.Token, *KeycloakUser, string, error) {
	token, err := k.oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, nil, "", fmt.Errorf("no id_token in token response")
	}

	idToken, err := k.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to verify ID token: %w", err)
	}

	var user KeycloakUser
	if err := idToken.Claims(&user); err != nil {
		return nil, nil, "", fmt.Errorf("failed to parse claims: %w", err)
	}

	return token, &user, rawIDToken, nil
}

func (k *KeycloakAuth) VerifyAccessToken(ctx context.Context, accessToken string) (*KeycloakUser, error) {
	idToken, err := k.verifier.Verify(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify access token: %w", err)
	}

	var user KeycloakUser
	if err := idToken.Claims(&user); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &user, nil
}

func (k *KeycloakAuth) ValidateIDToken(ctx context.Context, idToken string) (*KeycloakUser, error) {
	token, err := k.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var user KeycloakUser
	if err := token.Claims(&user); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &user, nil
}

// GetLogoutURL returns the Keycloak end session URL for logging out
func (k *KeycloakAuth) GetLogoutURL(postLogoutRedirectURI, idTokenHint string) string {
	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", k.config.ServerURL, k.config.Realm)

	params := make([]string, 0, 2)
	if postLogoutRedirectURI != "" {
		params = append(params, fmt.Sprintf("post_logout_redirect_uri=%s", postLogoutRedirectURI))
	}
	if idTokenHint != "" {
		params = append(params, fmt.Sprintf("id_token_hint=%s", idTokenHint))
	}

	if len(params) > 0 {
		logoutURL += "?" + strings.Join(params, "&")
	}

	return logoutURL
}

// KeycloakMiddleware validates Keycloak tokens and sets user context
func KeycloakMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		if keycloakAuth == nil {
			http.Error(w, "Authentication service unavailable", http.StatusInternalServerError)
			return
		}

		user, err := keycloakAuth.VerifyAccessToken(r.Context(), accessToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set user ID from Keycloak subject
		ctx := context.WithValue(r.Context(), UserIDKey, user.Sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
