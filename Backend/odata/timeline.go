package odata

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/NLstn/civo/auth"
	"github.com/NLstn/civo/models"
	odata "github.com/nlstn/go-odata"
	"gorm.io/gorm"
)

// registerTimelineHandlers sets up the virtual entity handlers for Timeline
func (s *Service) registerTimelineHandlers() error {
	return s.Service.SetEntityOverwrite("TimelineItems", &odata.EntityOverwrite{
		GetCollection: s.getTimelineCollection,
		GetEntity:     s.getTimelineEntity,
		// Create, Update, Delete are not implemented - Timeline is read-only
	})
}

// getUserClubs fetches all clubs the user is a member of in a single query
func (s *Service) getUserClubs(userID string) ([]string, map[string]string, error) {
	var members []models.Member
	if err := s.db.Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get user memberships: %w", err)
	}

	if len(members) == 0 {
		return []string{}, make(map[string]string), nil
	}

	// Extract club IDs
	clubIDs := make([]string, len(members))
	for i, member := range members {
		clubIDs[i] = member.ClubID
	}

	// Fetch clubs in a single query
	var clubs []models.Club
	if err := s.db.Where("id IN ? AND deleted = false", clubIDs).Find(&clubs).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	// Build club name map
	clubNameMap := make(map[string]string)
	validClubIDs := make([]string, 0, len(clubs))
	for _, club := range clubs {
		clubNameMap[club.ID] = club.Name
		validClubIDs = append(validClubIDs, club.ID)
	}

	return validClubIDs, clubNameMap, nil
}

// getTimelineCollection retrieves all timeline items (activities, events, news) for the user
func (s *Service) getTimelineCollection(ctx *odata.OverwriteContext) (*odata.CollectionResult, error) {
	// Get user ID from request context
	userID, ok := ctx.Request.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get all clubs the user is a member of (single query)
	userClubIDs, clubNameMap, err := s.getUserClubs(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user clubs: %w", err)
	}

	if len(userClubIDs) == 0 {
		// User is not a member of any clubs
		return &odata.CollectionResult{Items: []models.TimelineItem{}}, nil
	}

	var timelineItems []models.TimelineItem

	// Fetch activities
	activities, err := s.fetchActivities(userClubIDs, clubNameMap)
	if err != nil {
		s.logger.Error("Failed to fetch activities for timeline", "error", err)
		// Continue even if activities fail
	} else {
		timelineItems = append(timelineItems, activities...)
	}

	// Fetch events
	events, err := s.fetchEvents(userClubIDs, clubNameMap, userID)
	if err != nil {
		s.logger.Error("Failed to fetch events for timeline", "error", err)
		// Continue even if events fail
	} else {
		timelineItems = append(timelineItems, events...)
	}

	// Fetch news
	news, err := s.fetchNews(userClubIDs, clubNameMap)
	if err != nil {
		s.logger.Error("Failed to fetch news for timeline", "error", err)
		// Continue even if news fail
	} else {
		timelineItems = append(timelineItems, news...)
	}

	// Sort by timestamp (most recent first)
	sort.Slice(timelineItems, func(i, j int) bool {
		return timelineItems[i].Timestamp.After(timelineItems[j].Timestamp)
	})

	// Note: Virtual entity limitations
	// The OData library handles in-memory virtual entities but does not support
	// server-side filtering ($filter), pagination ($top, $skip), or sorting ($orderby).
	// These operations are limited by the in-memory data model.
	// For production use with large datasets, consider implementing manual filtering
	// and pagination logic or converting this to a database-backed entity.

	return &odata.CollectionResult{Items: timelineItems}, nil
}

