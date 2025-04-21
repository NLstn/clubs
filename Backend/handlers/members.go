package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// endpoint: GET /api/v1/clubs/{clubid}/members
func handleGetClubMembers(w http.ResponseWriter, r *http.Request) {

	type APIMember struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	clubID := extractPathParam(r, "clubs")
	if _, err := uuid.Parse(clubID); err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	members, err := models.GetClubMembers(clubID)
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
		apiMember.Name = user.Name
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
	userIDValue := r.Context().Value(auth.UserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized - authentication required", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)
	if !club.IsOwner(userID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	rowsAffected, err := models.DeleteMember(memberID, clubID)
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
