package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/database"
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
}

type NewsWithClub struct {
	models.News
	Club models.Club `json:"club"`
}

type EventWithClub struct {
	models.Event
	Club models.Club `json:"club"`
}

// GET /api/v1/dashboard/news
func handleGetDashboardNews(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	// Get all clubs where the user is a member
	var clubs []models.Club
	err := database.Db.Joins("JOIN members ON clubs.id = members.club_id").
		Where("members.user_id = ?", user.ID).
		Find(&clubs).Error
	if err != nil {
		http.Error(w, "Failed to get user's clubs", http.StatusInternalServerError)
		return
	}

	var allNews []NewsWithClub
	
	// Get news from each club
	for _, club := range clubs {
		var news []models.News
		err := database.Db.Where("club_id = ?", club.ID).
			Order("created_at DESC").
			Limit(10). // Limit per club to avoid too much data
			Find(&news).Error
		if err != nil {
			continue // Skip this club on error, don't fail entire request
		}

		// Add club info to each news item
		for _, newsItem := range news {
			allNews = append(allNews, NewsWithClub{
				News: newsItem,
				Club: club,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allNews)
}

// GET /api/v1/dashboard/events
func handleGetDashboardEvents(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	// Get all clubs where the user is a member
	var clubs []models.Club
	err := database.Db.Joins("JOIN members ON clubs.id = members.club_id").
		Where("members.user_id = ?", user.ID).
		Find(&clubs).Error
	if err != nil {
		http.Error(w, "Failed to get user's clubs", http.StatusInternalServerError)
		return
	}

	var allEvents []EventWithClub
	
	// Get upcoming events from each club
	for _, club := range clubs {
		events, err := club.GetUpcomingEvents()
		if err != nil {
			continue // Skip this club on error, don't fail entire request
		}

		// Add club info to each event
		for _, event := range events {
			allEvents = append(allEvents, EventWithClub{
				Event: event,
				Club:  club,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allEvents)
}