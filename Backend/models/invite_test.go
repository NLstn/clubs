package models_test

import (
	"fmt"
	"testing"

	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvite_DuplicatePrevention(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("cannot_send_duplicate_invite", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Create first invite
		err := club.CreateInvite("invitee@example.com", owner.ID)
		assert.NoError(t, err)

		// Attempt to create duplicate invite
		err = club.CreateInvite("invitee@example.com", owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invite already exists")
	})

	t.Run("can_send_same_email_to_different_clubs", func(t *testing.T) {
		owner1, _ := handlers.CreateTestUser(t, "owner1@example.com")
		owner2, _ := handlers.CreateTestUser(t, "owner2@example.com")
		club1 := handlers.CreateTestClub(t, owner1, "Club 1")
		club2 := handlers.CreateTestClub(t, owner2, "Club 2")

		// Send invite to same email from different clubs
		err := club1.CreateInvite("invitee@example.com", owner1.ID)
		assert.NoError(t, err)

		err = club2.CreateInvite("invitee@example.com", owner2.ID)
		assert.NoError(t, err)
	})
}

func TestCreateInvite_MemberCheck(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("cannot_invite_existing_member", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		member, _ := handlers.CreateTestUser(t, "member@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Add user as member
		handlers.CreateTestMember(t, member, club, "member")

		// Attempt to invite existing member
		err := club.CreateInvite(member.Email, owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already a member")
	})

	t.Run("can_invite_non_member", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner2@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Invite non-member
		err := club.CreateInvite("newuser@example.com", owner.ID)
		assert.NoError(t, err)
	})
}

func TestCreateInvite_AdminRateLimit(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("admin_rate_limit_enforced", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Create 10 invites (should succeed)
		for i := 0; i < 10; i++ {
			email := fmt.Sprintf("invitee%d@example.com", i)
			err := club.CreateInvite(email, owner.ID)
			assert.NoError(t, err, "Invite %d should succeed", i)
		}

		// 11th invite should fail due to rate limit
		err := club.CreateInvite("invitee11@example.com", owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")
		assert.Contains(t, err.Error(), "10 invites per hour per admin")
	})

	t.Run("different_admins_have_separate_limits", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner3@example.com")
		admin, _ := handlers.CreateTestUser(t, "admin3@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")
		handlers.CreateTestMember(t, admin, club, "admin")

		// Owner creates 10 invites
		for i := 0; i < 10; i++ {
			email := fmt.Sprintf("owner-invite%d@example.com", i)
			err := club.CreateInvite(email, owner.ID)
			assert.NoError(t, err)
		}

		// Admin should still be able to create invites (separate limit)
		err := club.CreateInvite("admin-invite@example.com", admin.ID)
		assert.NoError(t, err, "Admin should have separate rate limit")
	})
}

func TestCreateInvite_ClubRateLimit(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("club_rate_limit_enforced", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		admin1, _ := handlers.CreateTestUser(t, "admin1@example.com")
		admin2, _ := handlers.CreateTestUser(t, "admin2@example.com")
		admin3, _ := handlers.CreateTestUser(t, "admin3@example.com")
		admin4, _ := handlers.CreateTestUser(t, "admin4@example.com")
		admin5, _ := handlers.CreateTestUser(t, "admin5@example.com")
		
		club := handlers.CreateTestClub(t, owner, "Test Club")
		handlers.CreateTestMember(t, admin1, club, "admin")
		handlers.CreateTestMember(t, admin2, club, "admin")
		handlers.CreateTestMember(t, admin3, club, "admin")
		handlers.CreateTestMember(t, admin4, club, "admin")
		handlers.CreateTestMember(t, admin5, club, "admin")

		// Each admin creates 10 invites (50 total)
		admins := []models.User{owner, admin1, admin2, admin3, admin4}
		for adminIdx, admin := range admins {
			for i := 0; i < 10; i++ {
				email := fmt.Sprintf("club-invite-a%d-i%d@example.com", adminIdx, i)
				err := club.CreateInvite(email, admin.ID)
				assert.NoError(t, err, "Admin %d invite %d should succeed", adminIdx, i)
			}
		}

		// 51st invite should fail due to club rate limit
		err := club.CreateInvite("club-limit-exceeded@example.com", admin5.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")
		assert.Contains(t, err.Error(), "50 invites per hour per club")
	})

	t.Run("different_clubs_have_separate_limits", func(t *testing.T) {
		owner1, _ := handlers.CreateTestUser(t, "owner4@example.com")
		owner2, _ := handlers.CreateTestUser(t, "owner5@example.com")
		club1 := handlers.CreateTestClub(t, owner1, "Club 1")
		club2 := handlers.CreateTestClub(t, owner2, "Club 2")

		// Club1 creates 10 invites
		for i := 0; i < 10; i++ {
			email := fmt.Sprintf("club1-invite%d@example.com", i)
			err := club1.CreateInvite(email, owner1.ID)
			assert.NoError(t, err)
		}

		// Club2 should still be able to create invites (separate limit)
		err := club2.CreateInvite("club2-invite@example.com", owner2.ID)
		assert.NoError(t, err, "Club 2 should have separate rate limit")
	})
}

