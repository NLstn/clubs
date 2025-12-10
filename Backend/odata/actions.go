package odata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	"github.com/google/uuid"
	odata "github.com/nlstn/go-odata"
)

// Input validation helpers

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// isValidEmail validates email format
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// isValidUUID validates UUID format
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// isValidRole validates member role values
func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"member": true,
	}
	return validRoles[role]
}

// isValidRSVPResponse validates RSVP response values
func isValidRSVPResponse(response string) bool {
	validResponses := map[string]bool{
		"yes":   true,
		"no":    true,
		"maybe": true,
	}
	return validResponses[response]
}

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

	// Bound actions for Event entity - RSVP management
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:      "AddRSVP",
		IsBound:   true,
		EntitySet: "Events",
		Parameters: []odata.ParameterDefinition{
			{Name: "response", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: nil,
		Handler:    s.addRSVPAction,
	}); err != nil {
		return fmt.Errorf("failed to register AddRSVP action for Event: %w", err)
	}

	// Bound actions for Club entity - Additional operations
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:       "Join",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: nil,
		Handler:    s.joinClubAction,
	}); err != nil {
		return fmt.Errorf("failed to register Join action for Club: %w", err)
	}

	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:      "CreateInvite",
		IsBound:   true,
		EntitySet: "Clubs",
		Parameters: []odata.ParameterDefinition{
			{Name: "email", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: nil,
		Handler:    s.createInviteAction,
	}); err != nil {
		return fmt.Errorf("failed to register CreateInvite action for Club: %w", err)
	}

	// Bound actions for Member entity
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:      "UpdateRole",
		IsBound:   true,
		EntitySet: "Members",
		Parameters: []odata.ParameterDefinition{
			{Name: "newRole", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: nil,
		Handler:    s.updateMemberRoleAction,
	}); err != nil {
		return fmt.Errorf("failed to register UpdateRole action for Member: %w", err)
	}

	// Bound actions for Shift entity
	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:      "AddMember",
		IsBound:   true,
		EntitySet: "Shifts",
		Parameters: []odata.ParameterDefinition{
			{Name: "memberId", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: nil,
		Handler:    s.addShiftMemberAction,
	}); err != nil {
		return fmt.Errorf("failed to register AddMember action for Shift: %w", err)
	}

	if err := s.Service.RegisterAction(odata.ActionDefinition{
		Name:      "RemoveMember",
		IsBound:   true,
		EntitySet: "Shifts",
		Parameters: []odata.ParameterDefinition{
			{Name: "memberId", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: nil,
		Handler:    s.removeShiftMemberAction,
	}); err != nil {
		return fmt.Errorf("failed to register RemoveMember action for Shift: %w", err)
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

// addRSVPAction handles the AddRSVP action on Event entity
// POST /api/v2/Events('{eventId}')/AddRSVP
func (s *Service) addRSVPAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	event := ctx.(*models.Event)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get response parameter
	response, ok := params["response"].(string)
	if !ok {
		return fmt.Errorf("response parameter is required")
	}

	// Validate response value
	response = strings.TrimSpace(strings.ToLower(response))
	if !isValidRSVPResponse(response) {
		return fmt.Errorf("invalid response: must be 'yes', 'no', or 'maybe', got '%s'", response)
	}

	// Check if user is member of the club
	var club models.Club
	if err := s.db.Where("id = ?", event.ClubID).First(&club).Error; err != nil {
		return fmt.Errorf("failed to find club: %w", err)
	}

	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if !club.IsMember(user) {
		return fmt.Errorf("only club members can RSVP to events")
	}

	// Check if RSVP already exists
	var existingRSVP models.EventRSVP
	result := s.db.Where("event_id = ? AND user_id = ?", event.ID, userID).First(&existingRSVP)

	if result.Error == nil {
		// Update existing RSVP
		existingRSVP.Response = response
		if err := s.db.Save(&existingRSVP).Error; err != nil {
			return fmt.Errorf("failed to update RSVP: %w", err)
		}
	} else {
		// Create new RSVP
		newRSVP := models.EventRSVP{
			EventID:  event.ID,
			UserID:   userID,
			Response: response,
		}
		if err := s.db.Create(&newRSVP).Error; err != nil {
			return fmt.Errorf("failed to create RSVP: %w", err)
		}
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// joinClubAction handles the Join action on Club entity
// POST /api/v2/Clubs('{clubId}')/Join
func (s *Service) joinClubAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is already a member
	if club.IsMember(user) {
		return fmt.Errorf("you are already a member of this club")
	}

	// Create join request
	joinRequest := models.JoinRequest{
		ClubID: club.ID,
		UserID: userID,
	}

	if err := s.db.Create(&joinRequest).Error; err != nil {
		return fmt.Errorf("failed to create join request: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// createInviteAction handles the CreateInvite action on Club entity
// POST /api/v2/Clubs('{clubId}')/CreateInvite
func (s *Service) createInviteAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get email parameter
	email, ok := params["email"].(string)
	if !ok || email == "" {
		return fmt.Errorf("email parameter is required")
	}

	// Validate email format
	email = strings.TrimSpace(email)
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	// Get user for authorization
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is admin or owner
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		return fmt.Errorf("only club admins can send invites")
	}

	// Create invite using model function
	if err := club.CreateInvite(email, userID); err != nil {
		return fmt.Errorf("failed to create invite: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// updateMemberRoleAction handles the UpdateRole action on Member entity
// POST /api/v2/Members('{memberId}')/UpdateRole
func (s *Service) updateMemberRoleAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	member := ctx.(*models.Member)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get newRole parameter
	newRole, ok := params["newRole"].(string)
	if !ok || newRole == "" {
		return fmt.Errorf("newRole parameter is required")
	}

	// Validate role format
	newRole = strings.TrimSpace(strings.ToLower(newRole))
	if !isValidRole(newRole) {
		return fmt.Errorf("invalid role: must be 'owner', 'admin', or 'member', got '%s'", newRole)
	}

	// Get club
	var club models.Club
	if err := s.db.Where("id = ?", member.ClubID).First(&club).Error; err != nil {
		return fmt.Errorf("failed to find club: %w", err)
	}

	// Get current user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Get current user's role
	userRole, err := club.GetMemberRole(user)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Authorization checks
	if userRole != "owner" && userRole != "admin" {
		return fmt.Errorf("only club admins can change member roles")
	}

	// Owners can change any role
	// Admins can promote to admin or demote from admin, but cannot change owner roles
	if userRole == "admin" {
		if member.Role == "owner" || newRole == "owner" {
			return fmt.Errorf("admins cannot change owner roles")
		}
	}

	// Prevent last owner from being demoted
	if member.Role == "owner" && newRole != "owner" {
		ownerCount, err := club.CountOwners()
		if err != nil {
			return fmt.Errorf("failed to count owners: %w", err)
		}
		if ownerCount <= 1 {
			return fmt.Errorf("cannot demote the last owner")
		}
	}

	// Update the member's role
	member.Role = newRole
	if err := s.db.Save(member).Error; err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// addShiftMemberAction handles the AddMember action on Shift entity
// POST /api/v2/Shifts('{shiftId}')/AddMember
func (s *Service) addShiftMemberAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	shift := ctx.(*models.Shift)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get memberId parameter
	memberID, ok := params["memberId"].(string)
	if !ok || memberID == "" {
		return fmt.Errorf("memberId parameter is required")
	}

	// Validate UUID format
	if !isValidUUID(memberID) {
		return fmt.Errorf("invalid memberId format: must be a valid UUID")
	}

	// Get the event to find the club
	var event models.Event
	if err := s.db.Where("id = ?", shift.EventID).First(&event).Error; err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	// Get club
	var club models.Club
	if err := s.db.Where("id = ?", event.ClubID).First(&club).Error; err != nil {
		return fmt.Errorf("failed to find club: %w", err)
	}

	// Get current user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is admin or owner
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		return fmt.Errorf("only club admins can assign shift members")
	}

	// Verify the member exists and is part of the club
	var member models.Member
	if err := s.db.Where("id = ? AND club_id = ?", memberID, club.ID).First(&member).Error; err != nil {
		return fmt.Errorf("member not found in club")
	}

	// Check if already assigned
	var existing models.ShiftMember
	result := s.db.Where("shift_id = ? AND member_id = ?", shift.ID, memberID).First(&existing)
	if result.Error == nil {
		return fmt.Errorf("member is already assigned to this shift")
	}

	// Create shift member assignment
	shiftMember := models.ShiftMember{
		ShiftID:   shift.ID,
		UserID:    member.UserID,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := s.db.Create(&shiftMember).Error; err != nil {
		return fmt.Errorf("failed to assign member to shift: %w", err)
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// removeShiftMemberAction handles the RemoveMember action on Shift entity
// POST /api/v2/Shifts('{shiftId}')/RemoveMember
func (s *Service) removeShiftMemberAction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) error {
	shift := ctx.(*models.Shift)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get memberId parameter
	memberID, ok := params["memberId"].(string)
	if !ok || memberID == "" {
		return fmt.Errorf("memberId parameter is required")
	}

	// Validate UUID format
	if !isValidUUID(memberID) {
		return fmt.Errorf("invalid memberId format: must be a valid UUID")
	}

	// Get the event to find the club
	var event models.Event
	if err := s.db.Where("id = ?", shift.EventID).First(&event).Error; err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	// Get club
	var club models.Club
	if err := s.db.Where("id = ?", event.ClubID).First(&club).Error; err != nil {
		return fmt.Errorf("failed to find club: %w", err)
	}

	// Get current user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is admin or owner
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		return fmt.Errorf("only club admins can remove shift members")
	}

	// Get the member to find their UserID
	var member models.Member
	if err := s.db.Where("id = ? AND club_id = ?", memberID, club.ID).First(&member).Error; err != nil {
		return fmt.Errorf("member not found in club")
	}

	// Delete the shift member assignment using UserID
	result := s.db.Where("shift_id = ? AND user_id = ?", shift.ID, member.UserID).Delete(&models.ShiftMember{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove member from shift: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("member not found in shift")
	}

	w.Header().Set("OData-Version", "4.0")
	w.WriteHeader(http.StatusNoContent)
	return nil
}
