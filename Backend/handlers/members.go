package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

func handleClubMembers(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")

	if len(segments) != 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	clubID := segments[3]

	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		members, err := models.GetClubMembers(clubID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)

	case http.MethodPost:
		var member models.Member
		if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !member.Validate() {
			http.Error(w, "Email and name are required", http.StatusBadRequest)
			return
		}

		// Check if userID exists in context
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

		if err := models.AddMember(&member, clubID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		member.NotifyAdded(club.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(member)

	case http.MethodDelete:
		// Extract member ID from the URL path
		if len(segments) != 6 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		memberID := segments[5]

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
}

func handleClubMemberDelete(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")

	if len(segments) != 6 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	clubID := segments[3]
	memberID := segments[5]

	// Check if club exists
	club, err := models.GetClubByID(clubID)
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if userID exists in context
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

	// Delete the member
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
