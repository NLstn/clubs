package odata

import (
	"context"
	"fmt"

	"github.com/NLstn/clubs/auth"
)

// getUserIDFromContext extracts the user ID from the request context
// This is called by OData query handlers to determine access permissions
func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// AuthorizationModel defines the permission rules for each entity type.
// These rules are enforced through custom OData action handlers and middleware.
//
// AUTHORIZATION RULES:
//
// Users:
// - Users can only read themselves
// - Users can only update themselves
//
// Clubs:
// - Users can read clubs they're members of (non-deleted) or created (all states)
// - Users can create clubs
// - Only admins/owners can update clubs
// - Only owners can delete clubs
//
// Members:
// - Users can read members of clubs they're members of
// - Only club admins can create/update/delete members
//
// Teams:
// - Users can read teams in clubs they're members of (non-deleted) or created
// - Only club admins can create teams
// - Only admins can update/delete teams
//
// Events:
// - Users can read events in clubs they're members of
// - Only club admins can create/update/delete events
//
// Shifts:
// - Users can read shifts in clubs they're members of
// - Only club admins can create/update/delete shifts
//
// Fines:
// - Users can read fines in their clubs or their own fines
// - Only club admins can create/delete fines
//
// News:
// - Users can read news in clubs they're members of
// - Only club admins can create/update news
//
// Notifications:
// - Users can only read their own notifications
// - Users can only update their own notifications (mark as read)
// - Only backend can create notifications
//
// Invites:
// - Users can read invites sent to them or sent by them
// - Only club admins can create invites
// - Users can accept/reject invites sent to them
//
// JoinRequests:
// - Users can read their own join requests
// - Club admins can read join requests for their clubs
// - Only club admins can accept/reject join requests
//
// Privacy Settings & Preferences:
// - Users can only read/update their own settings
//
// Implementation Details for Phase 3:
// These rules will be enforced through OData entity lifecycle hooks:
//
// 1. BeforeReadCollection hooks filter collections by user context
//    Example: Club.BeforeReadCollection filters to clubs where user is member
//
// 2. BeforeReadEntity hooks validate access to specific entities
//    Example: Club.BeforeReadEntity checks membership before returning single club
//
// 3. BeforeCreate hooks validate write permissions before insertion
//    Example: Event.BeforeCreate checks user is admin of target club
//
// 4. BeforeUpdate hooks validate update permissions
//    Example: Club.BeforeUpdate checks user is admin/owner
//
// 5. BeforeDelete hooks validate delete permissions
//    Example: Club.BeforeDelete checks user is owner
//
// See: Documentation/Backend/OData_Hooks_Guide.md for complete implementation guide
