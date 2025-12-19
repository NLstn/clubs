package odata

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	odata "github.com/nlstn/go-odata"
)

// registerFunctions registers all OData bound and unbound functions
// Functions are GET operations that return computed values without side effects
func (s *Service) registerFunctions() error {
	// Bound functions for Club entity
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "IsAdmin",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf(map[string]bool{}),
		Handler:    s.isAdminFunction,
	}); err != nil {
		return fmt.Errorf("failed to register IsAdmin function for Club: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetOwnerCount",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf(map[string]int64{}),
		Handler:    s.getOwnerCountFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetOwnerCount function for Club: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetInviteLink",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf(map[string]string{}),
		Handler:    s.getInviteLinkFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetInviteLink function for Club: %w", err)
	}

	// Bound functions for Event entity
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:      "ExpandRecurrence",
		IsBound:   true,
		EntitySet: "Events",
		Parameters: []odata.ParameterDefinition{
			{Name: "startDate", Type: reflect.TypeOf(time.Time{}), Required: true},
			{Name: "endDate", Type: reflect.TypeOf(time.Time{}), Required: true},
		},
		ReturnType: reflect.TypeOf([]models.Event{}),
		Handler:    s.expandRecurrenceFunction,
	}); err != nil {
		return fmt.Errorf("failed to register ExpandRecurrence function for Event: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetRSVPCounts",
		IsBound:    true,
		EntitySet:  "Events",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf(map[string]int64{}),
		Handler:    s.getRSVPCountsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetRSVPCounts function for Event: %w", err)
	}

	// Unbound functions
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:    "SearchGlobal",
		IsBound: false,
		Parameters: []odata.ParameterDefinition{
			{Name: "query", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: reflect.TypeOf(SearchResponse{}),
		Handler:    s.searchGlobalFunction,
	}); err != nil {
		return fmt.Errorf("failed to register SearchGlobal function: %w", err)
	}

	// Bound functions for Team entity
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetOverview",
		IsBound:    true,
		EntitySet:  "Teams",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf(TeamOverviewResponse{}),
		Handler:    s.getTeamOverviewFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetOverview function for Team: %w", err)
	}

	return nil
}

// Helper types for search functions
type SearchResult struct {
	Type        string `json:"Type"`
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	ClubID      string `json:"ClubID,omitempty"`
	ClubName    string `json:"ClubName,omitempty"`
	StartTime   string `json:"StartTime,omitempty"`
	EndTime     string `json:"EndTime,omitempty"`
}

type SearchResponse struct {
	Clubs  []SearchResult `json:"Clubs"`
	Events []SearchResult `json:"Events"`
}

type TeamOverviewResponse struct {
	Team     models.Team            `json:"Team"`
	Stats    map[string]interface{} `json:"Stats"`
	UserRole string                 `json:"UserRole"`
	IsAdmin  bool                   `json:"IsAdmin"`
}

type EventWithRSVP struct {
	models.Event
	UserRSVP *models.EventRSVP `json:"UserRSVP,omitempty"`
}

// isAdminFunction checks if the current user is an admin of the club
// GET /api/v2/Clubs('{clubId}')/IsAdmin()
func (s *Service) isAdminFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	isAdmin := club.IsAdmin(user) || club.IsOwner(user)

	return map[string]bool{"IsAdmin": isAdmin}, nil
}

// getOwnerCountFunction returns the number of owners in the club
// GET /api/v2/Clubs('{clubId}')/GetOwnerCount()
func (s *Service) getOwnerCountFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	club := ctx.(*models.Club)

	count, err := club.CountOwners()
	if err != nil {
		return nil, fmt.Errorf("failed to count owners: %w", err)
	}

	return map[string]int64{"OwnerCount": count}, nil
}

// getInviteLinkFunction returns the invite link for the club
// GET /api/v2/Clubs('{clubId}')/GetInviteLink()
func (s *Service) getInviteLinkFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	club := ctx.(*models.Club)

	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user for authorization
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is admin or owner
	if !club.IsOwner(user) && !club.IsAdmin(user) {
		return nil, fmt.Errorf("unauthorized: only club admins can get invite link")
	}

	// Return invite link
	inviteLink := "/join/" + club.ID

	return map[string]string{"InviteLink": inviteLink}, nil
}

// searchGlobalFunction performs a global search across clubs and events
// GET /api/v2/SearchGlobal(query='search term')
func (s *Service) searchGlobalFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get query parameter
	query, ok := params["query"].(string)
	if !ok || strings.TrimSpace(query) == "" {
		return SearchResponse{
			Clubs:  []SearchResult{},
			Events: []SearchResult{},
		}, nil
	}

	// Search clubs
	clubResults, err := s.searchClubs(user, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search clubs: %w", err)
	}

	// Search events
	eventResults, err := s.searchEvents(user, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}

	return SearchResponse{
		Clubs:  clubResults,
		Events: eventResults,
	}, nil
}

