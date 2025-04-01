package handlers

import (
	"log"
	"net/http"

	"github.com/NLstn/clubs/notifications"
)

func Handler_v1() http.Handler {
	mux := http.NewServeMux()

	// dummy api to test sending emails
	mux.HandleFunc("/api/v1/send-email", func(w http.ResponseWriter, r *http.Request) {
		notifier, err := notifications.NewEmailNotifier(
			"endpoint=https://clubs-acs.germany.communication.azure.com/;accesskey=YOUR_ACCESS_KEY",
			"DoNotReply@9a1159be-1efb-4f84-b252-9e949a999910.azurecomm.net",
		)
		if err != nil {
			log.Fatalf("init notifier failed: %v", err)
		}

		err = notifier.SendEmailNotification(
			"niklas.lahnstein@outlook.com",
			"Test Email",
			"This is a plain text test email",
			"<strong>This is a test email from Clubs via ACS</strong>",
		)
		if err != nil {
			log.Fatalf("send failed: %v", err)
		}
		// Set the response header to indicate success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Email sent successfully"}`))
	})

	mux.HandleFunc("/api/v1/clubs", handleClubs)
	mux.HandleFunc("/api/v1/clubs/", handleClubs)
	mux.HandleFunc("/api/v1/clubs/{clubid}/members", handleClubMembers)
	mux.HandleFunc("/api/v1/clubs/{clubid}/events", handleClubEvents)
	mux.HandleFunc("/api/v1/clubs/{clubid}/events/", handleClubEvents)

	return mux
}
