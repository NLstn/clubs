package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	t.Run("User Struct", func(t *testing.T) {
		user := User{
			ID:    "test-id",
			Name:  "Test User",
			Email: "test@example.com",
		}
		
		assert.Equal(t, "test-id", user.ID)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("Club Struct", func(t *testing.T) {
		club := Club{
			ID:          "club-id",
			Name:        "Test Club",
			Description: "Test Description",
			CreatedBy:   "creator-id",
		}
		
		assert.Equal(t, "club-id", club.ID)
		assert.Equal(t, "Test Club", club.Name)
		assert.Equal(t, "Test Description", club.Description)
		assert.Equal(t, "creator-id", club.CreatedBy)
	})

	t.Run("Member Struct", func(t *testing.T) {
		member := Member{
			ID:     "member-id",
			UserID: "user-id",
			ClubID: "club-id",
			Role:   "admin",
		}
		
		assert.Equal(t, "member-id", member.ID)
		assert.Equal(t, "user-id", member.UserID)
		assert.Equal(t, "club-id", member.ClubID)
		assert.Equal(t, "admin", member.Role)
	})

	t.Run("Fine Struct", func(t *testing.T) {
		fine := Fine{
			ID:     "fine-id",
			UserID: "user-id",
			ClubID: "club-id",
			Reason: "Late arrival",
			Amount: 25.50,
			Paid:   false,
		}
		
		assert.Equal(t, "fine-id", fine.ID)
		assert.Equal(t, "user-id", fine.UserID)
		assert.Equal(t, "club-id", fine.ClubID)
		assert.Equal(t, "Late arrival", fine.Reason)
		assert.Equal(t, 25.50, fine.Amount)
		assert.False(t, fine.Paid)
	})

	t.Run("JoinRequest Struct", func(t *testing.T) {
		request := JoinRequest{
			ID:     "request-id",
			ClubID: "club-id",
			Email:  "requester@example.com",
		}
		
		assert.Equal(t, "request-id", request.ID)
		assert.Equal(t, "club-id", request.ClubID)
		assert.Equal(t, "requester@example.com", request.Email)
	})

	t.Run("RefreshToken Struct", func(t *testing.T) {
		token := RefreshToken{
			ID:     "token-id",
			UserID: "user-id",
			Token:  "refresh-token-value",
		}
		
		assert.Equal(t, "token-id", token.ID)
		assert.Equal(t, "user-id", token.UserID)
		assert.Equal(t, "refresh-token-value", token.Token)
	})

	t.Run("GetUsersByIDs - Empty Array", func(t *testing.T) {
		users, err := GetUsersByIDs([]string{})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})
}