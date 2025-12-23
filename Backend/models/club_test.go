package models_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NLstn/civo/database"
	"github.com/NLstn/civo/handlers"
	"github.com/NLstn/civo/models"
	"github.com/NLstn/civo/odata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllClubs(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("no clubs", func(t *testing.T) {
		clubs, err := models.GetAllClubs()
		assert.NoError(t, err)
		assert.Len(t, clubs, 0)
	})

	t.Run("with clubs", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "clubowner@example.com")
		club1 := handlers.CreateTestClub(t, user, "Club 1")
		club2 := handlers.CreateTestClub(t, user, "Club 2")

		clubs, err := models.GetAllClubs()
		assert.NoError(t, err)
		assert.Len(t, clubs, 2)

		// Check that both clubs are returned
		clubMap := make(map[string]models.Club)
		for _, club := range clubs {
			clubMap[club.ID] = club
		}
		assert.Contains(t, clubMap, club1.ID)
		assert.Contains(t, clubMap, club2.ID)
	})

	t.Run("excludes deleted clubs", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "deletedowner@example.com")
		club := handlers.CreateTestClub(t, user, "To Delete Club")

		// Soft delete the club
		err := club.SoftDelete(user.ID)
		assert.NoError(t, err)

		clubs, err := models.GetAllClubs()
		assert.NoError(t, err)

		// Verify deleted club is not included
		for _, c := range clubs {
			assert.NotEqual(t, club.ID, c.ID)
		}
	})
}

func TestGetAllClubsIncludingDeleted(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("includes deleted clubs", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "allclubsowner@example.com")
		activeClub := handlers.CreateTestClub(t, user, "Active Club")
		deletedClub := handlers.CreateTestClub(t, user, "Deleted Club")

		// Soft delete one club
		err := deletedClub.SoftDelete(user.ID)
		assert.NoError(t, err)

		clubs, err := models.GetAllClubsIncludingDeleted()
		assert.NoError(t, err)
		assert.Len(t, clubs, 2)

		// Check that both clubs are returned
		clubMap := make(map[string]models.Club)
		for _, club := range clubs {
			clubMap[club.ID] = club
		}
		assert.Contains(t, clubMap, activeClub.ID)
		assert.Contains(t, clubMap, deletedClub.ID)
	})
}

func TestGetClubByID(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("existing club", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "getclubuser@example.com")
		createdClub := handlers.CreateTestClub(t, user, "Test Club")

		club, err := models.GetClubByID(createdClub.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdClub.ID, club.ID)
		assert.Equal(t, createdClub.Name, club.Name)
	})

	t.Run("non-existent club", func(t *testing.T) {
		club, err := models.GetClubByID("non-existent-id")
		assert.Error(t, err)
		assert.Equal(t, "", club.ID)
	})

	t.Run("empty ID", func(t *testing.T) {
		club, err := models.GetClubByID("")
		assert.Error(t, err)
		assert.Equal(t, "", club.ID)
	})
}