// getTimelineEntity retrieves a single timeline item by ID
func (s *Service) getTimelineEntity(ctx *odata.OverwriteContext) (interface{}, error) {
	// Get user ID from request context
	userID, ok := ctx.Request.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Parse the timeline item ID to determine type and actual ID
	// Format: "activity-{id}", "event-{id}", or "news-{id}"
	timelineID := ctx.EntityKey

	// Parse ID to extract type and entity ID
	// Split on first hyphen only since UUIDs contain hyphens
	parts := strings.SplitN(timelineID, "-", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid timeline item ID format")
	}
	itemType := parts[0]
	itemID := parts[1]

	// Get user's clubs for authorization and club name mapping
	userClubIDs, clubNameMap, err := s.getUserClubs(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user clubs: %w", err)
	}

	// Convert to map for quick lookup
	clubIDSet := make(map[string]bool)
	for _, clubID := range userClubIDs {
		clubIDSet[clubID] = true
	}

	// Fetch specific item based on type
	switch itemType {
	case "activity":
		var activity models.Activity
		if err := s.db.Where("id = ?", itemID).First(&activity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("activity not found")
			}
			return nil, fmt.Errorf("failed to fetch activity: %w", err)
		}

		// Check authorization
		if !clubIDSet[activity.ClubID] {
			return nil, fmt.Errorf("access denied: user is not a member of this club")
		}

		// Parse metadata if it exists
		var metadata map[string]interface{}
		if activity.Metadata != "" {
			if err := json.Unmarshal([]byte(activity.Metadata), &metadata); err != nil {
				s.logger.Warn("Failed to parse activity metadata", "activityID", activity.ID, "error", err)
				metadata = make(map[string]interface{})
			}
		} else {
			metadata = make(map[string]interface{})
		}

		var actor *string
		var actorName *string
		if activity.ActorID != nil && *activity.ActorID != "" {
			actor = activity.ActorID
		}

		return models.TimelineItem{
			ID:        timelineID,
			ClubID:    activity.ClubID,
			ClubName:  clubNameMap[activity.ClubID],
			Type:      "activity",
			Title:     activity.Title,
			Content:   activity.Content,
			Timestamp: activity.CreatedAt,
			CreatedAt: activity.CreatedAt,
			UpdatedAt: activity.UpdatedAt,
			Actor:     actor,
			ActorName: actorName,
			Metadata:  metadata,
		}, nil

	case "event":
		var event models.Event
		if err := s.db.Where("id = ?", itemID).First(&event).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("event not found")
			}
			return nil, fmt.Errorf("failed to fetch event: %w", err)
		}

		// Check authorization
		if !clubIDSet[event.ClubID] {
			return nil, fmt.Errorf("access denied: user is not a member of this club")
		}

		// Get user's RSVP if available
		var userRSVP *models.EventRSVP
		var rsvp models.EventRSVP
		if err := s.db.Where("event_id = ? AND user_id = ?", event.ID, userID).First(&rsvp).Error; err == nil {
			userRSVP = &rsvp
		}

		metadata := make(map[string]interface{})
		if event.Description != nil {
			metadata["description"] = *event.Description
		}
		if event.Location != nil {
			metadata["location"] = *event.Location
		}

		return models.TimelineItem{
			ID:        timelineID,
			ClubID:    event.ClubID,
			ClubName:  clubNameMap[event.ClubID],
			Type:      "event",
			Title:     event.Name,
			Content:   "",
			Timestamp: event.StartTime,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
			StartTime: &event.StartTime,
			EndTime:   &event.EndTime,
			Location:  event.Location,
			UserRSVP:  userRSVP,
			Metadata:  metadata,
		}, nil

	case "news":
		var news models.News
		if err := s.db.Where("id = ?", itemID).First(&news).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("news not found")
			}
			return nil, fmt.Errorf("failed to fetch news: %w", err)
		}

		// Check authorization
		if !clubIDSet[news.ClubID] {
			return nil, fmt.Errorf("access denied: user is not a member of this club")
		}

		return models.TimelineItem{
			ID:        timelineID,
			ClubID:    news.ClubID,
			ClubName:  clubNameMap[news.ClubID],
			Type:      "news",
			Title:     news.Title,
			Content:   news.Content,
			Timestamp: news.CreatedAt,
			CreatedAt: news.CreatedAt,
			UpdatedAt: news.UpdatedAt,
			Metadata:  make(map[string]interface{}),
		}, nil

	default:
		return nil, fmt.Errorf("unknown timeline item type: %s", itemType)
	}
}

