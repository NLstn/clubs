package main

import (
	"log"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/NLstn/clubs/odata"
	frontend "github.com/NLstn/clubs/tools"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// ignore error if .env file is not found

	err := database.Init()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	// Run database migrations
	// NOTE: AutoMigrate must be called here (not in database package) to avoid circular dependency.
	// The models package imports database for DB connection, so database cannot import models.
	// This is the standard pattern for GORM applications.
	err = database.Db.AutoMigrate(&models.Club{},
		&models.Member{},
		&models.Team{},
		&models.TeamMember{},
		&models.MagicLink{},
		&models.User{},
		&models.JoinRequest{},
		&models.Invite{},
		&models.RefreshToken{},
		&models.Fine{},
		&models.FineTemplate{},
		&models.Shift{},
		&models.ShiftMember{},
		&models.Event{},
		&models.EventRSVP{},
		&models.News{},
		&models.ClubSettings{},
		&models.Notification{},
		&models.UserNotificationPreferences{},
		&models.UserPrivacySettings{},
		&models.MemberPrivacySettings{},
		&models.Activity{},
		&models.APIKey{},
	)
	if err != nil {
		log.Fatal("Could not migrate database:", err)
	}

	err = azure.Init()
	if err != nil {
		log.Fatal("Could not initialize Azure SDK:", err)
	}

	err = auth.Init()
	if err != nil {
		log.Fatal("Could not initialize auth:", err)
	}

	err = auth.InitKeycloak()
	if err != nil {
		log.Printf("Warning: Could not initialize Keycloak: %v", err)
	}

	err = frontend.Init()
	if err != nil {
		log.Fatal("Could not initialize frontend:", err)
	}

	mux := http.NewServeMux()

	// Add health check at root level for container monitoring
	mux.HandleFunc("/health", handlers.HealthCheck)

	mux.Handle("/api/v1/", handlers.Handler_v1())

	// Mount OData v2 API (Phase 2: With authentication and authorization)
	odataService, err := odata.NewService(database.Db)
	if err != nil {
		log.Fatal("Could not initialize OData service:", err)
	}

	// Get JWT secret for OData authentication middleware
	jwtSecret := []byte(auth.GetJWTSecret())

	// Create a submux for /api/v2/ to handle both OData and custom routes
	odataV2Mux := http.NewServeMux()

	// Register custom handlers (e.g., file uploads) that don't fit standard OData patterns
	odataService.RegisterCustomHandlers(odataV2Mux)

	// Register the OData service as the default handler
	odataV2Mux.Handle("/", odataService)

	// Wrap OData v2 service with authentication middleware
	// This enforces JWT token validation on all /api/v2/ endpoints
	// except for metadata and service document endpoints
	odataWithAuth := http.StripPrefix("/api/v2", odata.AuthMiddleware(jwtSecret)(odataV2Mux))
	mux.Handle("/api/v2/", odataWithAuth)

	handler := handlers.CorsMiddleware(mux)
	handlerWithLogging := handlers.LoggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
