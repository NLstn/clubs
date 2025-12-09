package odata

import (
	"encoding/json"
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

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetUpcomingEvents",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]models.Event{}),
		Handler:    s.getUpcomingEventsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetUpcomingEvents function for Club: %w", err)
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

	// Unbound functions
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetDashboardNews",
		IsBound:    false,
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]NewsWithClub{}),
		Handler:    s.getDashboardNewsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetDashboardNews function: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetDashboardEvents",
		IsBound:    false,
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]EventWithClub{}),
		Handler:    s.getDashboardEventsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetDashboardEvents function: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetDashboardActivities",
		IsBound:    false,
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]ActivityItem{}),
		Handler:    s.getDashboardActivitiesFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetDashboardActivities function: %w", err)
	}

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

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetEvents",
		IsBound:    true,
		EntitySet:  "Teams",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]models.Event{}),
		Handler:    s.getTeamEventsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetEvents function for Team: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetUpcomingEvents",
		IsBound:    true,
		EntitySet:  "Teams",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]models.Event{}),
		Handler:    s.getTeamUpcomingEventsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetUpcomingEvents function for Team: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetFines",
		IsBound:    true,
		EntitySet:  "Teams",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]models.Fine{}),
		Handler:    s.getTeamFinesFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetFines function for Team: %w", err)
	}

	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetMembers",
		IsBound:    true,
		EntitySet:  "Teams",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]map[string]interface{}{}),
		Handler:    s.getTeamMembersFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetMembers function for Team: %w", err)
	}

	// More bound functions for Event entity
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetRSVPs",
		IsBound:    true,
		EntitySet:  "Events",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf(EventRSVPResponse{}),
		Handler:    s.getEventRSVPsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetRSVPs function for Event: %w", err)
	}

	// More bound functions for Club entity
	if err := s.Service.RegisterFunction(odata.FunctionDefinition{
		Name:       "GetMyTeams",
		IsBound:    true,
		EntitySet:  "Clubs",
		Parameters: []odata.ParameterDefinition{},
		ReturnType: reflect.TypeOf([]models.Team{}),
		Handler:    s.getMyTeamsFunction,
	}); err != nil {
		return fmt.Errorf("failed to register GetMyTeams function for Club: %w", err)
	}

	return nil
}

// Helper types for dashboard functions
type NewsWithClub struct {
	models.News
	ClubName string `json:"clubName"`
	ClubID   string `json:"clubId"`
}

type EventWithClub struct {
	models.Event
	ClubName string            `json:"clubName"`
	ClubID   string            `json:"clubId"`
	UserRSVP *models.EventRSVP `json:"userRsvp,omitempty"`
}

