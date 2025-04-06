package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/NLstn/clubs/azure/acs"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	dbPort, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		log.Fatal("DATABASE_PORT environment variable is required")
	}

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal("DATABASE_PORT must be an integer")
	}

	dbUser := os.Getenv("DATABASE_USER")
	if dbUser == "" {
		log.Fatal("DATABASE_USER environment variable is required")
	}

	dbUserPassword := os.Getenv("DATABASE_USER_PASSWORD")
	if dbUserPassword == "" {
		log.Fatal("DATABASE_USER_PASSWORD environment variable is required")
	}

	config := &database.Config{
		Host:     dbUrl,
		Port:     dbPortInt,
		User:     dbUser,
		Password: dbUserPassword,
		DBName:   "clubs",
	}

	err = database.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	acs.SendTestMail()

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", handlers.Handler_v1())

	handler := handlers.CorsMiddleware(mux)
	handlerWithLogging := handlers.LoggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}
