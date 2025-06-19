package handlers

import (
	"encoding/json"
	"net/http"

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

// NewsWithClub represents news with club information
type NewsWithClub struct {
	models.News
	ClubName string `json:"club_name"`
	ClubID   string `json:"club_id"`
}

// EventWithClub represents event with club information and RSVP status
type EventWithClub struct {
	models.Event
	ClubName string             `json:"club_name"`
	ClubID   string             `json:"club_id"`
	UserRSVP *models.EventRSVP  `json:"user_rsvp,omitempty"`
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