type ActivityItem struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content,omitempty"`
	ClubName  string                 `json:"clubName"`
	ClubID    string                 `json:"clubId"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
	Actor     string                 `json:"actor,omitempty"`
	ActorName string                 `json:"actorName,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type SearchResult struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ClubID      string `json:"clubId,omitempty"`
	ClubName    string `json:"clubName,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
}

type SearchResponse struct {
	Clubs  []SearchResult `json:"clubs"`
	Events []SearchResult `json:"events"`
}

type TeamOverviewResponse struct {
	Team     models.Team            `json:"team"`
	Stats    map[string]interface{} `json:"stats"`
	UserRole string                 `json:"userRole"`
	IsAdmin  bool                   `json:"isAdmin"`
}

type EventRSVPResponse struct {
	Counts map[string]int          `json:"counts"`
	RSVPs  []models.EventRSVP      `json:"rsvps"`
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

	return map[string]bool{"isAdmin": isAdmin}, nil
}

// getOwnerCountFunction returns the number of owners in the club
// GET /api/v2/Clubs('{clubId}')/GetOwnerCount()
func (s *Service) getOwnerCountFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	club := ctx.(*models.Club)

	count, err := club.CountOwners()
	if err != nil {
		return nil, fmt.Errorf("failed to count owners: %w", err)
	}

	return map[string]int64{"ownerCount": count}, nil
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

	return map[string]string{"inviteLink": inviteLink}, nil
}

// getUpcomingEventsFunction returns upcoming events for the club
// GET /api/v2/Clubs('{clubId}')/GetUpcomingEvents()
func (s *Service) getUpcomingEventsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	club := ctx.(*models.Club)

	events, err := club.GetUpcomingEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	return events, nil
}

// getDashboardNewsFunction returns news from all clubs the user is a member of
// GET /api/v2/GetDashboardNews()
func (s *Service) getDashboardNewsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get all clubs
	clubs, err := models.GetAllClubs()
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	var allNews []NewsWithClub
	for _, club := range clubs {
		if club.IsMember(user) {
			clubNews, err := club.GetNews()
			if err != nil {
				continue // Skip clubs where we can't fetch news
			}

			for _, news := range clubNews {
				newsWithClub := NewsWithClub{
					News:     news,
					ClubName: club.Name,
					ClubID:   club.ID,
				}
				allNews = append(allNews, newsWithClub)
			}
		}
	}

	return allNews, nil
}

// getDashboardEventsFunction returns upcoming events from all clubs the user is a member of
// GET /api/v2/GetDashboardEvents()
func (s *Service) getDashboardEventsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get all clubs
	clubs, err := models.GetAllClubs()
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	var allEvents []EventWithClub
	for _, club := range clubs {
		if club.IsMember(user) {
			clubEvents, err := club.GetUpcomingEvents()
			if err != nil {
				continue // Skip clubs where we can't fetch events
			}

			for _, event := range clubEvents {
				eventWithClub := EventWithClub{
					Event:    event,
					ClubName: club.Name,
					ClubID:   club.ID,
				}

				// Add user's RSVP status if available
				userRSVP, err := user.GetUserRSVP(event.ID)
				if err == nil {
					eventWithClub.UserRSVP = userRSVP
				}

				allEvents = append(allEvents, eventWithClub)
			}
		}
	}

	return allEvents, nil
}

// getDashboardActivitiesFunction returns recent activities from all clubs the user is a member of
// GET /api/v2/GetDashboardActivities()
func (s *Service) getDashboardActivitiesFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	// Get user ID from request context
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get all clubs
	clubs, err := models.GetAllClubs()
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	var userClubIDs []string
	clubNameMap := make(map[string]string)

	// Collect club IDs for clubs the user is a member of
	for _, club := range clubs {
		if club.IsMember(user) {
			userClubIDs = append(userClubIDs, club.ID)
			clubNameMap[club.ID] = club.Name
		}
	}

	var activities []ActivityItem

	// Get activities from the activity store
	if len(userClubIDs) > 0 {
		storedActivities, err := models.GetRecentActivities(userClubIDs, 30, 50)
		if err == nil {
			for _, activity := range storedActivities {
				// Parse metadata if it exists
				var metadata map[string]interface{}
				if activity.Metadata != "" {
					json.Unmarshal([]byte(activity.Metadata), &metadata)
				}

				// Determine actor
				var createdBy string
				if (activity.Type == "member_promoted" || activity.Type == "member_demoted" || activity.Type == "role_changed") && activity.ActorID != nil {
					createdBy = *activity.ActorID
				} else {
					createdBy = activity.UserID
				}

				activityItem := ActivityItem{
					ID:        activity.ID,
					Type:      activity.Type,
					Title:     activity.Title,
					Content:   activity.Content,
					ClubName:  clubNameMap[activity.ClubID],
					ClubID:    activity.ClubID,
					CreatedAt: activity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: activity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
					Actor:     createdBy,
					Metadata:  metadata,
				}

				activities = append(activities, activityItem)
			}
		}
	}

	if activities == nil {
		activities = []ActivityItem{}
	}

	return activities, nil
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

// getTeamEventsFunction returns all events for the team
// GET /api/v2/Teams('{teamId}')/GetEvents()
func (s *Service) getTeamEventsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
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

	// Verify user has access (team member or club admin)
	club, err := models.GetClubByID(team.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !team.IsMember(user) && !club.IsAdmin(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this team")
	}

	events, err := team.GetEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

// getTeamUpcomingEventsFunction returns upcoming events for the team
// GET /api/v2/Teams('{teamId}')/GetUpcomingEvents()
func (s *Service) getTeamUpcomingEventsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
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

	// Verify user has access (team member or club admin)
	club, err := models.GetClubByID(team.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !team.IsMember(user) && !club.IsAdmin(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this team")
	}

	events, err := team.GetUpcomingEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	return events, nil
}

// getTeamFinesFunction returns all fines for the team
// GET /api/v2/Teams('{teamId}')/GetFines()
func (s *Service) getTeamFinesFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
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

	// Verify user has access (team member or club admin)
	club, err := models.GetClubByID(team.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !team.IsMember(user) && !club.IsAdmin(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this team")
	}

	// Get fines for this team
	var fines []models.Fine
	if err := s.db.Where("team_id = ?", team.ID).Preload("User").Find(&fines).Error; err != nil {
		return nil, fmt.Errorf("failed to get fines: %w", err)
	}

	return fines, nil
}

// getTeamMembersFunction returns all team members with user details
// GET /api/v2/Teams('{teamId}')/GetMembers()
func (s *Service) getTeamMembersFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
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

	// Get team members with user details using the existing method
	members, err := team.GetTeamMembersWithUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	return members, nil
}

// getEventRSVPsFunction returns all RSVPs for an event with counts
// GET /api/v2/Events('{eventId}')/GetRSVPs()
func (s *Service) getEventRSVPsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
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

	// Verify user is an admin of the club
	club, err := models.GetClubByID(event.ClubID)
	if err != nil {
		return nil, fmt.Errorf("failed to find club: %w", err)
	}

	if !club.IsOwner(user) && !club.IsAdmin(user) {
		return nil, fmt.Errorf("forbidden: only club admins can view RSVPs")
	}

	// Get RSVP counts
	counts, err := models.GetEventRSVPCounts(event.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get RSVP counts: %w", err)
	}

	// Get RSVPs with user details
	rsvps, err := models.GetEventRSVPs(event.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get RSVPs: %w", err)
	}

	return EventRSVPResponse{
		Counts: counts,
		RSVPs:  rsvps,
	}, nil
}

// getMyTeamsFunction returns teams the current user is a member of
// GET /api/v2/Clubs('{clubId}')/GetMyTeams()
func (s *Service) getMyTeamsFunction(w http.ResponseWriter, r *http.Request, ctx interface{}, params map[string]interface{}) (interface{}, error) {
	club := ctx.(*models.Club)

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
	if !club.IsMember(user) {
		return nil, fmt.Errorf("forbidden: user is not a member of this club")
	}

	// Get teams where user is a member by joining through team_members table
	var teams []models.Team
	if err := s.db.
		Table("teams").
		Joins("INNER JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ? AND teams.club_id = ?", userID, club.ID).
		Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("failed to get user's teams: %w", err)
	}

	return teams, nil
}
