package main

import (
	"log"
	"net/http"

	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	frontend "github.com/NLstn/clubs/tools"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
		// Note: Not fatal as env vars may be set through other means in production
	}

	err = database.Init()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	err = azure.Init()
	if err != nil {
		log.Fatal("Could not initialize Azure SDK:", err)
	}

	err = frontend.Init()
	if err != nil {
		log.Fatal("Could not initialize frontend:", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", handlers.Handler_v1())

	handler := handlers.CorsMiddleware(mux)
	handlerWithLogging := handlers.LoggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
