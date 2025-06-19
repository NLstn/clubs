package handlers

import (
	"net/http"
	"testing"

	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestCookieAuthenticationFlow(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)
	MockEnvironmentVariables(t)

	handler := GetTestHandler()

	t.Run("Full Cookie Authentication Flow", func(t *testing.T) {
		// Step 1: Create a user and magic link
		email := "integration@example.com"
		_, _ = CreateTestUser(t, email)
		token, err := models.CreateMagicLink(email)
		assert.NoError(t, err)

		// Step 2: Verify magic link - should set cookies
		verifyReq := MakeRequest(t, "GET", "/api/v1/auth/verifyMagicLink?token="+token, nil, "")
		verifyRR := ExecuteRequest(t, handler, verifyReq)

		CheckResponseCode(t, http.StatusNoContent, verifyRR.Code)

		// Extract cookies from response
		cookies := verifyRR.Result().Cookies()
		var accessCookie, refreshCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				accessCookie = cookie
			}
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie
			}
		}

		assert.NotNil(t, accessCookie)
		assert.NotNil(t, refreshCookie)

		// Step 3: Make authenticated request using cookies
		authReq := MakeRequest(t, "GET", "/api/v1/me", nil, "")
		authReq.AddCookie(accessCookie)
		authRR := ExecuteRequest(t, handler, authReq)

		CheckResponseCode(t, http.StatusOK, authRR.Code)

		// Step 4: Refresh token using cookies
		refreshReq := MakeRequest(t, "POST", "/api/v1/auth/refreshToken", nil, "")
		refreshReq.AddCookie(refreshCookie)
		refreshRR := ExecuteRequest(t, handler, refreshReq)

		CheckResponseCode(t, http.StatusNoContent, refreshRR.Code)

		// Should get new access token cookie
		refreshCookies := refreshRR.Result().Cookies()
		var newAccessCookie *http.Cookie
		for _, cookie := range refreshCookies {
			if cookie.Name == "access_token" {
				newAccessCookie = cookie
			}
		}
		assert.NotNil(t, newAccessCookie)
		// Note: tokens might be the same if generated at same time, but that's okay
		assert.NotEmpty(t, newAccessCookie.Value)

		// Step 5: Logout using cookies - should clear cookies
		logoutReq := MakeRequest(t, "POST", "/api/v1/auth/logout", nil, "")
		logoutReq.AddCookie(refreshCookie)
		logoutRR := ExecuteRequest(t, handler, logoutReq)

		CheckResponseCode(t, http.StatusNoContent, logoutRR.Code)

		// Should have cookies with MaxAge=-1 to clear them
		logoutCookies := logoutRR.Result().Cookies()
		var clearedAccessCookie, clearedRefreshCookie *http.Cookie
		for _, cookie := range logoutCookies {
			if cookie.Name == "access_token" {
				clearedAccessCookie = cookie
			}
			if cookie.Name == "refresh_token" {
				clearedRefreshCookie = cookie
			}
		}

		assert.NotNil(t, clearedAccessCookie)
		assert.NotNil(t, clearedRefreshCookie)
		assert.Equal(t, -1, clearedAccessCookie.MaxAge)
		assert.Equal(t, -1, clearedRefreshCookie.MaxAge)
	})
}