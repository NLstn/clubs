package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/NLstn/clubs/models"
)

func registerDashboardRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/dashboard/news", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetDashboardNews(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/dashboard/events", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetDashboardEvents(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/v1/dashboard/activities", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetDashboardActivities(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

// NewsWithClub represents news with club information
type NewsWithClub struct {
	models.News
	ClubName string `json:"club_name"`
	ClubID   string `json:"club_id"`
}

// EventWithClub represents event with club information and RSVP status
type EventWithClub struct {
	models.Event
	ClubName string            `json:"club_name"`
	ClubID   string            `json:"club_id"`
	UserRSVP *models.EventRSVP `json:"user_rsvp,omitempty"`
}

// ActivityItem represents a unified activity feed item
type ActivityItem struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // "news", "event", "role_changed", "member_promoted", "member_demoted"
	Title     string                 `json:"title"`
	Content   string                 `json:"content,omitempty"`
	ClubName  string                 `json:"club_name"`
	ClubID    string                 `json:"club_id"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Actor     string                 `json:"actor,omitempty"`      // User ID who created/initiated the activity
	ActorName string                 `json:"actor_name,omitempty"` // Name of the user who created/initiated the activity
	Metadata  map[string]interface{} `json:"metadata,omitempty"`   // For extensibility
}

// GET /api/v1/dashboard/news
func handleGetDashboardNews(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	clubs, err := models.GetAllClubs()
	if err != nil {
		http.Error(w, "Failed to get clubs", http.StatusInternalServerError)
		return
	}

	var allNews []NewsWithClub
	for _, club := range clubs {
		if club.IsMember(user) {
			clubNews, err := club.GetNews()
			if err != nil {
				continue // Skip clubs where we can't fetch news, don't fail entirely
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allNews)
}

// GET /api/v1/dashboard/events
func handleGetDashboardEvents(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	clubs, err := models.GetAllClubs()
	if err != nil {
		http.Error(w, "Failed to get clubs", http.StatusInternalServerError)
		return
	}

	var allEvents []EventWithClub
	for _, club := range clubs {
		if club.IsMember(user) {
			clubEvents, err := club.GetUpcomingEvents()
			if err != nil {
				continue // Skip clubs where we can't fetch events, don't fail entirely
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allEvents)
}

// GET /api/v1/dashboard/activities
func handleGetDashboardActivities(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	clubs, err := models.GetAllClubs()
	if err != nil {
		http.Error(w, "Failed to get clubs", http.StatusInternalServerError)
		return
	}

	var activities []ActivityItem
	var creatorIDs []string
	var userClubIDs []string

	// Collect club IDs for clubs the user is a member of
	for _, club := range clubs {
		if club.IsMember(user) {
			userClubIDs = append(userClubIDs, club.ID)
		}
	}

	// Get activities from the activity store (role changes, promotions, etc.)
	if len(userClubIDs) > 0 {
		storedActivities, err := models.GetRecentActivities(userClubIDs, 30, 50) // Last 30 days, max 50 items
		if err == nil {
			for _, activity := range storedActivities {
				// Find the club name
				var clubName string
				for _, club := range clubs {
					if club.ID == activity.ClubID {
						clubName = club.Name
						break
					}
				}

				// Parse metadata if it exists
				var metadata map[string]interface{}
				if activity.Metadata != "" {
					json.Unmarshal([]byte(activity.Metadata), &metadata)
				}

				// For role change activities, use the actor (who made the change) as created_by
				// For other activities, use the user who performed the action
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
					ClubName:  clubName,
					ClubID:    activity.ClubID,
					CreatedAt: activity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: activity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
					Actor:     createdBy,
					Metadata:  metadata,
				}

				// For role change activities, include the affected user ID
				if activity.Type == "member_promoted" || activity.Type == "member_demoted" || activity.Type == "role_changed" {
					activityItem.Metadata["affected_user_id"] = activity.UserID
				}

				activities = append(activities, activityItem)
				if createdBy != "" {
					creatorIDs = append(creatorIDs, createdBy)
				}
			}
		}
	}

	// Collect news as activities
	for _, club := range clubs {
		if club.IsMember(user) {
			clubNews, err := club.GetNews()
			if err != nil {
				continue // Skip clubs where we can't fetch news
			}

			for _, news := range clubNews {
				activity := ActivityItem{
					ID:        news.ID,
					Type:      "news",
					Title:     news.Title,
					Content:   news.Content,
					ClubName:  club.Name,
					ClubID:    club.ID,
					CreatedAt: news.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: news.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
					Actor:     news.CreatedBy,
				}
				activities = append(activities, activity)
				if news.CreatedBy != "" {
					creatorIDs = append(creatorIDs, news.CreatedBy)
				}
			}
		}
	}

	// Collect events as activities
	for _, club := range clubs {
		if club.IsMember(user) {
			clubEvents, err := club.GetUpcomingEvents()
			if err != nil {
				continue // Skip clubs where we can't fetch events
			}

			for _, event := range clubEvents {
				// Create event content description
				eventContent := "Event scheduled"
				if !event.StartTime.IsZero() && !event.EndTime.IsZero() {
					eventContent = "Event scheduled from " + event.StartTime.Format("2006-01-02 15:04") + " to " + event.EndTime.Format("2006-01-02 15:04")
				}

				// Add RSVP info to metadata if available
				metadata := make(map[string]interface{})
				userRSVP, err := user.GetUserRSVP(event.ID)
				if err == nil {
					metadata["user_rsvp"] = userRSVP
				}
				metadata["start_time"] = event.StartTime.Format("2006-01-02T15:04:05Z07:00")
				metadata["end_time"] = event.EndTime.Format("2006-01-02T15:04:05Z07:00")

				activity := ActivityItem{
					ID:        event.ID,
					Type:      "event",
					Title:     event.Name,
					Content:   eventContent,
					ClubName:  club.Name,
					ClubID:    club.ID,
					CreatedAt: event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: event.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
					Actor:     event.CreatedBy,
					Metadata:  metadata,
				}
				activities = append(activities, activity)
				if event.CreatedBy != "" {
					creatorIDs = append(creatorIDs, event.CreatedBy)
				}
			}
		}
	}

	// Fetch creator information for all activities
	if len(creatorIDs) > 0 {
		// Remove duplicates from creatorIDs
		uniqueCreatorIDs := make([]string, 0, len(creatorIDs))
		seen := make(map[string]bool)
		for _, id := range creatorIDs {
			if !seen[id] {
				uniqueCreatorIDs = append(uniqueCreatorIDs, id)
				seen[id] = true
			}
		}

		// Get user information for all creators
		creators, err := models.GetUsersByIDs(uniqueCreatorIDs)
		if err == nil {
			// Create a map for quick lookup
			creatorMap := make(map[string]models.User)
			for _, creator := range creators {
				creatorMap[creator.ID] = creator
			}

			// Add creator names to activities
			for i := range activities {
				if creator, exists := creatorMap[activities[i].Actor]; exists {
					activities[i].ActorName = creator.GetFullName()
					if activities[i].ActorName == "" {
						activities[i].ActorName = creator.Email // Fallback to email if name is empty
					}
				}
			}
		}
	}

	// Sort activities by creation date (most recent first)
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].CreatedAt > activities[j].CreatedAt
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}
