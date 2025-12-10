package odata

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/models"
	odata "github.com/nlstn/go-odata"
)

// registerTimelineHandlers sets up the virtual entity handlers for Timeline
func (s *Service) registerTimelineHandlers() error {
	return s.Service.SetEntityOverwrite("TimelineItems", &odata.EntityOverwrite{
		GetCollection: s.getTimelineCollection,
		GetEntity:     s.getTimelineEntity,
		// Create, Update, Delete are not implemented - Timeline is read-only
	})
}

// getTimelineCollection retrieves all timeline items (activities, events, news) for the user
func (s *Service) getTimelineCollection(ctx *odata.OverwriteContext) (*odata.CollectionResult, error) {
	// Get user ID from request context
	userID, ok := ctx.Request.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	// Get user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get all clubs the user is a member of
	clubs, err := models.GetAllClubs()
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	var userClubIDs []string
	clubNameMap := make(map[string]string)
	for _, club := range clubs {
		if club.IsMember(user) {
			userClubIDs = append(userClubIDs, club.ID)
			clubNameMap[club.ID] = club.Name
		}
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

	// Apply query options if needed
	// Note: For a full implementation, we would need to apply $filter, $top, $skip here
	// For now, we'll return all items and let the OData library handle basic filtering

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

	// Fetch all items and find the matching one
	// This is not the most efficient approach, but works for virtual entities
	result, err := s.getTimelineCollection(ctx)
	if err != nil {
		return nil, err
	}

	items, ok := result.Items.([]models.TimelineItem)
	if !ok {
		return nil, fmt.Errorf("unexpected timeline items type")
	}

	for _, item := range items {
		if item.ID == timelineID {
			return item, nil
		}
	}

	return nil, fmt.Errorf("timeline item not found")
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

	var items []models.TimelineItem
	for _, event := range events {
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