func TestGetClubsByFilter(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Set up OData service for testing
	service, err := odata.NewService(database.Db)
	require.NoError(t, err, "Failed to create OData service")

	odataV2Mux := http.NewServeMux()
	service.RegisterCustomHandlers(odataV2Mux)
	odataV2Mux.Handle("/", service)
	handler := http.StripPrefix("/api/v2", handlers.CompositeAuthMiddleware(odataV2Mux))

	t.Run("get multiple clubs via OData", func(t *testing.T) {
		user, token := handlers.CreateTestUser(t, "multiclubuser@example.com")
		club1 := handlers.CreateTestClub(t, user, "Club 1")
		club2 := handlers.CreateTestClub(t, user, "Club 2")
		club3 := handlers.CreateTestClub(t, user, "Club 3")

		// Get all clubs via OData
		req := httptest.NewRequest("GET", "/api/v2/Clubs", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		clubs := result["value"].([]interface{})
		assert.GreaterOrEqual(t, len(clubs), 3, "Should have at least 3 clubs")

		// Check that all clubs are returned
		clubIDs := make(map[string]bool)
		for _, c := range clubs {
			clubMap := c.(map[string]interface{})
			clubIDs[clubMap["ID"].(string)] = true
		}
		assert.True(t, clubIDs[club1.ID], "Club 1 should be in results")
		assert.True(t, clubIDs[club2.ID], "Club 2 should be in results")
		assert.True(t, clubIDs[club3.ID], "Club 3 should be in results")
	})

	t.Run("get club by ID via OData", func(t *testing.T) {
		user, token := handlers.CreateTestUser(t, "getclubuser@example.com")
		club1 := handlers.CreateTestClub(t, user, "Existing Club")

		// Get specific club via OData
		req := httptest.NewRequest("GET", "/api/v2/Clubs("+club1.ID+")", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var club map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&club)
		require.NoError(t, err)

		assert.Equal(t, club1.ID, club["ID"])
		assert.Equal(t, club1.Name, club["Name"])
	})

	t.Run("get non-existent club via OData", func(t *testing.T) {
		_, token := handlers.CreateTestUser(t, "nonexistuser@example.com")

		req := httptest.NewRequest("GET", "/api/v2/Clubs(non-existent-id)", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCreateClub(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Set up OData service for testing
	service, err := odata.NewService(database.Db)
	require.NoError(t, err, "Failed to create OData service")

	odataV2Mux := http.NewServeMux()
	service.RegisterCustomHandlers(odataV2Mux)
	odataV2Mux.Handle("/", service)
	handler := http.StripPrefix("/api/v2", handlers.CompositeAuthMiddleware(odataV2Mux))

	t.Run("create valid club via OData", func(t *testing.T) {
		user, token := handlers.CreateTestUser(t, "createclubuser@example.com")

		// Create club via OData POST
		clubData := map[string]interface{}{
			"Name":        "New Test Club",
			"Description": "Test Description",
		}
		body, err := json.Marshal(clubData)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v2/Clubs", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoError(t, err)

		assert.NotEmpty(t, created["ID"])
		assert.Equal(t, "New Test Club", created["Name"])
		assert.Equal(t, "Test Description", created["Description"])
		assert.Equal(t, user.ID, created["CreatedBy"])
		assert.NotEmpty(t, created["CreatedAt"])

		// Verify club was actually saved to database
		var dbClub models.Club
		err = database.Db.Where("id = ?", created["ID"]).First(&dbClub).Error
		assert.NoError(t, err)
		assert.Equal(t, "New Test Club", dbClub.Name)

		// Verify owner member was created
		var member models.Member
		err = database.Db.Where("user_id = ? AND club_id = ?", user.ID, created["ID"]).First(&member).Error
		assert.NoError(t, err)
		assert.Equal(t, "owner", member.Role)

		// Verify club settings were created
		var settings models.ClubSettings
		err = database.Db.Where("club_id = ?", created["ID"]).First(&settings).Error
		assert.NoError(t, err)
	})

	t.Run("create club with empty name via OData", func(t *testing.T) {
		_, token := handlers.CreateTestUser(t, "emptyclubuser@example.com")

		clubData := map[string]interface{}{
			"Name":        "",
			"Description": "Description",
		}
		body, err := json.Marshal(clubData)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v2/Clubs", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		// OData allows empty name since Name is not marked as required
		// This is different from the old direct function behavior
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

func TestClubUpdate(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("update club name", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "updateclubuser@example.com")
		club := handlers.CreateTestClub(t, user, "Original Name")

		err := club.Update("Updated Name", "Updated Description", user.ID)
		assert.NoError(t, err)

		// Verify update in database
		var dbClub models.Club
		err = database.Db.Where("id = ?", club.ID).First(&dbClub).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", dbClub.Name)
		assert.NotNil(t, dbClub.Description)
		assert.Equal(t, "Updated Description", *dbClub.Description)
	})

	t.Run("update non-existent club", func(t *testing.T) {
		club := models.Club{
			ID:   "non-existent-id",
			Name: "Non-existent Club",
		}
		// The current implementation doesn't validate club existence before update
		// So this will succeed (no rows affected but no error)
		err := club.Update("New Name", "New Description", "user-id")
		assert.NoError(t, err)
	})
}

func TestClubSoftDelete(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("soft delete existing club", func(t *testing.T) {
		user, _ := handlers.CreateTestUser(t, "softdeleteuser@example.com")
		club := handlers.CreateTestClub(t, user, "To Soft Delete")

		err := club.SoftDelete(user.ID)
		assert.NoError(t, err)

		// Verify club is marked as deleted
		var dbClub models.Club
		err = database.Db.Unscoped().Where("id = ?", club.ID).First(&dbClub).Error
		assert.NoError(t, err)
		assert.True(t, dbClub.Deleted)
		assert.NotNil(t, dbClub.DeletedAt)

		// Verify club is not returned by normal queries
		// Note: GetClubByID doesn't filter soft-deleted clubs in current implementation
		retrievedClub, err := models.GetClubByID(club.ID)
		assert.NoError(t, err)                // Will still find the club
		assert.True(t, retrievedClub.Deleted) // But it should be marked as deleted
	})

	t.Run("soft delete non-existent club", func(t *testing.T) {
		club := models.Club{
			ID: "non-existent-id",
		}
		// The current implementation doesn't validate club existence before soft delete
		// So this will succeed (no rows affected but no error)
		err := club.SoftDelete("user-id")
		assert.NoError(t, err)
	})
}

func TestSoftDeleteClub(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Set up OData service for testing
	service, err := odata.NewService(database.Db)
	require.NoError(t, err, "Failed to create OData service")

	odataV2Mux := http.NewServeMux()
	service.RegisterCustomHandlers(odataV2Mux)
	odataV2Mux.Handle("/", service)
	handler := http.StripPrefix("/api/v2", handlers.CompositeAuthMiddleware(odataV2Mux))

	t.Run("soft delete club via OData PATCH", func(t *testing.T) {
		user, token := handlers.CreateTestUser(t, "softdeleteuser@example.com")
		club := handlers.CreateTestClub(t, user, "To Soft Delete")

		// Soft delete via OData PATCH
		deleteData := map[string]interface{}{
			"Deleted": true,
		}
		body, err := json.Marshal(deleteData)
		require.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/v2/Clubs("+club.ID+")", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent)

		// Verify club is marked as deleted
		var dbClub models.Club
		err = database.Db.Unscoped().Where("id = ?", club.ID).First(&dbClub).Error
		assert.NoError(t, err)
		assert.True(t, dbClub.Deleted)
		assert.NotNil(t, dbClub.DeletedAt)
		assert.NotNil(t, dbClub.DeletedBy)
		assert.Equal(t, user.ID, *dbClub.DeletedBy)

		// Verify club is not returned by OData GET
		req = httptest.NewRequest("GET", "/api/v2/Clubs("+club.ID+")", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp = rec.Result()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("soft delete non-existent club via OData", func(t *testing.T) {
		_, token := handlers.CreateTestUser(t, "nonexistdelete@example.com")

		deleteData := map[string]interface{}{
			"Deleted": true,
		}
		body, err := json.Marshal(deleteData)
		require.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/v2/Clubs(non-existent-id)", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
