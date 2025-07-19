package main

import (
	"log"
	"net/http"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
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

	// FIXME: This should be in the database.go file, but importing the models there would result
	//        in a circular dependency.
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
		&models.Activity{},
	)
	if err != nil {
		log.Fatal("Could not migrate database:", err)
	}

	err = azure.Init()
	if err != nil {
		log.Fatal("Could not initialize Azure SDK:", err)
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

	handler := handlers.CorsMiddleware(mux)
	handlerWithLogging := handlers.LoggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
