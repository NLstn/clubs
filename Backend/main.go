package main

import (
	"log"
	"net/http"

	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	err := database.Init()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	err = azure.Init()
	if err != nil {
		log.Fatal("Could not initialize Azure SDK:", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", handlers.Handler_v1())

	handler := handlers.CorsMiddleware(mux)
	handlerWithLogging := handlers.LoggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
