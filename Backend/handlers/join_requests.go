package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
)

func handleJoinRequestCreate(w http.ResponseWriter, r *http.Request) {

	var joinRequest models.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)
	if !auth.IsOwnerOfClub(userID, joinRequest.ClubID) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if joinRequest.ClubID == "" || joinRequest.Email == "" {
		http.Error(w, "Missing club_id or email", http.StatusBadRequest)
		return
	}

	err := models.CreateJoinRequest(joinRequest.ClubID, joinRequest.Email)
	if err != nil {
		http.Error(w, "Failed to create join request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
