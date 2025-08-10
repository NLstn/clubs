package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

func registerClubRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAllClubs(w, r)
		case http.MethodPost:
			handleCreateClub(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubByID(w, r)
		case http.MethodPatch:
			handleUpdateClub(w, r)
		case http.MethodDelete:
			handleDeleteClub(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/hard-delete", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleHardDeleteClub(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/logo", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleUploadClubLogo(w, r)
		case http.MethodDelete:
			handleDeleteClubLogo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// ClubWithRole represents a club with the user's role in that club
type ClubWithRole struct {
	models.Club
	UserRole  string        `json:"user_role"`
	UserTeams []models.Team `json:"user_teams,omitempty"`
}

// endpoint: GET /api/v1/clubs
func handleGetAllClubs(w http.ResponseWriter, r *http.Request) {

	user := extractUser(r)

	// Get all clubs including deleted ones for owners to see
	clubs, err := models.GetAllClubsIncludingDeleted()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var authorizedClubs []ClubWithRole
	for _, club := range clubs {
		if club.IsMember(user) {
			// If club is deleted, only show to owners
			if club.Deleted && !club.IsOwner(user) {
				continue
			}

			// Get user's role in this club
			role, err := club.GetMemberRole(user)
			if err != nil {
				// If we can't get the role but they are a member, default to "member"
				role = "member"
			}

			// Get user's teams in this club
			userTeams, err := models.GetUserTeams(user.ID, club.ID)
			if err != nil {
				// If we can't get the teams, default to empty slice
				userTeams = []models.Team{}
			}

			clubWithRole := ClubWithRole{
				Club:      club,
				UserRole:  role,
				UserTeams: userTeams,
			}
			authorizedClubs = append(authorizedClubs, clubWithRole)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authorizedClubs)
}

// endpoint: GET /api/v1/clubs/{clubid}
func handleGetClubByID(w http.ResponseWriter, r *http.Request) {

	user := extractUser(r)

	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsMember(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// If club is deleted, only allow owners to access it
	if club.Deleted && !club.IsOwner(user) {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(club)
}

// endpoint: POST /api/v1/clubs
func handleCreateClub(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	user := extractUser(r)

	var payload Body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		http.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	club, err := models.CreateClub(payload.Name, payload.Description, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(club)
}

// endpoint: PATCH /api/v1/clubs/{clubid}
func handleUpdateClub(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var payload Body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := club.Update(payload.Name, payload.Description, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(club)
}

// endpoint: DELETE /api/v1/clubs/{clubid}
func handleDeleteClub(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized - only owners can delete clubs", http.StatusForbidden)
		return
	}

	if err := club.SoftDelete(user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/hard-delete
func handleHardDeleteClub(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized - only owners can permanently delete clubs", http.StatusForbidden)
		return
	}

	// Only allow hard delete if club is already soft deleted
	if !club.Deleted {
		http.Error(w, "Club must be soft deleted before permanent deletion", http.StatusBadRequest)
		return
	}

	if err := models.DeleteClubPermanently(clubID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: POST /api/v1/clubs/{clubid}/logo
func handleUploadClubLogo(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	// Get club and verify it exists
	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		log.Printf("ERROR: Club not found: %s", clubID)
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Database error getting club: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is an admin or owner of the club
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		log.Printf("ERROR: Unauthorized logo upload attempt by user %s for club %s", user.ID, clubID)
		http.Error(w, "Unauthorized - only club admins and owners can upload logos", http.StatusForbidden)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("ERROR: Unable to parse form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

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

	logoURL, err := azure.UploadClubLogo(clubID, file, header)
	if err != nil {
		log.Printf("ERROR: Failed to upload logo to Azure: %v", err)
		http.Error(w, "Failed to upload logo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update club with new logo URL
	err = club.UpdateLogo(&logoURL, user.ID)
	if err != nil {
		// Try to delete the uploaded file if database update fails
		log.Printf("ERROR: Failed to update club in database: %v", err)
		azure.DeleteClubLogo(logoURL)
		http.Error(w, "Failed to update club: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated club
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"logo_url": logoURL,
		"message":  "Logo uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/logo
func handleDeleteClubLogo(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	// Get club and verify it exists
	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is an admin or owner of the club
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		http.Error(w, "Unauthorized - only club admins and owners can delete logos", http.StatusForbidden)
		return
	}

	// Check if club has a logo
	if club.LogoURL == nil || *club.LogoURL == "" {
		http.Error(w, "Club has no logo to delete", http.StatusBadRequest)
		return
	}

	// Delete logo from storage
	err = azure.DeleteClubLogo(*club.LogoURL)
	if err != nil {
		http.Error(w, "Failed to delete logo from storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update club to remove logo URL
	err = club.UpdateLogo(nil, user.ID)
	if err != nil {
		http.Error(w, "Failed to update club: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Logo deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
