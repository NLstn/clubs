package odata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

// RegisterCustomHandlers registers custom HTTP handlers that don't fit standard OData patterns
// These handlers are mounted alongside the OData service
func (s *Service) RegisterCustomHandlers(mux *http.ServeMux) {
	// Club logo upload - requires multipart/form-data which OData doesn't support natively
	mux.HandleFunc("/Clubs(", s.handleClubCustomRoutes)
}

// handleClubCustomRoutes handles custom routes for Club entity
// This handler intercepts requests matching /api/v2/Clubs({id})/... patterns
func (s *Service) handleClubCustomRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract club ID and action from URL
	// Expected format: /api/v2/Clubs('{clubId}')/UploadLogo
	path := r.URL.Path

	// Parse the club ID from the path
	clubID, action := parseClubCustomRoute(path)
	if clubID == "" || action == "" {
		// Not a custom route, let OData handle it
		http.NotFound(w, r)
		return
	}

	// Route to appropriate handler based on action
	switch action {
	case "UploadLogo":
		s.handleUploadClubLogo(w, r, clubID)
	default:
		http.NotFound(w, r)
	}
}

// parseClubCustomRoute extracts club ID and action from custom route path
// Example: /api/v2/Clubs('abc-123')/UploadLogo -> ("abc-123", "UploadLogo")
func parseClubCustomRoute(path string) (clubID, action string) {
	// Remove /api/v2/Clubs( prefix
	path = strings.TrimPrefix(path, "/api/v2/Clubs(")

	// Find the closing parenthesis
	closeIdx := strings.Index(path, ")")
	if closeIdx == -1 {
		return "", ""
	}

	// Extract club ID (remove quotes if present)
	clubID = strings.Trim(path[:closeIdx], "'\"")

	// Extract action (remove leading /)
	remainder := path[closeIdx+1:]
	action = strings.TrimPrefix(remainder, "/")

	return clubID, action
}

// handleUploadClubLogo handles multipart file upload for club logos
// POST /api/v2/Clubs('{clubID}')/UploadLogo
//
// This is a custom endpoint because OData v4 doesn't natively support multipart/form-data.
// For proper OData media entities, see: https://www.odata.org/getting-started/advanced-tutorial/#media
//
// Request: multipart/form-data with "logo" field containing image file
// Response: 200 OK with JSON containing logo_url
//
// Authorization: User must be admin or owner of the club
func (s *Service) handleUploadClubLogo(w http.ResponseWriter, r *http.Request, clubID string) {
	// Only accept POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user from context
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			log.Printf("ERROR: Database error getting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Get club and verify it exists
	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		log.Printf("ERROR: Club not found: %s", clubID)
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Database error getting club: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if user is an admin or owner of the club
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		log.Printf("ERROR: Unauthorized logo upload attempt by user %s for club %s", userID, clubID)
		http.Error(w, "Forbidden - only club admins and owners can upload logos", http.StatusForbidden)
		return
	}

	// Parse multipart form (max 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("ERROR: Unable to parse form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the logo file
	file, header, err := r.FormFile("logo")
	if err != nil {
		log.Printf("ERROR: No logo file provided: %v", err)
		http.Error(w, "No logo file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Delete existing logo if present
	if club.LogoURL != nil && *club.LogoURL != "" {
		if err := azure.DeleteClubLogo(*club.LogoURL); err != nil {
			// Log error but don't fail the upload
			log.Printf("WARNING: Failed to delete existing logo: %v", err)
		}
	}

	// Upload to Azure Blob Storage
	logoURL, err := azure.UploadClubLogo(clubID, file, header)
	if err != nil {
		log.Printf("ERROR: Failed to upload logo to Azure: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload logo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Update club with new logo URL in database
	err = club.UpdateLogo(&logoURL, userID)
	if err != nil {
		// Try to delete the uploaded file if database update fails
		log.Printf("ERROR: Failed to update club in database: %v", err)
		if deleteErr := azure.DeleteClubLogo(logoURL); deleteErr != nil {
			log.Printf("ERROR: Failed to cleanup uploaded logo after database error: %v", deleteErr)
		}
		http.Error(w, fmt.Sprintf("Failed to update club: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return OData-compatible response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"@odata.context": "/api/v2/$metadata#Clubs/$entity",
		"logo_url":       logoURL,
		"message":        "Logo uploaded successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
	}
}
