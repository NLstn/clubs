package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NLstn/clubs/database"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// HealthCheck handles health check requests for container monitoring
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status: "healthy",
		Services: map[string]string{
			"api": "healthy",
		},
	}

	// Check database connectivity
	if database.Db != nil {
		sqlDB, err := database.Db.DB()
		if err != nil || sqlDB.Ping() != nil {
			response.Status = "unhealthy"
			response.Services["database"] = "unhealthy"
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			response.Services["database"] = "healthy"
		}
	} else {
		response.Status = "unhealthy"
		response.Services["database"] = "unavailable"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	// Set status code to 200 if healthy, already set to 503 if unhealthy
	if response.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}