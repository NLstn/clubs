package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NLstn/clubs/models"
)

func registerSearchRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/search", RateLimitMiddleware(apiLimiter)(withAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGlobalSearch(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

type SearchResult struct {
	Type        string `json:"type"` // "club" or "event"
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ClubID      string `json:"club_id,omitempty"`    // For events
	ClubName    string `json:"club_name,omitempty"`  // For events
	StartTime   string `json:"start_time,omitempty"` // For events
	EndTime     string `json:"end_time,omitempty"`   // For events
}

type SearchResponse struct {
	Clubs  []SearchResult `json:"clubs"`
	Events []SearchResult `json:"events"`
}

// endpoint: GET /api/v1/search?q=query
func handleGlobalSearch(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		response := SearchResponse{
			Clubs:  []SearchResult{},
			Events: []SearchResult{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Search clubs
	clubResults, err := searchClubs(user, query)
	if err != nil {
		http.Error(w, "Failed to search clubs", http.StatusInternalServerError)
		return
	}

	// Search events
	eventResults, err := searchEvents(user, query)
	if err != nil {
		http.Error(w, "Failed to search events", http.StatusInternalServerError)
		return
	}

	response := SearchResponse{
		Clubs:  clubResults,
		Events: eventResults,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func searchClubs(user models.User, query string) ([]SearchResult, error) {
	// Get all clubs that the user is a member of
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
		if strings.Contains(strings.ToLower(club.Name), queryLower) ||
			strings.Contains(strings.ToLower(club.Description), queryLower) {

			results = append(results, SearchResult{
				Type:        "club",
				ID:          club.ID,
				Name:        club.Name,
				Description: club.Description,
			})
		}
	}

	return results, nil
}

func searchEvents(user models.User, query string) ([]SearchResult, error) {
	// Get all events from clubs the user is a member of
	events, err := models.SearchEventsForUser(user.ID, query)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, event := range events {
		results = append(results, SearchResult{
			Type:      "event",
			ID:        event.ID,
			Name:      event.Name,
			ClubID:    event.ClubID,
			ClubName:  event.ClubName,
			StartTime: event.StartTime.Format("2006-01-02T15:04:05Z07:00"),
			EndTime:   event.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return results, nil
}
