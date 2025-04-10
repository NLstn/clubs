package frontend

import (
	"fmt"
	"os"
)

var frontendUrl string

func Init() error {
	frontendUrl = os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		return fmt.Errorf("FRONTEND_URL environment variable is not set")
	}

	return nil
}

func MakeMagicLink(token string) string {
	return fmt.Sprintf("%s/auth/magic?token=%s", frontendUrl, token)
}

func MakeClubLink(clubID string) string {
	return fmt.Sprintf("%s/clubs/%s", frontendUrl, clubID)
}
