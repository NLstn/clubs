package odata

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

// entitySetToFeature maps entity set names to their corresponding feature names
var entitySetToFeature = map[string]string{
	"Fines":         "fines",
	"FineTemplates": "fines",
	"Shifts":        "shifts",
	"ShiftMembers":  "shifts",
	"Teams":         "teams",
	"TeamMembers":   "teams",
}

// entitySetPattern matches OData entity set paths with optional key predicates
// Examples: /Fines, /Fines('uuid'), /Fines('uuid')/User
var entitySetPattern = regexp.MustCompile(`^/?([A-Z][a-zA-Z]*)(?:\([^)]+\))?`)

// FeatureCheckMiddleware checks if the requested feature is enabled for the club
// Returns HTTP 400 if a feature is disabled
func FeatureCheckMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip metadata and service document
			path := r.URL.Path
			if path == "" || path == "/" || path == "$metadata" || path == "/$metadata" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract entity set name from path
			matches := entitySetPattern.FindStringSubmatch(path)
			if len(matches) < 2 {
				// Not a standard entity set request, let it pass
				next.ServeHTTP(w, r)
				return
			}

			entitySet := matches[1]
			featureName, requiresCheck := entitySetToFeature[entitySet]
			if !requiresCheck {
				// This entity doesn't require feature check
				next.ServeHTTP(w, r)
				return
			}

			// Get user ID from context (set by auth middleware)
			userID, ok := r.Context().Value(auth.UserIDKey).(string)
			if !ok || userID == "" {
				// No user context, let auth middleware handle it
				next.ServeHTTP(w, r)
				return
			}

			// Get club ID from the entity or query
			clubID, err := getClubIDFromRequest(r, entitySet)
			if err != nil {
				// Can't determine club ID, let OData handle the request
				// The OData hooks will handle authorization
				log.Printf("WARNING: Could not determine club ID for feature check: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			// Check if the feature is enabled
			if err := models.CheckFeatureEnabled(clubID, featureName); err != nil {
				var featureErr *models.FeatureDisabledError
				if errors.As(err, &featureErr) {
					// Feature is disabled, return 400
					writeODataError(w, http.StatusBadRequest, "FeatureDisabled", featureErr.Message)
					return
				}
				// Other error, let it continue (OData will handle it)
				log.Printf("ERROR: Failed to check feature status: %v", err)
			}

			// Feature is enabled or check was skipped, continue
			next.ServeHTTP(w, r)
		})
	}
}

// getClubIDFromRequest attempts to extract the club ID from the request
func getClubIDFromRequest(r *http.Request, entitySet string) (string, error) {
	// For GET requests to collections, check if there's a club filter in query
	// For POST requests, we'd need to parse the body, but that's complex and not reliable
	// For GET requests with entity key, we need to look up the entity in the database
	
	// Strategy: For most entities, we can query the members table to get all clubs the user is in
	// Then check settings for all those clubs. However, this is inefficient.
	// Better approach: For now, use a simplified check based on entity relationships

	path := r.URL.Path
	
	// Check if path contains an entity key (e.g., /Fines('uuid'))
	if strings.Contains(path, "(") && strings.Contains(path, ")") {
		// Extract the key value
		start := strings.Index(path, "('") + 2
		end := strings.Index(path[start:], "')")
		if start > 2 && end > 0 {
			entityID := path[start : start+end]
			return getClubIDFromEntity(entitySet, entityID)
		}
	}

	// For collection requests (no key), we can't easily determine a single club ID
	// The OData hooks will filter by club membership anyway
	// So we'll skip the check here and let the hooks handle it
	// This means we can't return 400 for collection requests, only for specific entity requests
	
	// However, for creates (POST), the ClubID should be in the request body
	// Let's handle that case
	if r.Method == "POST" {
		// We could parse the body here, but it's complex and might interfere with OData processing
		// Skip for now - POST requests will be checked at the OData hook level
		return "", fmt.Errorf("cannot determine club ID for POST requests in middleware")
	}

	// For collection queries, we can't determine a single club ID
	// Return error to skip the check
	return "", fmt.Errorf("cannot determine club ID for collection queries")
}

// getClubIDFromEntity looks up the club ID for a specific entity
func getClubIDFromEntity(entitySet, entityID string) (string, error) {
	var clubID string
	var err error

	switch entitySet {
	case "Fines":
		var fine models.Fine
		err = database.Db.Select("club_id").Where("id = ?", entityID).First(&fine).Error
		clubID = fine.ClubID
	case "FineTemplates":
		var template models.FineTemplate
		err = database.Db.Select("club_id").Where("id = ?", entityID).First(&template).Error
		clubID = template.ClubID
	case "Shifts":
		var shift models.Shift
		err = database.Db.Select("club_id").Where("id = ?", entityID).First(&shift).Error
		clubID = shift.ClubID
	case "ShiftMembers":
		// ShiftMember doesn't have ClubID directly, need to join through Shift
		var shiftMember models.ShiftMember
		err = database.Db.Preload("Shift").Where("id = ?", entityID).First(&shiftMember).Error
		if err == nil && shiftMember.Shift != nil {
			clubID = shiftMember.Shift.ClubID
		}
	case "Teams":
		var team models.Team
		err = database.Db.Select("club_id").Where("id = ?", entityID).First(&team).Error
		clubID = team.ClubID
	case "TeamMembers":
		// TeamMember doesn't have ClubID directly, need to join through Team
		var teamMember models.TeamMember
		err = database.Db.Preload("Team").Where("id = ?", entityID).First(&teamMember).Error
		if err == nil && teamMember.Team != nil {
			clubID = teamMember.Team.ClubID
		}
	default:
		return "", fmt.Errorf("unknown entity set: %s", entitySet)
	}

	if err == gorm.ErrRecordNotFound {
		return "", fmt.Errorf("entity not found")
	}
	if err != nil {
		return "", err
	}
	if clubID == "" {
		return "", fmt.Errorf("club ID not found for entity")
	}

	return clubID, nil
}

// writeODataError writes an OData v4 compliant error response
func writeODataError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("OData-Version", "4.01")
	w.WriteHeader(statusCode)
	
	errorJSON := fmt.Sprintf(`{"error":{"code":"%s","message":"%s"}}`, code, message)
	w.Write([]byte(errorJSON))
}
