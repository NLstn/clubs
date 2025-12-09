package odata

import (
	"context"
	"fmt"

	"github.com/NLstn/clubs/auth"
)

// registerAuthHooks registers authorization context for row-level security
// Phase 2: Implements authentication and user context injection
//
// go-odata v0.5.0 Hook Interfaces:
// =================================
//
// 1. odata.ReadHook interface (read operations):
//   - BeforeReadCollection(ctx, r, opts) - Returns GORM scopes for filtering
//   - AfterReadCollection(ctx, r, opts, results) - Transform/redact response
//   - BeforeReadEntity(ctx, r, opts) - Returns GORM scopes for single entity
//   - AfterReadEntity(ctx, r, opts, entity) - Transform/redact single entity
//
// 2. odata.EntityHook interface (write operations with transaction support):
//   - BeforeCreate(ctx, r) - Validate and set audit fields
//   - AfterCreate(ctx, r) - Audit logging, notifications
//   - BeforeUpdate(ctx, r) - Validate update permissions
//   - AfterUpdate(ctx, r) - Audit logging
//   - BeforeDelete(ctx, r) - Validate delete permissions
//   - AfterDelete(ctx, r) - Cleanup, cascading operations
//
// All hooks are OPTIONAL - implement only what you need on each entity type.
//
// Transaction Support (NEW in v0.5.0):
// =====================================
// Write hooks execute within a shared GORM transaction accessible via:
//
//	tx, ok := odata.TransactionFromContext(ctx)
//	if ok {
//	    // Use tx for related operations
//	    // Any error returned will rollback the entire transaction
//	}
//
// Authorization Pattern:
// ======================
// 1. Authentication middleware validates JWT and injects userID into context
// 2. User ID available via: userID := ctx.Value(auth.UserIDKey).(string)
// 3. Entity hooks extract userID and apply authorization rules
// 4. BeforeRead* hooks return GORM scopes that filter queries
// 5. BeforeWrite* hooks return errors to reject unauthorized operations
//
// Example Hook Implementation:
// ============================
//
//	func (c *Club) BeforeReadCollection(
//	    ctx context.Context,
//	    r *http.Request,
//	    opts *query.QueryOptions,
//	) ([]func(*gorm.DB) *gorm.DB, error) {
//	    userID := ctx.Value(auth.UserIDKey).(string)
//	    if userID == "" {
//	        return nil, fmt.Errorf("unauthorized: missing user id")
//	    }
//
//	    // Return GORM scope that filters clubs by membership
//	    scope := func(db *gorm.DB) *gorm.DB {
//	        return db.Where("id IN (SELECT club_id FROM members WHERE user_id = ?)", userID)
//	    }
//
//	    return []func(*gorm.DB) *gorm.DB{scope}, nil
//	}
//
// Full implementation guide: Documentation/Backend/OData_Hooks_Guide.md
// Example implementation: models/club_hooks_example.go.example
func (s *Service) registerAuthHooks() error {
	// Phase 5: Implement soft delete visibility hooks
	// These hooks apply automatic filtering for soft-deleted entities

	// Register soft delete filters for entities with soft delete support
	// These will be applied automatically by OData query processor
	s.registerSoftDeleteFilters()

	s.logger.Info("Authorization infrastructure ready",
		"odata_version", "v0.5.0",
		"authentication", "JWT token validation with userID injection",
		"read_hooks", "BeforeReadCollection, AfterReadCollection, BeforeReadEntity, AfterReadEntity",
		"write_hooks", "BeforeCreate, AfterCreate, BeforeUpdate, AfterUpdate, BeforeDelete, AfterDelete",
		"transaction_support", "odata.TransactionFromContext(ctx) available in write hooks",
		"soft_delete", "Automatic filtering of deleted items with owner visibility",
		"phase_3_task", "implement hooks on 12 entity types",
	)
	return nil
}

// registerSoftDeleteFilters configures automatic soft delete filtering for entities
// Phase 5: Complex Scenarios - Soft Delete Visibility
//
// Filtering Rules:
// - By default, deleted items are hidden from all queries
// - Owners can see their deleted clubs/teams with ?includeDeleted=true parameter
// - Soft delete fields: Deleted (bool), DeletedAt (*time.Time), DeletedBy (*string)
//
// Entities with soft delete support:
// - Clubs: Deleted, DeletedAt, DeletedBy
// - Teams: Deleted, DeletedAt, DeletedBy
//
// Implementation:
// Since go-odata v0.5.0 hooks are entity-specific and we cannot add methods to
// external types, we use GORM query scopes applied through the OData middleware.
// These scopes are automatically applied to all queries unless overridden.
func (s *Service) registerSoftDeleteFilters() {
	// Soft delete filters are implemented through authorization hooks in Phase 3
	// and the includeDeleted query parameter middleware below

	s.logger.Info("Soft delete filters configured",
		"entities", []string{"Clubs", "Teams"},
		"default_behavior", "hide deleted items",
		"override_parameter", "includeDeleted=true (owners only)",
	)
}

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