func (s *Service) searchClubs(user models.User, query string) ([]SearchResult, error) {
	clubs, err := models.GetAllClubsIncludingDeleted()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	queryLower := strings.ToLower(query)

	for _, club := range clubs {
		// Only show clubs the user is a member of
		if !club.IsMember(user) {
			continue
		}

		// Skip deleted clubs unless user is owner
		if club.Deleted && !club.IsOwner(user) {
			continue
		}

		// Check if club name or description contains the query
		description := ""
		if club.Description != nil {
			description = *club.Description
		}

		if strings.Contains(strings.ToLower(club.Name), queryLower) ||
			strings.Contains(strings.ToLower(description), queryLower) {
			results = append(results, SearchResult{
				Type:        "club",
				ID:          club.ID,
				Name:        club.Name,
				Description: description,
			})
		}
	}

	return results, nil
}

// SearchClubsForTest exposes searchClubs for testing
func (s *Service) SearchClubsForTest(user models.User, query string) ([]SearchResult, error) {
	return s.searchClubs(user, query)
}

func (s *Service) searchEvents(user models.User, query string) ([]SearchResult, error) {
	// Get events from clubs the user is a member of
	clubs, err := models.GetAllClubs()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	queryLower := strings.ToLower(query)

	for _, club := range clubs {
		if !club.IsMember(user) {
			continue
		}

		events, err := club.GetUpcomingEvents()
		if err != nil {
			continue
		}

		for _, event := range events {
			description := ""
			if event.Description != nil {
				description = *event.Description
			}

			if strings.Contains(strings.ToLower(event.Name), queryLower) ||
				strings.Contains(strings.ToLower(description), queryLower) {
				result := SearchResult{
					Type:        "event",
					ID:          event.ID,
					Name:        event.Name,
					Description: description,
					ClubID:      event.ClubID,
					ClubName:    club.Name,
				}

				if !event.StartTime.IsZero() {
					result.StartTime = event.StartTime.Format("2006-01-02T15:04:05Z07:00")
				}
				if !event.EndTime.IsZero() {
					result.EndTime = event.EndTime.Format("2006-01-02T15:04:05Z07:00")
				}

				results = append(results, result)
			}
		}
	}

	return results, nil
}

// SearchEventsForTest exposes searchEvents for testing
func (s *Service) SearchEventsForTest(user models.User, query string) ([]SearchResult, error) {
	return s.searchEvents(user, query)
}

// expandRecurrenceFunction generates recurring event instances for a given date range
// GET /api/v2/Events('{eventId}')/ExpandRecurrence(startDate=2024-01-01T00:00:00Z,endDate=2024-12-31T23:59:59Z)
//
// This function takes a recurring event pattern and expands it into individual event instances
// within the specified date range. This is useful for displaying recurring events in calendars
// without creating all instances in the database upfront.
//
// Parameters:
// - startDate: The start of the date range to generate instances for
// - endDate: The end of the date range to generate instances for
//
// Returns:
// - Array of Event objects representing the expanded instances
//
// Authorization: User must be a member of the club that owns the event
func (s *Service) expandRecurrenceFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	event := ctx.(*models.Event)

	// Extract parameters
	startDate, ok := params["startDate"].(time.Time)
	if !ok {
		return nil, fmt.Errorf("startDate parameter is required")
	}

	endDate, ok := params["endDate"].(time.Time)
	if !ok {
		return nil, fmt.Errorf("endDate parameter is required")
	}

	// Get user ID from request context
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: missing user id")
	}

	// Get user from database
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get club and verify user has access
	club, err := models.GetClubByID(event.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !club.IsMember(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this club")
	}

	// Check if event is recurring
	if !event.IsRecurring {
		// If not recurring, just return the event itself if it falls within range
		if (event.StartTime.After(startDate) || event.StartTime.Equal(startDate)) &&
			(event.StartTime.Before(endDate) || event.StartTime.Equal(endDate)) {
			return []models.Event{*event}, nil
		}
		// Event is outside the requested range
		return []models.Event{}, nil
	}

	// Generate recurring instances
	instances, err := generateRecurringInstances(event, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recurring instances: %w", err)
	}

	// Return the expanded instances
	return instances, nil
}