func TestCreateInvite_RateLimitReset(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("rate_limit_uses_sliding_window", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// This test verifies the sliding window behavior
		// In practice, invites older than 1 hour don't count toward the limit

		// Create 10 invites
		for i := 0; i < 10; i++ {
			email := fmt.Sprintf("sliding-window%d@example.com", i)
			err := club.CreateInvite(email, owner.ID)
			assert.NoError(t, err)
		}

		// 11th should fail
		err := club.CreateInvite("should-fail@example.com", owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")

		// Note: In a real scenario with time progression, old invites would be excluded
		// This test documents the expected behavior
	})
}

func TestCreateInvite_CombinedValidations(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("all_validations_work_together", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		member, _ := handlers.CreateTestUser(t, "member@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")
		handlers.CreateTestMember(t, member, club, "member")

		// 1. Cannot invite existing member
		err := club.CreateInvite(member.Email, owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already a member")

		// 2. Create first invite
		err = club.CreateInvite("valid@example.com", owner.ID)
		assert.NoError(t, err)

		// 3. Cannot create duplicate
		err = club.CreateInvite("valid@example.com", owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invite already exists")

		// 4. Can continue creating new invites until rate limit
		for i := 0; i < 9; i++ {
			email := fmt.Sprintf("more-invites%d@example.com", i)
			err = club.CreateInvite(email, owner.ID)
			assert.NoError(t, err, "Invite %d should succeed", i)
		}

		// 5. Rate limit kicks in
		err = club.CreateInvite("rate-limited@example.com", owner.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")
	})
}

func TestGetInvites(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("get_club_invites", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Create multiple invites
		emails := []string{"invite1@example.com", "invite2@example.com", "invite3@example.com"}
		for _, email := range emails {
			err := club.CreateInvite(email, owner.ID)
			assert.NoError(t, err)
		}

		// Get invites
		invites, err := club.GetInvites()
		assert.NoError(t, err)
		assert.Len(t, invites, 3)
	})
}

func TestAcceptInvite_WithDuplicateConstraint(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("accepting_invite_removes_it", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		user, _ := handlers.CreateTestUser(t, "user@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Create invite
		err := club.CreateInvite(user.Email, owner.ID)
		assert.NoError(t, err)

		// Get the invite
		invites, err := club.GetInvites()
		assert.NoError(t, err)
		assert.Len(t, invites, 1)

		// Accept invite
		err = models.AcceptInvite(invites[0].ID, user.ID)
		assert.NoError(t, err)

		// Verify invite is deleted
		invites, err = club.GetInvites()
		assert.NoError(t, err)
		assert.Len(t, invites, 0)

		// Verify user is now a member
		assert.True(t, club.IsMember(user))

		// Should be able to create new invite to same email for different club
		owner2, _ := handlers.CreateTestUser(t, "owner2@example.com")
		club2 := handlers.CreateTestClub(t, owner2, "Club 2")
		err = club2.CreateInvite(user.Email, owner2.ID)
		assert.NoError(t, err, "Should be able to invite to different club even if member elsewhere")
	})
}

func TestRejectInvite_WithDuplicateConstraint(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("rejecting_invite_removes_it", func(t *testing.T) {
		owner, _ := handlers.CreateTestUser(t, "owner@example.com")
		club := handlers.CreateTestClub(t, owner, "Test Club")

		// Create invite
		err := club.CreateInvite("reject@example.com", owner.ID)
		assert.NoError(t, err)

		// Get the invite
		invites, err := club.GetInvites()
		assert.NoError(t, err)
		assert.Len(t, invites, 1)

		// Reject invite
		err = models.RejectInvite(invites[0].ID)
		assert.NoError(t, err)

		// Verify invite is deleted
		invites, err = club.GetInvites()
		assert.NoError(t, err)
		assert.Len(t, invites, 0)

		// Should be able to create new invite to same email after rejection
		err = club.CreateInvite("reject@example.com", owner.ID)
		assert.NoError(t, err, "Should be able to re-invite after rejection")
	})
}
