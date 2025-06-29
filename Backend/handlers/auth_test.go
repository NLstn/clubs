package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpoints(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Request Magic Link", func(t *testing.T) {
		tests := []struct {
			name           string
			requestBody    interface{}
			expectedStatus int
			shouldContain  string
		}{
			{
				name:           "Valid email",
				requestBody:    map[string]string{"email": "test@example.com"},
				expectedStatus: http.StatusNoContent,
			},
			{
				name:           "Missing email",
				requestBody:    map[string]string{},
				expectedStatus: http.StatusBadRequest,
				shouldContain:  "Email required",
			},
			{
				name:           "Empty email",
				requestBody:    map[string]string{"email": ""},
				expectedStatus: http.StatusBadRequest,
				shouldContain:  "Email required",
			},
			{
				name:           "Invalid JSON",
				requestBody:    "invalid-json",
				expectedStatus: http.StatusBadRequest, // Empty email after decode failure
				shouldContain:  "Email required",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := MakeRequest(t, "POST", "/api/v1/auth/requestMagicLink", tt.requestBody, "")
				rr := ExecuteRequest(t, handler, req)

				CheckResponseCode(t, tt.expectedStatus, rr.Code)
				if tt.shouldContain != "" {
					AssertContains(t, rr.Body.String(), tt.shouldContain)
				}
			})
		}
	})

	t.Run("Verify Magic Link", func(t *testing.T) {
		// First create a magic link
		email := "verify@example.com"
		// Create the user first (since FindOrCreateUser has issues with SQLite)
		_, _ = CreateTestUser(t, email)

		token, err := models.CreateMagicLink(email)
		assert.NoError(t, err)

		tests := []struct {
			name           string
			token          string
			expectedStatus int
			shouldContain  string
		}{
			{
				name:           "Valid token",
				token:          token,
				expectedStatus: http.StatusOK,
				shouldContain:  "access",
			},
			{
				name:           "Missing token",
				token:          "",
				expectedStatus: http.StatusBadRequest,
				shouldContain:  "Token required",
			},
			{
				name:           "Invalid token",
				token:          "invalid-token",
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Invalid or expired token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				url := "/api/v1/auth/verifyMagicLink"
				if tt.token != "" {
					url += "?token=" + tt.token
				}

				req := MakeRequest(t, "GET", url, nil, "")
				rr := ExecuteRequest(t, handler, req)

				CheckResponseCode(t, tt.expectedStatus, rr.Code)
				if tt.shouldContain != "" {
					AssertContains(t, rr.Body.String(), tt.shouldContain)
				}

				// If successful, check response structure
				if tt.expectedStatus == http.StatusOK {
					var response map[string]string
					ParseJSONResponse(t, rr, &response)
					assert.Contains(t, response, "access")
					assert.Contains(t, response, "refresh")
					assert.NotEmpty(t, response["access"])
					assert.NotEmpty(t, response["refresh"])
				}
			})
		}
	})

	t.Run("Refresh Token", func(t *testing.T) {
		// Create a user and refresh token
		user, _ := CreateTestUser(t, "refresh@example.com")

		// Generate a refresh token
		refreshToken, err := auth.GenerateRefreshToken(user.ID)
		assert.NoError(t, err)

		// Store it in database
		err = user.StoreRefreshToken(refreshToken, "test-user-agent", "127.0.0.1")
		assert.NoError(t, err)

		tests := []struct {
			name           string
			authHeader     string
			expectedStatus int
			shouldContain  string
		}{
			{
				name:           "Valid refresh token",
				authHeader:     refreshToken,
				expectedStatus: http.StatusOK,
				shouldContain:  "access",
			},
			{
				name:           "Missing refresh token",
				authHeader:     "",
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Refresh token required",
			},
			{
				name:           "Invalid refresh token",
				authHeader:     "invalid-token",
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Invalid refresh token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := MakeRequest(t, "POST", "/api/v1/auth/refreshToken", nil, "")
				if tt.authHeader != "" {
					req.Header.Set("Authorization", tt.authHeader)
				}

				rr := ExecuteRequest(t, handler, req)

				CheckResponseCode(t, tt.expectedStatus, rr.Code)
				if tt.shouldContain != "" {
					AssertContains(t, rr.Body.String(), tt.shouldContain)
				}

				// If successful, check response structure
				if tt.expectedStatus == http.StatusOK {
					var response map[string]string
					ParseJSONResponse(t, rr, &response)
					assert.Contains(t, response, "access")
					assert.Contains(t, response, "refresh")
					assert.NotEmpty(t, response["access"])
					assert.NotEmpty(t, response["refresh"])

					// Verify the new refresh token is different from the old one
					assert.NotEqual(t, refreshToken, response["refresh"])

					// Verify the old refresh token is now invalid
					req2 := MakeRequest(t, "POST", "/api/v1/auth/refreshToken", nil, "")
					req2.Header.Set("Authorization", refreshToken)
					rr2 := ExecuteRequest(t, handler, req2)
					CheckResponseCode(t, http.StatusUnauthorized, rr2.Code)
					AssertContains(t, rr2.Body.String(), "Invalid refresh token")

					// Verify the new refresh token works
					req3 := MakeRequest(t, "POST", "/api/v1/auth/refreshToken", nil, "")
					req3.Header.Set("Authorization", response["refresh"])
					rr3 := ExecuteRequest(t, handler, req3)
					CheckResponseCode(t, http.StatusOK, rr3.Code)
				}
			})
		}
	})

	t.Run("Logout", func(t *testing.T) {
		// Create a user and refresh token
		user, _ := CreateTestUser(t, "logout@example.com")

		// Generate a refresh token
		refreshToken, err := auth.GenerateRefreshToken(user.ID)
		assert.NoError(t, err)

		// Store it in database
		err = user.StoreRefreshToken(refreshToken, "test-user-agent", "127.0.0.1")
		assert.NoError(t, err)

		tests := []struct {
			name           string
			authHeader     string
			expectedStatus int
			shouldContain  string
		}{
			{
				name:           "Valid refresh token",
				authHeader:     refreshToken,
				expectedStatus: http.StatusNoContent,
			},
			{
				name:           "Missing refresh token",
				authHeader:     "",
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Refresh token required",
			},
			{
				name:           "Invalid refresh token",
				authHeader:     "invalid-token",
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Invalid refresh token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := MakeRequest(t, "POST", "/api/v1/auth/logout", nil, "")
				if tt.authHeader != "" {
					req.Header.Set("Authorization", tt.authHeader)
				}

				rr := ExecuteRequest(t, handler, req)

				CheckResponseCode(t, tt.expectedStatus, rr.Code)
				if tt.shouldContain != "" {
					AssertContains(t, rr.Body.String(), tt.shouldContain)
				}
			})
		}
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		endpoints := []string{
			"/api/v1/auth/requestMagicLink",
			"/api/v1/auth/verifyMagicLink",
			"/api/v1/auth/refreshToken",
			"/api/v1/auth/logout",
		}

		for _, endpoint := range endpoints {
			req := MakeRequest(t, "PUT", endpoint, nil, "")
			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
		}
	})
}