// generateRecurringInstances generates event instances from a recurring pattern
func generateRecurringInstances(parentEvent *models.Event, startDate, endDate time.Time) ([]models.Event, error) {
	if parentEvent.RecurrencePattern == nil || *parentEvent.RecurrencePattern == "" {
		return nil, fmt.Errorf("event has no recurrence pattern")
	}

	pattern := *parentEvent.RecurrencePattern
	interval := parentEvent.RecurrenceInterval
	if interval < 1 {
		interval = 1
	}

	var instances []models.Event
	currentStart := parentEvent.StartTime
	duration := parentEvent.EndTime.Sub(parentEvent.StartTime)

	// If parent event is within range, include it
	if (currentStart.After(startDate) || currentStart.Equal(startDate)) &&
		(currentStart.Before(endDate) || currentStart.Equal(endDate)) {
		instances = append(instances, *parentEvent)
	}

	// Generate instances until we reach the end date or recurrence end
	maxEnd := endDate
	if parentEvent.RecurrenceEnd != nil && parentEvent.RecurrenceEnd.Before(endDate) {
		maxEnd = *parentEvent.RecurrenceEnd
	}

	// Start from the next occurrence after the parent event
	currentStart = calculateNextOccurrence(currentStart, pattern, interval)

	for currentStart.Before(maxEnd) || currentStart.Equal(maxEnd) {
		currentEnd := currentStart.Add(duration)

		// Only include if within requested range
		if (currentStart.After(startDate) || currentStart.Equal(startDate)) &&
			(currentStart.Before(endDate) || currentStart.Equal(endDate)) {

			// Create instance (not saved to DB, just for response)
			instance := models.Event{
				ID:            fmt.Sprintf("%s-%s", parentEvent.ID, currentStart.Format("20060102T150405")),
				ClubID:        parentEvent.ClubID,
				TeamID:        parentEvent.TeamID,
				Name:          parentEvent.Name,
				Description:   parentEvent.Description,
				Location:      parentEvent.Location,
				StartTime:     currentStart,
				EndTime:       currentEnd,
				CreatedBy:     parentEvent.CreatedBy,
				UpdatedBy:     parentEvent.UpdatedBy,
				CreatedAt:     parentEvent.CreatedAt,
				UpdatedAt:     parentEvent.UpdatedAt,
				IsRecurring:   false,
				ParentEventID: &parentEvent.ID,
			}

			instances = append(instances, instance)
		}

		// Calculate next occurrence
		currentStart = calculateNextOccurrence(currentStart, pattern, interval)
	}

	return instances, nil
}

// calculateNextOccurrence calculates the next occurrence based on pattern and interval
func calculateNextOccurrence(current time.Time, pattern string, interval int) time.Time {
	switch pattern {
	case "daily":
		return current.AddDate(0, 0, interval)
	case "weekly":
		return current.AddDate(0, 0, 7*interval)
	case "monthly":
		return current.AddDate(0, interval, 0)
	default:
		return current
	}
}

// getTeamOverviewFunction returns team overview with stats and user role
// GET /api/v2/Teams('{teamId}')/GetOverview()
func (s *Service) getTeamOverviewFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	team := ctx.(*models.Team)

	// Get user ID from request context
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: missing user id")
	}

	// Get user from database
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify user is a member of the club
	club, err := models.GetClubByID(team.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !club.IsMember(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this club")
	}

	// Get team stats
	stats, err := team.GetTeamStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get team stats: %w", err)
	}

	// Get user's role in the team
	userRole := ""
	if team.IsMember(user) {
		userRole, _ = team.GetUserRole(user)
	}

	return TeamOverviewResponse{
		Team:     *team,
		Stats:    stats,
		UserRole: userRole,
		IsAdmin:  team.IsAdmin(user),
	}, nil
}

// getRSVPCountsFunction returns RSVP counts for an event without fetching all RSVPs
// GET /api/v2/Events('{eventId}')/GetRSVPCounts()
// Returns: {"Yes": 10, "No": 3, "Maybe": 5}
func (s *Service) getRSVPCountsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	event := ctx.(*models.Event)

	// Get user ID from request context
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: missing user id")
	}

	// Get user from database
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify user has access to this event by checking club membership
	club, err := models.GetClubByID(event.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !club.IsMember(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this club")
	}

	// Count RSVPs by response type using SQL aggregation
	type RSVPCount struct {
		Response string
		Count    int64
	}

	var counts []RSVPCount
	if err := s.db.Table("event_rsvps").
		Select("response, COUNT(*) as count").
		Where("event_id = ?", event.ID).
		Group("response").
		Scan(&counts).Error; err != nil {
		return nil, fmt.Errorf("failed to count RSVPs: %w", err)
	}

	// Build response map with PascalCase keys
	result := make(map[string]int64)
	for _, count := range counts {
		// Capitalize first letter to match OData v2 convention
		response := strings.ToUpper(count.Response[:1]) + strings.ToLower(count.Response[1:])
		result[response] = count.Count
	}

	// Ensure all response types are present (even if zero)
	if _, ok := result["Yes"]; !ok {
		result["Yes"] = 0
	}
	if _, ok := result["No"]; !ok {
		result["No"] = 0
	}
	if _, ok := result["Maybe"]; !ok {
		result["Maybe"] = 0
	}

	return result, nil
}
