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
				expectedStatus: http.StatusNoContent,
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

				// If successful, check cookies are set (since we use cookie-only auth)
				if tt.expectedStatus == http.StatusNoContent {
					cookies := rr.Result().Cookies()
					var accessCookie, refreshCookie *http.Cookie
					
					for _, cookie := range cookies {
						if cookie.Name == "access_token" {
							accessCookie = cookie
						}
						if cookie.Name == "refresh_token" {
							refreshCookie = cookie
						}
					}
					
					assert.NotNil(t, accessCookie, "Access token cookie should be set")
					assert.NotNil(t, refreshCookie, "Refresh token cookie should be set")
					assert.NotEmpty(t, accessCookie.Value, "Access token should have value")
					assert.NotEmpty(t, refreshCookie.Value, "Refresh token should have value")
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
		err = user.StoreRefreshToken(refreshToken)
		assert.NoError(t, err)

		tests := []struct {
			name           string
			cookie         *http.Cookie
			expectedStatus int
			shouldContain  string
		}{
			{
				name:           "Valid refresh token",
				cookie:         &http.Cookie{Name: "refresh_token", Value: refreshToken},
				expectedStatus: http.StatusNoContent,
			},
			{
				name:           "Missing refresh token",
				cookie:         nil,
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Refresh token required",
			},
			{
				name:           "Invalid refresh token",
				cookie:         &http.Cookie{Name: "refresh_token", Value: "invalid-token"},
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Invalid refresh token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := MakeRequest(t, "POST", "/api/v1/auth/refreshToken", nil, "")
				if tt.cookie != nil {
					req.AddCookie(tt.cookie)
				}

				rr := ExecuteRequest(t, handler, req)

				CheckResponseCode(t, tt.expectedStatus, rr.Code)
				if tt.shouldContain != "" {
					AssertContains(t, rr.Body.String(), tt.shouldContain)
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
		err = user.StoreRefreshToken(refreshToken)
		assert.NoError(t, err)

		tests := []struct {
			name           string
			cookie         *http.Cookie
			expectedStatus int
			shouldContain  string
		}{
			{
				name:           "Valid refresh token",
				cookie:         &http.Cookie{Name: "refresh_token", Value: refreshToken},
				expectedStatus: http.StatusNoContent,
			},
			{
				name:           "Missing refresh token",
				cookie:         nil,
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Refresh token required",
			},
			{
				name:           "Invalid refresh token",
				cookie:         &http.Cookie{Name: "refresh_token", Value: "invalid-token"},
				expectedStatus: http.StatusUnauthorized,
				shouldContain:  "Invalid refresh token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := MakeRequest(t, "POST", "/api/v1/auth/logout", nil, "")
				if tt.cookie != nil {
					req.AddCookie(tt.cookie)
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

	t.Run("Cookie Authentication", func(t *testing.T) {
		t.Run("Verify Magic Link Sets Cookies", func(t *testing.T) {
			// Create a user and magic link
			email := "cookie-test@example.com"  
			_, _ = CreateTestUser(t, email)
			
			token, err := models.CreateMagicLink(email)
			assert.NoError(t, err)

			url := "/api/v1/auth/verifyMagicLink?token=" + token
			req := MakeRequest(t, "GET", url, nil, "")
			rr := ExecuteRequest(t, handler, req)

			CheckResponseCode(t, http.StatusNoContent, rr.Code)

			// Check that cookies are set in the response
			cookies := rr.Result().Cookies()
			var accessCookie, refreshCookie *http.Cookie

			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					accessCookie = cookie
				}
				if cookie.Name == "refresh_token" {
					refreshCookie = cookie
				}
			}

			// Verify cookies are set with secure attributes
			assert.NotNil(t, accessCookie, "Access token cookie should be set")
			assert.NotNil(t, refreshCookie, "Refresh token cookie should be set")
			
			if accessCookie != nil {
				assert.True(t, accessCookie.HttpOnly, "Access token cookie should be HttpOnly")
				assert.Equal(t, http.SameSiteStrictMode, accessCookie.SameSite, "Access token cookie should have SameSite=Strict")
				assert.NotEmpty(t, accessCookie.Value, "Access token cookie should have a value")
			}
			
			if refreshCookie != nil {
				assert.True(t, refreshCookie.HttpOnly, "Refresh token cookie should be HttpOnly")
				assert.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite, "Refresh token cookie should have SameSite=Strict")
				assert.NotEmpty(t, refreshCookie.Value, "Refresh token cookie should have a value")
			}
		})

		t.Run("Refresh Token From Cookie", func(t *testing.T) {
			// Create a user and refresh token
			user, _ := CreateTestUser(t, "cookie-refresh@example.com")
			
			refreshToken, err := auth.GenerateRefreshToken(user.ID)
			assert.NoError(t, err)
			
			err = user.StoreRefreshToken(refreshToken)
			assert.NoError(t, err)

			// Make request with refresh token in cookie
			req := MakeRequest(t, "POST", "/api/v1/auth/refreshToken", nil, "")
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: refreshToken,
			})

			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusNoContent, rr.Code)

			// Check that new access token cookie is set
			cookies := rr.Result().Cookies()
			var accessCookie *http.Cookie

			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					accessCookie = cookie
				}
			}

			assert.NotNil(t, accessCookie, "New access token cookie should be set")
			if accessCookie != nil {
				assert.True(t, accessCookie.HttpOnly, "Access token cookie should be HttpOnly")
				assert.NotEmpty(t, accessCookie.Value, "Access token cookie should have a value")
			}
		})

		t.Run("Logout Clears Cookies", func(t *testing.T) {
			// Create a user and refresh token
			user, _ := CreateTestUser(t, "cookie-logout@example.com")
			
			refreshToken, err := auth.GenerateRefreshToken(user.ID)
			assert.NoError(t, err)
			
			err = user.StoreRefreshToken(refreshToken)
			assert.NoError(t, err)

			// Make logout request with refresh token in cookie
			req := MakeRequest(t, "POST", "/api/v1/auth/logout", nil, "")
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: refreshToken,
			})

			rr := ExecuteRequest(t, handler, req)
			CheckResponseCode(t, http.StatusNoContent, rr.Code)

			// Check that cookies are cleared (should have MaxAge=-1)
			cookies := rr.Result().Cookies()
			var accessCookie, refreshCookie *http.Cookie

			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					accessCookie = cookie
				}
				if cookie.Name == "refresh_token" {
					refreshCookie = cookie
				}
			}

			// Cookies should be present with MaxAge=-1 to clear them
			assert.NotNil(t, accessCookie, "Access token cookie should be present for clearing")
			assert.NotNil(t, refreshCookie, "Refresh token cookie should be present for clearing")
			
			if accessCookie != nil {
				assert.Equal(t, -1, accessCookie.MaxAge, "Access token cookie should have MaxAge=-1 for clearing")
			}
			if refreshCookie != nil {
				assert.Equal(t, -1, refreshCookie.MaxAge, "Refresh token cookie should have MaxAge=-1 for clearing")
			}
		})
	})
}