// fetchActivities fetches activities from the database and converts them to timeline items
func (s *Service) fetchActivities(clubIDs []string, clubNameMap map[string]string) ([]models.TimelineItem, error) {
	if len(clubIDs) == 0 {
		return []models.TimelineItem{}, nil
	}

	activities, err := models.GetRecentActivities(clubIDs, 30, 50)
	if err != nil {
		return nil, err
	}

	var items []models.TimelineItem
	for _, activity := range activities {
		// Parse metadata if it exists
		var metadata map[string]interface{}
		if activity.Metadata != "" {
			if err := json.Unmarshal([]byte(activity.Metadata), &metadata); err != nil {
				s.logger.Warn("Failed to parse activity metadata", "activityID", activity.ID, "error", err)
				metadata = make(map[string]interface{})
			}
		} else {
			metadata = make(map[string]interface{})
		}

		// Determine actor
		var actor *string
		var actorName *string
		if activity.ActorID != nil && *activity.ActorID != "" {
			actor = activity.ActorID
			// Optionally fetch actor name if needed
		}

		item := models.TimelineItem{
			ID:        fmt.Sprintf("activity-%s", activity.ID),
			ClubID:    activity.ClubID,
			ClubName:  clubNameMap[activity.ClubID],
			Type:      "activity",
			Title:     activity.Title,
			Content:   activity.Content,
			Timestamp: activity.CreatedAt,
			CreatedAt: activity.CreatedAt,
			UpdatedAt: activity.UpdatedAt,
			Actor:     actor,
			ActorName: actorName,
			Metadata:  metadata,
		}

		items = append(items, item)
	}

	return items, nil
}

// fetchEvents fetches upcoming events from the database and converts them to timeline items
func (s *Service) fetchEvents(clubIDs []string, clubNameMap map[string]string, userID string) ([]models.TimelineItem, error) {
	if len(clubIDs) == 0 {
		return []models.TimelineItem{}, nil
	}

	var events []models.Event
	now := time.Now()
	err := s.db.Where("club_id IN ? AND start_time >= ?", clubIDs, now).
		Order("start_time ASC").
		Limit(50).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	// Batch fetch RSVPs for these events to avoid N+1 query
	eventIDs := make([]string, len(events))
	for i, event := range events {
		eventIDs[i] = event.ID
	}

	var rsvps []models.EventRSVP
	if len(eventIDs) > 0 {
		if err := s.db.Where("event_id IN ? AND user_id = ?", eventIDs, userID).Find(&rsvps).Error; err != nil {
			s.logger.Warn("Failed to fetch RSVPs", "error", err)
		}
	}

	// Build a map from event_id to RSVP for quick lookup
	rsvpMap := make(map[string]*models.EventRSVP)
	for i := range rsvps {
		r := rsvps[i]
		rsvpMap[r.EventID] = &r
	}

	var items []models.TimelineItem
	for _, event := range events {
		// Lookup user's RSVP if available
		var userRSVP *models.EventRSVP
		if rsvp, ok := rsvpMap[event.ID]; ok {
			userRSVP = rsvp
		}

		metadata := make(map[string]interface{})
		if event.Description != nil {
			metadata["description"] = *event.Description
		}
		if event.Location != nil {
			metadata["location"] = *event.Location
		}

		item := models.TimelineItem{
			ID:        fmt.Sprintf("event-%s", event.ID),
			ClubID:    event.ClubID,
			ClubName:  clubNameMap[event.ClubID],
			Type:      "event",
			Title:     event.Name,
			Content:   "", // Events don't have content field
			Timestamp: event.StartTime,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
			StartTime: &event.StartTime,
			EndTime:   &event.EndTime,
			Location:  event.Location,
			UserRSVP:  userRSVP,
			Metadata:  metadata,
		}

		items = append(items, item)
	}

	return items, nil
}

// fetchNews fetches recent news from the database and converts them to timeline items
func (s *Service) fetchNews(clubIDs []string, clubNameMap map[string]string) ([]models.TimelineItem, error) {
	if len(clubIDs) == 0 {
		return []models.TimelineItem{}, nil
	}

	var newsList []models.News
	err := s.db.Where("club_id IN ?", clubIDs).
		Order("created_at DESC").
		Limit(50).
		Find(&newsList).Error

	if err != nil {
		return nil, err
	}

	var items []models.TimelineItem
	for _, news := range newsList {
		item := models.TimelineItem{
			ID:        fmt.Sprintf("news-%s", news.ID),
			ClubID:    news.ClubID,
			ClubName:  clubNameMap[news.ClubID],
			Type:      "news",
			Title:     news.Title,
			Content:   news.Content,
			Timestamp: news.CreatedAt,
			CreatedAt: news.CreatedAt,
			UpdatedAt: news.UpdatedAt,
			Metadata:  make(map[string]interface{}),
		}

		items = append(items, item)
	}

	return items, nil
}
