package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func handleClubs(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(auth.UserIDKey).(string)

	switch r.Method {
	case http.MethodGet:
		path := strings.Trim(r.URL.Path, "/")
		segments := strings.Split(path, "/")

		if len(segments) == 4 && segments[3] != "" {
			id := segments[3]
			var club models.Club
			result := database.Db.First(&club, "id = ?", id)

			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Club not found", http.StatusNotFound)
				return
			}
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			if !auth.IsAuthorizedForClub(userID, club.ID) {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(club)
			return
		}

		var clubs []models.Club
		if result := database.Db.Find(&clubs); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clubs)

	case http.MethodPost:
		var club models.Club
		if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		club.ID = uuid.New().String()

		if result := database.Db.Create(&club); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(club)
	}
}
