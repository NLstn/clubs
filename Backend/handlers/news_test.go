package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestNewsRoutes(t *testing.T) {
	// Setup test database
	SetupTestDB(t)
	defer TeardownTestDB(t)

	// Create test user and club
	user, accessToken := CreateTestUser(t, "test@example.com")
	club := CreateTestClub(t, user, "Test Club")

	t.Run("Create News", func(t *testing.T) {
		payload := map[string]string{
			"title":   "Test News",
			"content": "This is test news content",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/clubs/"+club.ID+"/news", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		handler := Handler_v1()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.News
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Test News", response.Title)
		assert.Equal(t, "This is test news content", response.Content)
		assert.Equal(t, club.ID, response.ClubID)
	})

	t.Run("Get News", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/clubs/"+club.ID+"/news", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		handler := Handler_v1()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.News
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
		assert.Equal(t, "Test News", response[0].Title)
	})

	t.Run("Update News", func(t *testing.T) {
		// First create a news item
		news, err := club.CreateNews("Original Title", "Original Content", user.ID)
		assert.NoError(t, err)

		payload := map[string]string{
			"title":   "Updated Title",
			"content": "Updated Content",
		}

		body, _ := json.Marshal(payload)

		// Debug: Print the URL being constructed
		url := "/api/v1/clubs/" + club.ID + "/news/" + news.ID
		t.Logf("Update URL: %s", url)

		req := httptest.NewRequest("PUT", url, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		handler := Handler_v1()
		handler.ServeHTTP(w, req)

		t.Logf("Update response status: %d", w.Code)
		t.Logf("Update response body: %s", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.News
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", response.Title)
		assert.Equal(t, "Updated Content", response.Content)
	})

	t.Run("Delete News", func(t *testing.T) {
		// First create a news item
		news, err := club.CreateNews("To Delete", "Content to delete", user.ID)
		assert.NoError(t, err)

		url := "/api/v1/clubs/" + club.ID + "/news/" + news.ID
		t.Logf("Delete URL: %s", url)

		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))

		w := httptest.NewRecorder()
		handler := Handler_v1()
		handler.ServeHTTP(w, req)

		t.Logf("Delete response status: %d", w.Code)
		t.Logf("Delete response body: %s", w.Body.String())

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
