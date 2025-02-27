package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
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

	var club models.Club
	if result := database.Db.First(&club, "id = ?", clubID); result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var members []models.Member
		if result := database.Db.Where("club_id = ?", clubID).Find(&members); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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

		if member.Email == "" || member.Name == "" {
			http.Error(w, "Email and name are required", http.StatusBadRequest)
			return
		}

		member.ID = uuid.New().String()
		member.ClubID = clubID

		if result := database.Db.Create(&member); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(member)

	case http.MethodDelete:
		memberID := r.URL.Query().Get("id")
		if memberID == "" {
			http.Error(w, "Member ID parameter is required", http.StatusBadRequest)
			return
		}

		result := database.Db.Where("id = ? AND club_id = ?", memberID, clubID).Delete(&models.Member{})
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
