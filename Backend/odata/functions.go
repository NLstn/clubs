package odata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

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
		Name:       "SearchGlobal",
		IsBound:    false,
		Parameters: []odata.ParameterDefinition{
			{Name: "query", Type: reflect.TypeOf(""), Required: true},
		},
		ReturnType: reflect.TypeOf(SearchResponse{}),
		Handler:    s.searchGlobalFunction,
	}); err != nil {
		return fmt.Errorf("failed to register SearchGlobal function: %w", err)
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
