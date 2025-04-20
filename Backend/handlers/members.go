package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// endpoint: GET /api/v1/clubs/{clubid}/members
func handleGetClubMembers(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")

	if len(segments) != 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
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
