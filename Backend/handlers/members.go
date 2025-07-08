package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerMemberRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/clubs/{clubid}/members", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetClubMembers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/isAdmin", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleCheckAdminRights(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/clubs/{clubid}/members/{memberid}", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleClubMemberDelete(w, r)
		case http.MethodPatch:
			handleUpdateMemberRole(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// endpoint: GET /api/v1/clubs/{clubid}/members
func handleGetClubMembers(w http.ResponseWriter, r *http.Request) {

	type APIMember struct {
		ID       string `json:"id"`
		UserId   string `json:"userId"`
		Name     string `json:"name"`
		Role     string `json:"role"`
		JoinedAt string `json:"joinedAt"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	var club models.Club
	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := club.GetClubMembers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Loop at members and load their names
	var apiMembers []APIMember
	for i := range members {
		user, err := models.GetUserByID(members[i].UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var apiMember APIMember
		apiMember.ID = members[i].ID
		apiMember.UserId = members[i].UserID
		apiMember.Name = user.GetFullName()
		apiMember.Role = members[i].Role
		apiMember.JoinedAt = members[i].CreatedAt.Format("2006-01-02T15:04:05Z")
		apiMembers = append(apiMembers, apiMember)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiMembers)
}

// endpoint: DELETE /api/v1/clubs/{clubid}/members/{memberid}
func handleClubMemberDelete(w http.ResponseWriter, r *http.Request) {

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	memberID := extractPathParam(r, "members")
	if _, err := uuid.Parse(memberID); err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !club.IsOwner(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	rowsAffected, err := club.DeleteMember(memberID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: PATCH /api/v1/clubs/{clubid}/members/{memberid}
func handleUpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Role string `json:"role"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	memberID := extractPathParam(r, "members")
	if _, err := uuid.Parse(memberID); err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	user := extractUser(r)
	if user.ID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload Body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Role != "owner" && payload.Role != "admin" && payload.Role != "member" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !club.IsAdmin(user) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = club.UpdateMemberRole(user, memberID, payload.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// endpoint: GET /api/v1/clubs/{clubid}/isAdmin
func handleCheckAdminRights(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	clubID := extractPathParam(r, "clubs")

	club, err := models.GetClubByID(clubID)
	if err != nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	role, err := club.GetMemberRole(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAdmin := role == "owner" || role == "admin"
	isOwner := role == "owner"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"isAdmin": isAdmin,
		"isOwner": isOwner,
	})
}
