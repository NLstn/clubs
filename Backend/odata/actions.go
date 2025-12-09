package odata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	odata "github.com/nlstn/go-odata"
)

// registerActions registers all OData bound and unbound actions
// Actions are POST operations that can have side effects
func (s *Service) registerActions() error {
	// Bound actions for Invite entity
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "Accept",
		IsBound:    true,
		EntitySet:  "Invites",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.acceptInviteAction,
	}); err != nil {
		return fmt.Errorf("failed to register Accept action for Invite: %w", err)
	}

	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "Reject",
		IsBound:    true,
		EntitySet:  "Invites",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.rejectInviteAction,
	}); err != nil {
		return fmt.Errorf("failed to register Reject action for Invite: %w", err)
	}

	// Bound actions for JoinRequest entity
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "Accept",
		IsBound:    true,
		EntitySet:  "JoinRequests",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.acceptJoinRequestAction,
	}); err != nil {
		return fmt.Errorf("failed to register Accept action for JoinRequest: %w", err)
	}

	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "Reject",
		IsBound:    true,
		EntitySet:  "JoinRequests",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.rejectJoinRequestAction,
	}); err != nil {
		return fmt.Errorf("failed to register Reject action for JoinRequest: %w", err)
	}

	// Bound actions for Club entity
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "Leave",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.leaveClubAction,
	}); err != nil {
		return fmt.Errorf("failed to register Leave action for Club: %w", err)
	}

	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "DeleteLogo",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.deleteLogoAction,
	}); err != nil {
		return fmt.Errorf("failed to register DeleteLogo action for Club: %w", err)
	}

	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "HardDelete",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.hardDeleteClubAction,
	}); err != nil {
		return fmt.Errorf("failed to register HardDelete action for Club: %w", err)
	}

	// Bound actions for Notification entity
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "MarkAsRead",
		IsBound:    true,
		EntitySet:  "Notifications",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.markNotificationReadAction,
	}); err != nil {
		return fmt.Errorf("failed to register MarkAsRead action for Notification: %w", err)
	}

	// Unbound action
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "MarkAllNotificationsRead",
		IsBound:    false,
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.markAllNotificationsReadAction,
	}); err != nil {
		return fmt.Errorf("failed to register MarkAllNotificationsRead action: %w", err)
	}

	return nil
}

// acceptInviteAction handles the Accept action on Invite entity
// POST /api/v2/Invites('{inviteId}')/Accept
func (s *Service) acceptInviteAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	invite := ctx.(*models.Invite)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Verify the invite is for this user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	canEdit, err := user.CanUserEditInvite(invite.ID)
	if err != nil || !canEdit {
		return fmt.Errorf("unauthorized: invite is not for this user")
	}

	// Accept the invite using model function
	if err := models.AcceptInvite(invite.ID, userID); err != nil {
		return fmt.Errorf("failed to accept invite: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// rejectInviteAction handles the Reject action on Invite entity
// POST /api/v2/Invites('{inviteId}')/Reject
func (s *Service) rejectInviteAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	invite := ctx.(*models.Invite)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Verify the invite is for this user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	canEdit, err := user.CanUserEditInvite(invite.ID)
	if err != nil || !canEdit {
		return fmt.Errorf("unauthorized: invite is not for this user")
	}

	// Reject the invite using model function
	if err := models.RejectInvite(invite.ID); err != nil {
		return fmt.Errorf("failed to reject invite: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// acceptJoinRequestAction handles the Accept action on JoinRequest entity
// POST /api/v2/JoinRequests('{requestId}')/Accept
func (s *Service) acceptJoinRequestAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	joinRequest := ctx.(*models.JoinRequest)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Accept the join request using model function (it handles authorization internally)
	if err := models.AcceptJoinRequest(joinRequest.ID, userID); err != nil {
		return fmt.Errorf("failed to accept join request: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// rejectJoinRequestAction handles the Reject action on JoinRequest entity
// POST /api/v2/JoinRequests('{requestId}')/Reject
func (s *Service) rejectJoinRequestAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	joinRequest := ctx.(*models.JoinRequest)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Reject the join request using model function (it handles authorization internally)
	if err := models.RejectJoinRequest(joinRequest.ID, userID); err != nil {
		return fmt.Errorf("failed to reject join request: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// leaveClubAction handles the Leave action on Club entity
// POST /api/v2/Clubs('{clubId}')/Leave
func (s *Service) leaveClubAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user for authorization checks
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is a member of the club
	if !club.IsMember(user) {
		return fmt.Errorf("you are not a member of this club")
	}

	// Get the user's member record to check their role
	userRole, err := club.GetMemberRole(user)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if user is the last owner - prevent leaving if so
	if userRole == "owner" {
		ownerCount, err := club.CountOwners()
		if err != nil {
			return fmt.Errorf("failed to check owner count: %w", err)
		}
		if ownerCount <= 1 {
			return fmt.Errorf("cannot leave club: you are the last owner. Transfer ownership or delete the club first")
		}
	}

	// Find the user's member record and delete it
	if err := club.DeleteMemberByUserID(userID); err != nil {
		return fmt.Errorf("failed to leave club: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// deleteLogoAction handles the DeleteLogo action on Club entity
// POST /api/v2/Clubs('{clubId}')/DeleteLogo
func (s *Service) deleteLogoAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user for authorization checks
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is admin or owner
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		return fmt.Errorf("unauthorized: only club admins can delete club logo")
	}

	// Delete the logo (LogoURL is a pointer to string)
	club.LogoURL = nil
	if err := s.db.Save(club).Error; err != nil {
		return fmt.Errorf("failed to delete logo: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// hardDeleteClubAction handles the HardDelete action on Club entity
// POST /api/v2/Clubs('{clubId}')/HardDelete
func (s *Service) hardDeleteClubAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user for authorization checks
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is owner (only owners can hard delete)
	if !club.IsOwner(user) {
		return fmt.Errorf("unauthorized: only club owners can hard delete clubs")
	}

	// Hard delete the club (permanently delete, bypassing soft delete)
	if err := s.db.Unscoped().Delete(club).Error; err != nil {
		return fmt.Errorf("failed to hard delete club: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// markNotificationReadAction handles the MarkAsRead action on Notification entity
// POST /api/v2/Notifications('{notificationId}')/MarkAsRead
func (s *Service) markNotificationReadAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	notification := ctx.(*models.Notification)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Verify the notification belongs to this user
	if notification.UserID != userID {
		return fmt.Errorf("unauthorized: notification does not belong to this user")
	}

	// Mark as read
	notification.Read = true
	if err := s.db.Save(notification).Error; err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// markAllNotificationsReadAction handles the unbound MarkAllNotificationsRead action
// POST /api/v2/MarkAllNotificationsRead
func (s *Service) markAllNotificationsReadAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Mark all notifications as read for this user
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error; err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	// Return the count of notifications marked as read
	var count int64
	s.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, true).
		Count(&count)

	response := map[string]interface{}{
		"@odata.context": "$metadata#Edm.Int64",
		"value":          count,
	}

	w.Header().Set("Content-Type", "application/json;odata.metadata=minimal")
	w.Header().Set("OData-Version", "4.0")
	return json.NewEncoder(w).Encode(response)
}
